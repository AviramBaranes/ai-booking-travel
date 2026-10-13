package queries

import (
	"context"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/reservation_pricing"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var ErrCanceledReservationNotFound = api_errors.NewErrorWithDetail(
	errs.NotFound,
	"No canceled unpaid reservation with this supplier reservation id belongs to the billing entity",
	api_errors.ErrorDetails{Code: api_errors.CodeCanceledReservationNotFound},
)

// FindCanceledUnpaidReservationParams looks up a canceled, unpaid reservation of a billing unit by
// its supplier reservation id. Exactly one of OfficeID or OrgID must be provided.
type FindCanceledUnpaidReservationParams struct {
	OfficeID            int64  `query:"office_id" encore:"optional"`
	OrgID               int64  `query:"org_id" encore:"optional"`
	BrokerReservationID string `query:"broker_reservation_id" validate:"required"`
}

func (r *FindCanceledUnpaidReservationParams) Validate() error {
	if err := validation.ValidateStruct(r); err != nil {
		return err
	}

	if (r.OfficeID == 0 && r.OrgID == 0) || (r.OfficeID != 0 && r.OrgID != 0) {
		return ErrInvalidBillingEntity
	}
	return nil
}

type FindCanceledUnpaidReservationResponse struct {
	Reservation reservation_pricing.BillingReservation `json:"reservation"`
}

// FindCanceledUnpaidReservation returns a reservation that was canceled before it was paid for.
// Such a reservation is not open for billing, but an entity sometimes pays for it anyway, and the
// accountant brings it into the bill by hand to record that payment.
func (s *QueryService) FindCanceledUnpaidReservation(ctx context.Context, p *FindCanceledUnpaidReservationParams) (*FindCanceledUnpaidReservationResponse, error) {
	brokerReservationID := strings.TrimSpace(p.BrokerReservationID)
	rows, err := s.query.GetCanceledUnpaidReservationsByBillingEntity(ctx, db.GetCanceledUnpaidReservationsByBillingEntityParams{
		OfficeID:            nonZero(p.OfficeID),
		OrganizationID:      nonZero(p.OrgID),
		BrokerReservationID: &brokerReservationID,
	})
	if err != nil {
		rlog.Error("failed to fetch canceled unpaid reservation", "error", err, "officeID", p.OfficeID, "orgID", p.OrgID, "brokerReservationID", brokerReservationID)
		return nil, api_errors.ErrInternalError
	}

	if len(rows) == 0 {
		return nil, ErrCanceledReservationNotFound
	}

	return &FindCanceledUnpaidReservationResponse{
		Reservation: toBillingReservation(db.GetPaymentPendingReservationsByBillingEntityRow(rows[0])),
	}, nil
}

// ListCanceledUnpaidReservationsByIDsParams selects canceled, unpaid reservations of a billing unit
// by id. Exactly one of OfficeID or OrgID must be provided.
type ListCanceledUnpaidReservationsByIDsParams struct {
	OfficeID int64   `json:"office_id"`
	OrgID    int64   `json:"org_id"`
	IDs      []int64 `json:"ids" validate:"required,min=1"`
}

func (r *ListCanceledUnpaidReservationsByIDsParams) Validate() error {
	if err := validation.ValidateStruct(r); err != nil {
		return err
	}

	if (r.OfficeID == 0 && r.OrgID == 0) || (r.OfficeID != 0 && r.OrgID != 0) {
		return ErrInvalidBillingEntity
	}
	return nil
}

type ListCanceledUnpaidReservationsByIDsResponse struct {
	Reservations []reservation_pricing.BillingReservation `json:"reservations"`
}

// ListCanceledUnpaidReservationsByIDs returns those of the given ids that are canceled, unpaid
// reservations of the billing unit. Ids that are not are left out rather than reported.
func (s *QueryService) ListCanceledUnpaidReservationsByIDs(ctx context.Context, p *ListCanceledUnpaidReservationsByIDsParams) (*ListCanceledUnpaidReservationsByIDsResponse, error) {
	rows, err := s.query.GetCanceledUnpaidReservationsByBillingEntity(ctx, db.GetCanceledUnpaidReservationsByBillingEntityParams{
		OfficeID:       nonZero(p.OfficeID),
		OrganizationID: nonZero(p.OrgID),
		Ids:            p.IDs,
	})
	if err != nil {
		rlog.Error("failed to fetch canceled unpaid reservations", "error", err, "officeID", p.OfficeID, "orgID", p.OrgID, "ids", p.IDs)
		return nil, api_errors.ErrInternalError
	}

	reservations := make([]reservation_pricing.BillingReservation, 0, len(rows))
	for _, r := range rows {
		reservations = append(reservations, toBillingReservation(db.GetPaymentPendingReservationsByBillingEntityRow(r)))
	}

	return &ListCanceledUnpaidReservationsByIDsResponse{Reservations: reservations}, nil
}

// nonZero returns nil for a zero id, so an unset billing entity reaches the query as NULL.
func nonZero(id int64) *int64 {
	if id == 0 {
		return nil
	}
	return &id
}
