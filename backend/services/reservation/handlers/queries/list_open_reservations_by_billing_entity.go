package queries

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/reservation_pricing"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrInvalidBillingEntity = api_errors.NewValidationError("Invalid billing entity: exactly one of office_id or org_id must be provided")

	ErrOfficeInOrganicOrg = api_errors.NewErrorWithDetail(
		errs.FailedPrecondition,
		"This office belongs to an organic organization; fetch reservations at the organization level instead",
		api_errors.ErrorDetails{Code: api_errors.CodeOfficeInOrganicOrg},
	)

	ErrOrgIsInorganic = api_errors.NewErrorWithDetail(
		errs.FailedPrecondition,
		"This organization is inorganic; fetch reservations at the office level instead",
		api_errors.ErrorDetails{Code: api_errors.CodeOrgIsInorganic},
	)
)

// ListOpenReservationsByBillingEntityParams filters open reservations by a billing unit.
// Exactly one of OfficeID or OrgID must be provided.
type ListOpenReservationsByBillingEntityParams struct {
	OfficeID int64 `query:"office_id" encore:"optional"`
	OrgID    int64 `query:"org_id" encore:"optional"`
}

func (r *ListOpenReservationsByBillingEntityParams) Validate() error {
	if (r.OfficeID == 0 && r.OrgID == 0) || (r.OfficeID != 0 && r.OrgID != 0) {
		return ErrInvalidBillingEntity
	}
	return nil
}

// ListOpenReservationsByBillingEntityResponse holds the open reservations for a billing unit,
// grouped by currency.
type ListOpenReservationsByBillingEntityResponse struct {
	CurrencyGroups []CurrencyGroup `json:"currencyGroups"`
}

// CurrencyGroup is a set of billing reservations sharing the same currency.
type CurrencyGroup struct {
	CurrencyCode string                                   `json:"currencyCode"`
	Reservations []reservation_pricing.BillingReservation `json:"reservations"`
}

// ListOpenReservationsByBillingEntity returns all unpaid/refund-pending reservations
// for a given billing unit (an organic organization or an office of an inorganic organization).
func (s *QueryService) ListOpenReservationsByBillingEntity(ctx context.Context, p *ListOpenReservationsByBillingEntityParams) (*ListOpenReservationsByBillingEntityResponse, error) {
	var officeID, orgID *int64
	if p.OfficeID != 0 {
		officeID = &p.OfficeID
	} else {
		orgID = &p.OrgID
	}
	rows, err := s.query.GetPaymentPendingReservationsByBillingEntity(ctx, db.GetPaymentPendingReservationsByBillingEntityParams{
		OfficeID:       officeID,
		OrganizationID: orgID,
	})
	if err != nil {
		rlog.Error("failed to fetch reservations by billing entity", "error", err, "officeID", p.OfficeID, "orgID", p.OrgID)
		return nil, err
	}

	return &ListOpenReservationsByBillingEntityResponse{
		CurrencyGroups: toCurrencyGroups(rows),
	}, nil
}

// toCurrencyGroups maps db rows to CurrencyGroup response objects, grouping reservations
// by their currency code while preserving the order of first appearance.
func toCurrencyGroups(rows []db.GetPaymentPendingReservationsByBillingEntityRow) []CurrencyGroup {
	var groups []CurrencyGroup
	for _, r := range rows {
		groupIndex := -1
		for j, group := range groups {
			if group.CurrencyCode == r.CurrencyCode {
				groupIndex = j
				break
			}
		}

		if groupIndex == -1 {
			groups = append(groups, CurrencyGroup{
				CurrencyCode: r.CurrencyCode,
				Reservations: []reservation_pricing.BillingReservation{},
			})
			groupIndex = len(groups) - 1
		}

		pd := reservation_pricing.GetReservationPriceDetails(r)
		groups[groupIndex].Reservations = append(groups[groupIndex].Reservations, reservation_pricing.BillingReservation{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			PaymentStatus:       string(r.PaymentStatus),
			ReservationStatus:   string(r.ReservationStatus),
			CarPurchasePrice:    pd.CarPurchasePrice,
			CarSellingPrice:     pd.CarSellingPrice,
			ERPSellingPrice:     pd.ErpSellingPrice,
			TotalProfit:         pd.TotalProfit,
			TotalPrice:          pd.TotalPrice,
			CurrencyCode:        r.CurrencyCode,
			CurrencyRate:        dbadapters.NumericToFloat64(r.CurrencyRate),
			CreatedAt:           dbadapters.TimestamptzToString(r.CreatedAt),
			PickupDate:          dbadapters.DateToString(r.PickupDate),
		})
	}
	return groups
}
