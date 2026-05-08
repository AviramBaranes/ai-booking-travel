package reservation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/notifications"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/rlog"
)

// ApplyVoucherRequest is the request payload type for the apply voucher EP
type ApplyVoucherRequest struct {
	Voucher string `json:"voucher" validate:"required,notblank"`
}

func (r ApplyVoucherRequest) Validate() error {
	return validation.ValidateStruct(r)
}

// ApplyVoucher is the EP for applying a voucher on an agent order
//
// encore:api auth method=POST path=/reservations/:id/voucher tag:agent
func (s *Service) ApplyVoucher(ctx context.Context, id int64, p ApplyVoucherRequest) error {
	authData := auth.Data().(*accounts.AuthData)
	reservation, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            id,
		UserID:        authData.UserID,
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

	err = sendVoucher(b, reservation)
	if err != nil {
		rlog.Error("sending voucher", "error", err, "id", id, "voucher", p.Voucher)
		notifyVoucherError(ctx, "Voucher Send Failed", id, reservation.Broker, p.Voucher, err)
	}

	return nil
}

// notifyVoucherError publishes a critical error notification when voucher generation or sending fails.
func notifyVoucherError(ctx context.Context, subject string, id int64, b db.Broker, voucher string, err error) {
	notifications.CriticalErrorEventTopic.Publish(ctx, &notifications.CriticalErrorEvent{
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

func sendVoucher(b broker.VoucherProvider, reservation db.Reservation) error {
	voucherData, err := toVoucherData(reservation)
	if err != nil {
		return fmt.Errorf("converting to voucher data: %w", err)
	}

	htmlVoucher, err := b.GenerateVoucher(voucherData)
	if err != nil {
		return fmt.Errorf("generating voucher: %w", err)
	}

	_ = htmlVoucher
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
		PickupDate:         db.DateToString(reservation.PickupDate),
		PickupTime:         reservation.PickupTime,
		DropoffLoc:         reservation.DropoffLocationName,
		DropoffDate:        db.DateToString(reservation.ReturnDate),
		DropoffTime:        reservation.DropoffTime,
		CarGroupDesc:       carDetails.CarGroup,
		LeadModel:          carDetails.Model,
		Passengers:         carDetails.Seats,
		Suitcases:          carDetails.Bags,
		PrepaidIncludes:    reservation.PlanInclusions,
	}, nil
}
