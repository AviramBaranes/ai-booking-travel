package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/accounts/handlers/office"
	"encore.app/services/accounts/handlers/organization"
	"encore.app/services/notifications"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrNotEnoughCredits = api_errors.NewErrorWithDetail(errs.PermissionDenied, "Not enough credits for this order", api_errors.ErrorDetails{
		Code: api_errors.CodeNotEnoughCredits,
	})
	ErrNoObligo = api_errors.NewErrorWithDetail(errs.PermissionDenied, "Organization/office has no oblogo", api_errors.ErrorDetails{
		Code: api_errors.CodeNoObligo,
	})
)

// ApplyVoucherParams is the request payload type for the apply voucher EP.
type ApplyVoucherParams struct {
	Voucher string `json:"voucher" validate:"required,notblank"`
}

func (r ApplyVoucherParams) Validate() error {
	return validation.ValidateStruct(r)
}

func (s *ActionService) ApplyVoucher(ctx context.Context, id int64, p ApplyVoucherParams) error {
	authData := auth.Data().(*accounts.AuthData)
	reservation, err := s.getReservation(ctx, id, authData.UserID)
	if err != nil {
		return err
	}

	currencyRate, err := s.currencyCache.GetCurrencyRate(ctx, reservation.CurrencyCode)
	if err != nil {
		rlog.Error("failed to get currency rate", "error", err, "currency_code", reservation.CurrencyCode)
		return api_errors.ErrInternalError
	}

	tp := dbadapters.NumericToFloat64(reservation.TotalPrice)
	err = checkCredit(ctx, tp, currencyRate)
	if err != nil {
		return err
	}

	err = s.applyVoucher(ctx, id, authData.UserID, p.Voucher, currencyRate)
	if err != nil {
		return err
	}

	err = updateUserBalance(ctx, authData, tp, currencyRate)
	if err != nil {
		notifyVoucherError(ctx, "Update User Balance Failed", reservation.ID, reservation.Broker, p.Voucher, err)
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: authData.UserID})
	if err != nil {
		rlog.Error("getting user email for voucher", "error", err, "id", reservation.ID, "userID", authData.UserID)
		notifyVoucherError(ctx, "Voucher Send Failed — Could Not Resolve Recipient", reservation.ID, reservation.Broker, p.Voucher, err)
		return nil
	}

	return sendReservationVoucherToUser(ctx, reservation, userEmail.Email, p.Voucher)
}

func (s *ActionService) getReservation(ctx context.Context, id, userID int64) (db.Reservation, error) {
	reservation, err := s.query.GetReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return db.Reservation{}, api_errors.ErrNotFound
		}
		rlog.Error("getting reservation by ID", "error", err, "id", id)
		return db.Reservation{}, api_errors.ErrInternalError
	}
	if reservation.UserID != userID {
		return db.Reservation{}, api_errors.ErrNotFound
	}

	return reservation, nil
}

func (s *ActionService) applyVoucher(ctx context.Context, id, userID int64, voucher string, currencyRate float64) error {
	err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            id,
		UserID:        userID,
		CurrencyRate:  dbadapters.NumericFromFloat64(currencyRate),
		VoucherNumber: &voucher,
	})

	if err != nil {
		rlog.Error("applying voucher", "error", err, "id", id, "voucher", voucher)
		return api_errors.ErrInternalError
	}
	return nil
}

// sendReservationVoucherToUser generates and sends the reservation voucher to the user's email.
// It logs and notifies critical errors without returning them, as voucher sending failure should not block the main flow.
func sendReservationVoucherToUser(
	ctx context.Context,
	reservation db.Reservation,
	email string,
	voucherNumber string,
) error {
	b, err := getVoucherProvider(reservation.Broker)
	if err != nil {
		rlog.Error("getting voucher provider", "error", err, "id", reservation.ID, "voucher", voucherNumber)
		notifyVoucherError(ctx, "Voucher Send Failed — Unsupported Broker", reservation.ID, reservation.Broker, voucherNumber, err)
		return nil
	}

	err = sendVoucher(ctx, b, reservation, email)
	if err != nil {
		rlog.Error("sending voucher", "error", err, "id", reservation.ID, "voucher", voucherNumber)
		notifyVoucherError(ctx, "Voucher Send Failed", reservation.ID, reservation.Broker, voucherNumber, err)
	}

	return nil
}

// notifyVoucherError publishes a critical error notification when voucher generation or sending fails.
func notifyVoucherError(ctx context.Context, subject string, id int64, b db.Broker, voucher string, err error) {
	if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
		Subject: subject,
		Message: fmt.Sprintf("Reservation %d (broker: %s, voucher: %s): %v", id, b, voucher, err),
	}); publishErr != nil {
		rlog.Error("failed to publish critical error email event", "reservationId", id, "error", publishErr)
	}
}

// getVoucherProvider returns the appropriate VoucherProvider implementation based on the broker type.
func getVoucherProvider(b db.Broker) (broker.VoucherProvider, error) {
	switch b {
	case db.BrokerFlex:
		return broker.NewFlex(), nil
	case db.BrokerHertz:
		return broker.NewHertz(), nil
	default:
		return nil, errors.New("unsupported broker")
	}
}

// sendVoucher generates a voucher using the broker's VoucherProvider and sends it to the recipient's email.
func sendVoucher(ctx context.Context, b broker.VoucherProvider, reservation db.Reservation, recipientEmail string) error {
	voucherData, err := toVoucherData(reservation)
	if err != nil {
		return fmt.Errorf("converting to voucher data: %w", err)
	}

	htmlVoucher, err := b.GenerateVoucher(voucherData)
	if err != nil {
		return fmt.Errorf("generating voucher: %w", err)
	}

	if err = notifications.SendVoucher(ctx, notifications.SendVoucherParams{
		RecipientEmail:     recipientEmail,
		BookingReferenceID: reservation.BrokerReservationID,
		DriverFullName:     fmt.Sprintf("%s %s %s", reservation.DriverTitle, reservation.DriverFirstName, reservation.DriverLastName),
		VoucherNumber:      reservation.BrokerReservationID,
		VoucherHTML:        htmlVoucher,
		Broker:             notifications.VoucherBroker(reservation.Broker),
	}); err != nil {
		return fmt.Errorf("sending voucher email: %w", err)
	}

	return nil
}

func toVoucherData(reservation db.Reservation) (*broker.VoucherData, error) {
	var carDetails broker.CarDetails
	if err := json.Unmarshal(reservation.CarDetails, &carDetails); err != nil {
		return nil, fmt.Errorf("unmarshalling car details: %w", err)
	}

	return &broker.VoucherData{
		ReservationNum:     reservation.BrokerReservationID,
		BookingReferenceID: reservation.BrokerReservationID,
		CustomerName:       strings.TrimSpace(reservation.DriverTitle + " " + reservation.DriverFirstName + " " + reservation.DriverLastName),
		Supplier:           reservation.SupplierCode,
		PickupLoc:          reservation.PickupLocationName,
		PickupDate:         dbadapters.DateToString(reservation.PickupDate),
		PickupTime:         reservation.PickupTime,
		DropoffLoc:         reservation.DropoffLocationName,
		DropoffDate:        dbadapters.DateToString(reservation.DropoffDate),
		DropoffTime:        reservation.DropoffTime,
		CarGroupDesc:       carDetails.CarGroup,
		LeadModel:          carDetails.Model,
		Passengers:         carDetails.Seats,
		Suitcases:          carDetails.Bags,
		PrepaidIncludes:    reservation.PlanInclusions,
	}, nil
}

// checkCredit retrieves the balance of the user billing entity and checks if it is sufficient to cover the total price of the reservation after applying the currency rate.
func checkCredit(ctx context.Context, totalPrice float64, currencyRate float64) error {
	balance, err := accounts.GetUserCredit(ctx)
	if err != nil {
		rlog.Error("getting user credit", "error", err)
		return api_errors.ErrInternalError
	}

	if balance.Obligo <= 0 {
		return ErrNoObligo
	}

	adjustedPrice := totalPrice * currencyRate
	credit := balance.Obligo - balance.BalanceDue
	if credit < adjustedPrice {
		return ErrNotEnoughCredits
	}

	return nil
}

// updateUserBalance updates the balance due for the user's billing entity (organization or office) based on the total price of the reservation and the currency rate. It handles both organic and inorganic organizations appropriately.
func updateUserBalance(ctx context.Context, authData *accounts.AuthData, totalPrice float64, currencyRate float64) error {
	adjustedPrice := totalPrice * currencyRate
	if authData.OrganizationContext.IsOrganic {
		if err := accounts.UpdateOrganizationBalanceDue(ctx, organization.UpdateOrganizationBalanceDueParams{
			ID:            authData.OrganizationContext.OrganizationID,
			BalanceChange: adjustedPrice,
		}); err != nil {
			rlog.Error("updating organization balance due", "error", err)
			return err
		}
		return nil
	}
	if err := accounts.UpdateOfficeBalanceDue(ctx, office.UpdateOfficeBalanceDueParams{
		ID:            authData.OrganizationContext.OfficeID,
		BalanceChange: adjustedPrice,
	}); err != nil {
		rlog.Error("updating office balance due", "error", err)
		return err
	}
	return nil
}
