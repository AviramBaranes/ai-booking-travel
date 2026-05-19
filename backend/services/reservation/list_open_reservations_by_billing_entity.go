package reservation

import (
	"context"
	"math"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
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
	TotalPrice          float64 `json:"totalPrice"`
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
func (s *Service) ListOpenReservationsByBillingEntity(ctx context.Context, p *ListOpenReservationsByBillingEntityParams) (*ListOpenReservationsByBillingEntityResponse, error) {
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
			TotalPrice:          pd.totalPrice,
			CurrencyCode:        r.CurrencyCode,
			CreatedAt:           dbadapters.TimestamptzToString(r.CreatedAt),
			PickupDate:          dbadapters.DateToString(r.PickupDate),
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
	totalPrice       float64
}

// roundPrice rounds a price to 2 decimal places.
func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}

// getReservationPriceDetails computes purchase price, selling price, profit, and ERP price from a db row.
func getReservationPriceDetails(row db.GetPaymentPendingReservationsByBillingEntityRow) priceDetails {
	carPurchasePrice := dbadapters.NumericToFloat64(row.PurchasePrice) + dbadapters.NumericToFloat64(row.BrokerErpPrice)
	mp := dbadapters.NumericToFloat64(row.MarkupPercentage)
	carSellingPrice := pricing.ApplyMarkup(carPurchasePrice, mp)

	return priceDetails{
		carPurchasePrice: roundPrice(carPurchasePrice),
		carSellingPrice:  roundPrice(carSellingPrice),
		carProfit:        roundPrice(carSellingPrice - carPurchasePrice),
		erpSellingPrice:  roundPrice(float64(row.BtErpPrice)),
		totalPrice:       roundPrice(carSellingPrice + float64(row.BtErpPrice)),
	}
}
