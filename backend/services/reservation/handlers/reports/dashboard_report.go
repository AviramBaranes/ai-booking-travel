package reports

import (
	"context"
	"encoding/json"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

// Gear types reported to the dashboard. Derived from car_details, which stores the
// gearbox as two booleans rather than a single field.
const (
	GearTypeAuto     = "auto"
	GearTypeManual   = "manual"
	GearTypeElectric = "electric"
)

type DashboardParams struct {
	// From and To are calendar dates ("2006-01-02") in the business timezone, inclusive.
	From string `query:"from" validate:"required"`
	To   string `query:"to" validate:"required"`
}

func (p *DashboardParams) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}

	from, err := time.ParseInLocation("2006-01-02", p.From, reportLocation)
	if err != nil {
		return api_errors.NewValidationError("from must be a calendar date formatted as YYYY-MM-DD")
	}
	to, err := time.ParseInLocation("2006-01-02", p.To, reportLocation)
	if err != nil {
		return api_errors.NewValidationError("to must be a calendar date formatted as YYYY-MM-DD")
	}
	if to.Before(from) {
		return api_errors.NewValidationError("to must not be earlier than from")
	}

	return nil
}

// DashboardReservation is one reservation as the dashboard sees it: the dimensions it
// slices by, plus a money breakdown that is already derived and already converted to ILS.
// The client only ever sums, groups and averages these numbers — it never re-derives them.
type DashboardReservation struct {
	ReservationID int64  `json:"reservationId"`
	CreatedAt     string `json:"createdAt"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentStatus"`
	// IsVouchered reports whether a voucher was ever issued, not the current status: a
	// reservation cancelled after ticketing stays vouchered here, so the dashboard can
	// narrow to ticketed business and still account for its cancellations.
	IsVouchered bool `json:"isVouchered"`

	IsBusiness            bool   `json:"isBusiness"`
	UserID                int64  `json:"userId"`
	OfficeID              *int64 `json:"officeId,omitempty"`
	OrganizationID        *int64 `json:"organizationId,omitempty"`
	IsOrganizationOrganic *bool  `json:"isOrganizationOrganic,omitempty"`

	Broker       string `json:"broker"`
	SupplierCode string `json:"supplierCode"`
	SupplierName string `json:"supplierName"`
	CarType      string `json:"carType"`
	CarGroup     string `json:"carGroup"`
	GearType     string `json:"gearType"`
	CountryCode  string `json:"countryCode"`

	PickupDate   string `json:"pickupDate"`
	RentalDays   int32  `json:"rentalDays"`
	LeadTimeDays int32  `json:"leadTimeDays"`
	DriverAge    int32  `json:"driverAge"`
	CouponName   string `json:"couponName"`
	HasERP       bool   `json:"hasErp"`
	CurrencyCode string `json:"currencyCode"`

	// Money, all in ILS.
	RevenueILS    float64 `json:"revenueIls"`
	CostILS       float64 `json:"costIls"`
	ProfitILS     float64 `json:"profitIls"`
	ErpRevenueILS float64 `json:"erpRevenueIls"`
	ErpCostILS    float64 `json:"erpCostIls"`
	DiscountILS   float64 `json:"discountIls"`

	SupplierPaid        bool    `json:"supplierPaid"`
	PenaltyType         string  `json:"penaltyType,omitempty"`
	PenaltyAmountILS    float64 `json:"penaltyAmountIls"`
	PenaltyPaid         bool    `json:"penaltyPaid"`
	PenaltySupplierPaid bool    `json:"penaltySupplierPaid"`
}

// DashboardResponse carries the rows plus the id→name lookups they reference. Names are
// sent once here rather than repeated on every row.
type DashboardResponse struct {
	Reservations  []DashboardReservation `json:"reservations"`
	Organizations []accounts.AccountName `json:"organizations"`
	Offices       []accounts.AccountName `json:"offices"`
	Users         []accounts.AccountName `json:"users"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
}

func (s *ReportsService) GetDashboardReport(ctx context.Context, p *DashboardParams) (*DashboardResponse, error) {
	rows, err := s.query.ListReservationsForDashboard(ctx, db.ListReservationsForDashboardParams{
		CreatedFrom: timestamptzFromString(p.From, false),
		CreatedTo:   timestamptzFromString(p.To, true),
	})
	if err != nil {
		rlog.Error("failed to list reservations for dashboard", "error", err)
		return nil, api_errors.ErrInternalError
	}

	accountsLookup, err := accounts.GetAccountsLookup(ctx, accounts.GetAccountsLookupParams{
		OrganizationIDs: idsFromSet(dashboardOrganizationIDs(rows)),
		OfficeIDs:       idsFromSet(dashboardOfficeIDs(rows)),
		UserIDs:         idsFromSet(dashboardUserIDs(rows)),
	})
	if err != nil {
		rlog.Error("failed to get accounts lookup for dashboard", "error", err)
		return nil, err
	}

	reservations, err := buildDashboardRows(rows)
	if err != nil {
		return nil, err
	}

	return &DashboardResponse{
		Reservations:  reservations,
		Organizations: accountsLookup.Organizations,
		Offices:       accountsLookup.Offices,
		Users:         accountsLookup.Users,
		From:          p.From,
		To:            p.To,
	}, nil
}

func buildDashboardRows(rows []db.ListReservationsForDashboardRow) ([]DashboardReservation, error) {
	out := make([]DashboardReservation, 0, len(rows))

	for _, r := range rows {
		var carDetails broker.CarDetails
		if err := json.Unmarshal(r.CarDetails, &carDetails); err != nil {
			rlog.Error("failed to unmarshal car details for dashboard", "reservation_id", r.ID, "error", err)
			return nil, api_errors.ErrInternalError
		}

		price := pricing.NewWithNumerics(pricing.NumericParams{
			PurchasePrice:      r.PurchasePrice,
			BrokerErpPrice:     r.BrokerErpPrice,
			BtErpPrice:         r.BtErpPrice,
			MarkupPercentage:   r.MarkupPercentage,
			DiscountPercentage: r.DiscountPercentage,
			CurrencyCode:       r.CurrencyCode,
			CurrencyRate:       r.CurrencyRate,
			TotalPrice:         r.TotalPrice,
		}).Details()

		out = append(out, DashboardReservation{
			ReservationID: r.ID,
			CreatedAt:     dbadapters.TimestamptzToString(r.CreatedAt),
			Status:        string(r.ReservationStatus),
			PaymentStatus: string(r.PaymentStatus),
			IsVouchered:   r.VoucheredAt.Valid,

			IsBusiness:            r.OfficeID != nil && r.OrganizationID != nil,
			UserID:                r.UserID,
			OfficeID:              r.OfficeID,
			OrganizationID:        r.OrganizationID,
			IsOrganizationOrganic: r.IsOrganizationOrganic,

			Broker:       string(r.Broker),
			SupplierCode: r.SupplierCode,
			SupplierName: carDetails.SupplierName,
			CarType:      carDetails.CarType,
			CarGroup:     carDetails.CarGroup,
			GearType:     gearType(carDetails),
			CountryCode:  r.CountryCode,

			PickupDate:   dbadapters.DateToString(r.PickupDate),
			RentalDays:   r.RentalDays,
			LeadTimeDays: leadTimeDays(r.CreatedAt, r.PickupDate),
			DriverAge:    r.DriverAge,
			CouponName:   r.CouponName,
			HasERP:       price.BtErpPrice.Value > 0 || price.BrokerErpCost.Value > 0,
			CurrencyCode: price.CurrencyCode,

			RevenueILS:    price.TotalPrice.ILS,
			CostILS:       price.TotalCost.ILS,
			ProfitILS:     price.TotalProfit.ILS,
			ErpRevenueILS: price.BtErpPrice.ILS,
			ErpCostILS:    price.BrokerErpCost.ILS,
			DiscountILS:   price.TotalDiscount.ILS,

			SupplierPaid:        r.SupplierPaidAt.Valid,
			PenaltyType:         penaltyType(r.PenaltyType),
			PenaltyAmountILS:    dbadapters.NumericToFloat64(r.PenaltyAmount) * dbadapters.NumericToFloat64(r.PenaltyCurrencyRate),
			PenaltyPaid:         r.PenaltyPaidAt.Valid,
			PenaltySupplierPaid: r.PenaltySupplierPaidAt.Valid,
		})
	}

	return out, nil
}

func gearType(carDetails broker.CarDetails) string {
	switch {
	case carDetails.IsElectric:
		return GearTypeElectric
	case carDetails.IsAutoGear:
		return GearTypeAuto
	default:
		return GearTypeManual
	}
}

// leadTimeDays is how far ahead of pickup the reservation was made. Negative values are
// clamped to zero: a same-day booking made after midnight would otherwise read as -1.
func leadTimeDays(createdAt dbadapters.Timestamptz, pickupDate dbadapters.Date) int32 {
	if !createdAt.Valid || !pickupDate.Valid {
		return 0
	}

	created := dbadapters.TimeFromDB(createdAt).In(reportLocation)
	createdDay := time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, reportLocation)
	pickup := pickupDate.Time.In(time.UTC)
	pickupDay := time.Date(pickup.Year(), pickup.Month(), pickup.Day(), 0, 0, 0, 0, reportLocation)

	days := int32(pickupDay.Sub(createdDay).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func penaltyType(t *db.PenaltyType) string {
	if t == nil {
		return ""
	}
	return string(*t)
}

func dashboardOrganizationIDs(rows []db.ListReservationsForDashboardRow) map[int64]struct{} {
	ids := make(map[int64]struct{})
	for _, r := range rows {
		if r.OrganizationID != nil {
			ids[*r.OrganizationID] = struct{}{}
		}
	}
	return ids
}

func dashboardOfficeIDs(rows []db.ListReservationsForDashboardRow) map[int64]struct{} {
	ids := make(map[int64]struct{})
	for _, r := range rows {
		if r.OfficeID != nil {
			ids[*r.OfficeID] = struct{}{}
		}
	}
	return ids
}

func dashboardUserIDs(rows []db.ListReservationsForDashboardRow) map[int64]struct{} {
	ids := make(map[int64]struct{})
	for _, r := range rows {
		ids[r.UserID] = struct{}{}
	}
	return ids
}
