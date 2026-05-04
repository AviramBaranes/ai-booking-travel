package reservation

import (
	"context"
	"math"

	"encore.app/internal/api_errors"
	"encore.app/internal/pricing"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrInvalidBillingEntity = api_errors.NewValidationError("Invalid billing entity: exactly one of office_id or org_id must be provided")
	// ErrOfficeInOrganicOrg is returned when trying to fetch reservations by office for an office
	// that belongs to an organic organization — organic orgs are billed at the org level, not office level.
	ErrOfficeInOrganicOrg = api_errors.NewErrorWithDetail(
		errs.FailedPrecondition,
		"This office belongs to an organic organization; fetch reservations at the organization level instead",
		api_errors.ErrorDetails{Code: api_errors.CodeOfficeInOrganicOrg},
	)
	// ErrOrgIsInorganic is returned when trying to fetch reservations by org for an inorganic organization —
	// inorganic orgs are billed per office, so the accountant must specify an office.
	ErrOrgIsInorganic = api_errors.NewErrorWithDetail(
		errs.FailedPrecondition,
		"This organization is inorganic; fetch reservations at the office level instead",
		api_errors.ErrorDetails{Code: api_errors.CodeOrgIsInorganic},
	)
)

// ListOpenReservationsByBillingEntityRequest filters open reservations by a billing unit.
// Exactly one of OfficeID or OrgID must be provided.
type ListOpenReservationsByBillingEntityRequest struct {
	OfficeID int32 `query:"office_id" encore:"optional"`
	OrgID    int32 `query:"org_id" encore:"optional"`
}

func (r *ListOpenReservationsByBillingEntityRequest) Validate() error {
	if (r.OfficeID == 0 && r.OrgID == 0) || (r.OfficeID != 0 && r.OrgID != 0) {
		return ErrInvalidBillingEntity
	}

	return nil
}

// BillingReservation is a reservation summary tailored for accountant billing workflows.
type BillingReservation struct {
	ID                  int64   `json:"id"`
	BrokerReservationID string  `json:"brokerReservationId"`
	PaymentStatus       string  `json:"paymentStatus"`
	ReservationStatus   string  `json:"reservationStatus"`
	CarPurchasePrice    float64 `json:"carPurchasePrice"`
	CarSellingPrice     float64 `json:"carSellingPrice"`
	ERPSellingPrice     float64 `json:"erpSellingPrice"`
	ProfitOnCar         float64 `json:"profitOnCar"`
	CurrencyCode        string  `json:"currencyCode"`
	CreatedAt           string  `json:"createdAt"`
	PickupDate          string  `json:"pickupDate"`
}

// ListOpenReservationsByBillingEntityResponse holds the open reservations for a billing unit,
// grouped by currency.
type ListOpenReservationsByBillingEntityResponse struct {
	CurrencyGroups []CurrencyGroup `json:"currencyGroups"`
}

// CurrencyGroup is a set of billing reservations sharing the same currency.
type CurrencyGroup struct {
	CurrencyCode string               `json:"currencyCode"`
	Reservations []BillingReservation `json:"reservations"`
}

// ListOpenReservationsByBillingEntity returns all unpaid/refund-pending reservations
// for a given billing unit (an organic organization or an office of an inorganic organization).
//
//encore:api auth method=GET path=/reservations-for-billing tag:accountant
func (s *Service) ListOpenReservationsByBillingEntity(ctx context.Context, req *ListOpenReservationsByBillingEntityRequest) (*ListOpenReservationsByBillingEntityResponse, error) {
	agentIDs, err := s.getAgentsByBillingEntity(ctx, req.OfficeID, req.OrgID)
	if err != nil {
		return nil, err
	}

	rows, err := s.query.GetPaymentPendingReservationsByAgentsIDs(ctx, agentIDs)
	if err != nil {
		rlog.Error("failed to fetch reservations by billing entity", "error", err, "agent_ids", agentIDs)
		return nil, err
	}

	return &ListOpenReservationsByBillingEntityResponse{
		CurrencyGroups: toCurrencyGroups(rows),
	}, nil
}

// getAgentsByBillingEntity resolves agent IDs from the given billing unit and validates
// that the billing unit type matches the organization's organic setting.
func (s *Service) getAgentsByBillingEntity(ctx context.Context, officeID, orgID int32) ([]int32, error) {
	if officeID != 0 {
		r, err := accounts.GetAgentsByOfficeID(ctx, accounts.GetAgentsByOfficeIDRequest{
			OfficeID: officeID,
		})
		if err != nil {
			return nil, err
		}

		if r.IsOrganic {
			return nil, ErrOfficeInOrganicOrg
		}

		return r.IDs, nil
	}

	r, err := accounts.GetAgentsByOrganizationID(ctx, accounts.GetAgentsByOrganizationIDRequest{
		OrgID: orgID,
	})
	if err != nil {
		return nil, err
	}

	if !r.IsOrganic {
		return nil, ErrOrgIsInorganic
	}

	return r.IDs, nil
}

// toCurrencyGroups maps db rows to CurrencyGroup response objects, grouping reservations
// by their currency code while preserving the order of first appearance.
func toCurrencyGroups(rows []db.GetPaymentPendingReservationsByAgentsIDsRow) []CurrencyGroup {
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
				Reservations: []BillingReservation{},
			})
			groupIndex = len(groups) - 1
		}

		pd := getReservationPriceDetails(r)
		groups[groupIndex].Reservations = append(groups[groupIndex].Reservations, BillingReservation{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			PaymentStatus:       string(r.PaymentStatus),
			ReservationStatus:   string(r.ReservationStatus),
			CarPurchasePrice:    pd.carPurchasePrice,
			CarSellingPrice:     pd.carSellingPrice,
			ERPSellingPrice:     pd.erpSellingPrice,
			ProfitOnCar:         pd.carProfit,
			CurrencyCode:        r.CurrencyCode,
			CreatedAt:           db.TimestamptzToString(r.CreatedAt),
			PickupDate:          db.DateToString(r.PickupDate),
		})
	}
	return groups
}

// priceDetails holds the computed price breakdown for a single reservation.
type priceDetails struct {
	carPurchasePrice float64
	carSellingPrice  float64
	carProfit        float64
	erpSellingPrice  float64
}

// roundPrice rounds a price to 2 decimal places.
func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}

// getReservationPriceDetails computes purchase price, selling price, profit, and ERP price from a db row.
func getReservationPriceDetails(row db.GetPaymentPendingReservationsByAgentsIDsRow) priceDetails {
	carPurchasePrice := db.NumericToFloat64(row.PurchasePrice) + db.NumericToFloat64(row.BrokerErpPrice)
	mp := db.NumericToFloat64(row.MarkupPercentage)
	carSellingPrice := pricing.ApplyMarkup(carPurchasePrice, mp)

	return priceDetails{
		carPurchasePrice: roundPrice(carPurchasePrice),
		carSellingPrice:  roundPrice(carSellingPrice),
		carProfit:        roundPrice(carSellingPrice - carPurchasePrice),
		erpSellingPrice:  roundPrice(float64(row.BtErpPrice)),
	}
}
