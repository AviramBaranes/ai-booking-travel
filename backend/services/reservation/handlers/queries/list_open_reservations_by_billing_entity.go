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

// CurrencyGroup is a set of billing reservations and penalties sharing the same currency.
// They are grouped together because an invoice covers a single currency, so the accountant
// settles reservations and fees of one currency in the same document.
type CurrencyGroup struct {
	CurrencyCode string                                   `json:"currencyCode"`
	Reservations []reservation_pricing.BillingReservation `json:"reservations"`
	Penalties    []reservation_pricing.BillingPenalty     `json:"penalties"`
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

	penaltyRows, err := s.query.GetPaymentPendingPenaltiesByBillingEntity(ctx, db.GetPaymentPendingPenaltiesByBillingEntityParams{
		OfficeID:       officeID,
		OrganizationID: orgID,
	})
	if err != nil {
		rlog.Error("failed to fetch penalties by billing entity", "error", err, "officeID", p.OfficeID, "orgID", p.OrgID)
		return nil, err
	}

	return &ListOpenReservationsByBillingEntityResponse{
		CurrencyGroups: toCurrencyGroups(rows, penaltyRows),
	}, nil
}

// toCurrencyGroups maps db rows to CurrencyGroup response objects, grouping reservations and
// penalties by their currency code while preserving the order of first appearance.
func toCurrencyGroups(
	rows []db.GetPaymentPendingReservationsByBillingEntityRow,
	penaltyRows []db.GetPaymentPendingPenaltiesByBillingEntityRow,
) []CurrencyGroup {
	var groups []CurrencyGroup

	// groupIndexFor returns the index of the group for a currency, creating it if needed.
	groupIndexFor := func(currencyCode string) int {
		for i, group := range groups {
			if group.CurrencyCode == currencyCode {
				return i
			}
		}

		groups = append(groups, CurrencyGroup{
			CurrencyCode: currencyCode,
			Reservations: []reservation_pricing.BillingReservation{},
			Penalties:    []reservation_pricing.BillingPenalty{},
		})
		return len(groups) - 1
	}

	for _, r := range rows {
		groupIndex := groupIndexFor(r.CurrencyCode)

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
			VoucheredAt:         dbadapters.TimestamptzToString(r.VoucheredAt),
		})
	}

	for _, p := range penaltyRows {
		groupIndex := groupIndexFor(p.CurrencyCode)

		groups[groupIndex].Penalties = append(groups[groupIndex].Penalties, reservation_pricing.BillingPenalty{
			ID:                  p.ID,
			ReservationID:       p.ReservationID,
			BrokerReservationID: p.BrokerReservationID,
			Type:                string(p.PenaltyType),
			Amount:              dbadapters.NumericToFloat64(p.Amount),
			CurrencyCode:        p.CurrencyCode,
			CurrencyRate:        dbadapters.NumericToFloat64(p.CurrencyRate),
			CreatedAt:           dbadapters.TimestamptzToString(p.CreatedAt),
		})
	}

	return groups
}
