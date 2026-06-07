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
	emailevents "encore.app/internal/email_events"
	"encore.app/internal/icount"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/notifications"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/rlog"
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

	currencyRate, err := s.getCurrencyRate(ctx, id)
	if err != nil {
		return err
	}
	reservation, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            id,
		UserID:        authData.UserID,
		CurrencyRate:  dbadapters.NumericFromFloat64(currencyRate),
		VoucherNumber: &p.Voucher,
	})

	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("applying voucher", "error", err, "id", id, "voucher", p.Voucher)
		return api_errors.ErrInternalError
	}

	b, err := getVoucherProvider(reservation.Broker)
	if err != nil {
		rlog.Error("getting voucher provider", "error", err, "id", id, "voucher", p.Voucher)
		notifyVoucherError(ctx, "Voucher Send Failed — Unsupported Broker", id, reservation.Broker, p.Voucher, err)
		return nil
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: authData.UserID})
	if err != nil {
		rlog.Error("getting user email for voucher", "error", err, "id", id, "userID", authData.UserID)
		notifyVoucherError(ctx, "Voucher Send Failed — Could Not Resolve Recipient", id, reservation.Broker, p.Voucher, err)
		return nil
	}

	err = sendVoucher(ctx, b, reservation, userEmail.Email)
	if err != nil {
		rlog.Error("sending voucher", "error", err, "id", id, "voucher", p.Voucher)
		notifyVoucherError(ctx, "Voucher Send Failed", id, reservation.Broker, p.Voucher, err)
	}

	return nil
}

// notifyVoucherError publishes a critical error notification when voucher generation or sending fails.
func notifyVoucherError(ctx context.Context, subject string, id int64, b db.Broker, voucher string, err error) {
	emailevents.PublishEmailEvent(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
		Subject: subject,
		Message: fmt.Sprintf("Reservation %d (broker: %s, voucher: %s): %v", id, b, voucher, err),
	})
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

// getCurrencyRate retrieves the currency code for the reservation and fetches the corresponding exchange rate from Icount.
func (s *ActionService) getCurrencyRate(ctx context.Context, reservationID int64) (float64, error) {
	currencyCode, err := s.query.GetReservationCurrencyCode(ctx, reservationID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return 0, api_errors.ErrNotFound
		}
		rlog.Error("getting reservation currency code", "error", err, "id", reservationID)
		return 0, api_errors.ErrInternalError
	}

	ic := icount.NewIcount(s.cfg.IcountCID, s.cfg.IcountUser)
	resp, err := ic.FetchCurrencies()
	if err != nil {
		rlog.Error("fetching currency rates from icount", "error", err)
		return 0, api_errors.ErrInternalError
	}

	rate, ok := resp.Rates[currencyCode]
	if !ok {
		rlog.Error("currency code not found in icount rates", "currencyCode", currencyCode)
		return 0, api_errors.ErrInternalError
	}

	return rate, nil
}
