package reports

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

func (p ReportParams) Validate() error {
	return validation.ValidateStruct(p)
}

type BusinessReservationReportRow struct {
	ReservationID                  int64   `json:"reservationId"`
	BrokerReservationID            string  `json:"brokerReservationId"`
	Status                         string  `json:"status"`
	OrganizationName               string  `json:"organizationName"`
	OfficeName                     string  `json:"officeName"`
	UserName                       string  `json:"userName"`
	AdminName                      *string `json:"adminName,omitempty"`
	BrokerName                     string  `json:"brokerName"`
	SupplierName                   string  `json:"supplierName"`
	CountryCode                    string  `json:"countryCode"`
	PickupDate                     string  `json:"pickupDate"`
	DropoffDate                    string  `json:"dropoffDate"`
	RentalDays                     int32   `json:"rentalDays"`
	DriverName                     string  `json:"driverName"`
	CurrencyCode                   string  `json:"currencyCode"`
	CurrencyRate                   float64 `json:"currencyRate"`
	CarSellPriceWithBrokerERP      float64 `json:"carSellPriceWithBrokerERP"`
	CarSellPriceWithBrokerERPInILS float64 `json:"carSellPriceWithBrokerERPInILS"`
	BtERPPrice                     float64 `json:"btERPPrice"`
	BtERPPriceInILS                float64 `json:"btERPPriceInILS"`
	TotalPrice                     float64 `json:"totalPrice"`
	TotalPriceInILS                float64 `json:"totalPriceInILS"`
	CouponName                     string  `json:"couponName"`
	VoucherNumber                  *string `json:"voucherNumber,omitempty"`
	VoucheredAt                    *string `json:"voucheredAt,omitempty"`
	CreatedAt                      string  `json:"createdAt"`
}

type BusinessReportResponse struct {
	Reservations []BusinessReservationReportRow `json:"reservations"`
	Count        int64                          `json:"count"`
	TotalSales   float64                        `json:"totalSales"`
}

type accountsSet struct {
	organizationIDs map[int64]struct{}
	officeIDs       map[int64]struct{}
	userIDs         map[int64]struct{}
}

func (s *ReportsService) GetBusinessReport(ctx context.Context, p ReportParams) (*BusinessReportResponse, error) {
	result, err := s.getReports(ctx, p)
	if err != nil {
		return nil, err
	}

	accountsSet := buildAccountsSet(result.Reservations)
	accountsLookup, err := accounts.GetAccountsLookup(ctx, accounts.GetAccountsLookupParams{
		OrganizationIDs: idsFromSet(accountsSet.organizationIDs),
		OfficeIDs:       idsFromSet(accountsSet.officeIDs),
		UserIDs:         idsFromSet(accountsSet.userIDs),
	})
	if err != nil {
		rlog.Error("failed to get accounts lookup for business report", "error", err)
		return nil, err
	}

	reservations, err := buildBusinessReportRows(result.Reservations, accountsLookup)
	if err != nil {
		return nil, err
	}

	return &BusinessReportResponse{
		Reservations: reservations,
		Count:        result.Count,
		TotalSales:   result.TotalSales,
	}, nil
}

func buildAccountsSet(reservations []db.Reservation) accountsSet {
	orgIDs := make(map[int64]struct{})
	officeIDs := make(map[int64]struct{})
	usersIDs := make(map[int64]struct{})

	for _, r := range reservations {
		if r.OrganizationID != nil {
			orgIDs[*r.OrganizationID] = struct{}{}
		}
		if r.OfficeID != nil {
			officeIDs[*r.OfficeID] = struct{}{}
		}
		usersIDs[r.UserID] = struct{}{}
		if r.AdminRefID != nil {
			usersIDs[*r.AdminRefID] = struct{}{}
		}
	}

	return accountsSet{
		organizationIDs: orgIDs,
		officeIDs:       officeIDs,
		userIDs:         usersIDs,
	}
}

func buildBusinessReportRows(reservations []db.Reservation, accountsLookup *accounts.GetAccountsLookupResponse) ([]BusinessReservationReportRow, error) {
	organizationNames := namesByID(accountsLookup.Organizations)
	officeNames := namesByID(accountsLookup.Offices)
	userNames := namesByID(accountsLookup.Users)

	rows := make([]BusinessReservationReportRow, 0, len(reservations))
	for _, r := range reservations {
		var carDetails broker.CarDetails
		if err := json.Unmarshal(r.CarDetails, &carDetails); err != nil {
			rlog.Error("failed to unmarshal car details for business report", "reservation_id", r.ID, "error", err)
			return nil, api_errors.ErrInternalError
		}

		currencyRate := dbadapters.NumericToFloat64(r.CurrencyRate)
		carSellPriceWithBrokerERP := calculateCarSellPriceWithBrokerERP(r)
		btERPPrice := dbadapters.NumericToFloat64(r.BtErpPrice)
		totalPrice := dbadapters.NumericToFloat64(r.TotalPrice)
		adminName := namePtr(userNames, r.AdminRefID)
		voucheredAt := optionalString(dbadapters.TimestamptzToString(r.VoucheredAt))

		rows = append(rows, BusinessReservationReportRow{
			ReservationID:                  r.ID,
			BrokerReservationID:            r.BrokerReservationID,
			Status:                         string(r.ReservationStatus),
			OrganizationName:               nameForID(organizationNames, r.OrganizationID),
			OfficeName:                     nameForID(officeNames, r.OfficeID),
			UserName:                       userNames[r.UserID],
			AdminName:                      adminName,
			BrokerName:                     string(r.Broker),
			SupplierName:                   carDetails.SupplierName,
			CountryCode:                    r.CountryCode,
			PickupDate:                     dbadapters.DateToString(r.PickupDate),
			DropoffDate:                    dbadapters.DateToString(r.DropoffDate),
			RentalDays:                     r.RentalDays,
			DriverName:                     fmt.Sprintf("%s %s %s", r.DriverTitle, r.DriverFirstName, r.DriverLastName),
			CurrencyCode:                   r.CurrencyCode,
			CurrencyRate:                   currencyRate,
			CarSellPriceWithBrokerERP:      carSellPriceWithBrokerERP,
			CarSellPriceWithBrokerERPInILS: carSellPriceWithBrokerERP * currencyRate,
			BtERPPrice:                     btERPPrice,
			BtERPPriceInILS:                btERPPrice * currencyRate,
			TotalPrice:                     totalPrice,
			TotalPriceInILS:                totalPrice * currencyRate,
			CouponName:                     r.CouponName,
			VoucherNumber:                  r.VoucherNumber,
			VoucheredAt:                    voucheredAt,
			CreatedAt:                      dbadapters.TimestamptzToString(r.CreatedAt),
		})
	}

	return rows, nil
}

func calculateCarSellPriceWithBrokerERP(reservation db.Reservation) float64 {
	purchasePrice := dbadapters.NumericToFloat64(reservation.PurchasePrice)
	markupPercentage := dbadapters.NumericToFloat64(reservation.MarkupPercentage)
	brokerERPPrice := dbadapters.NumericToFloat64(reservation.BrokerErpPrice)

	return pricing.ApplyMarkup(purchasePrice, markupPercentage) + pricing.ApplyMarkup(brokerERPPrice, markupPercentage)
}

func namesByID(rows []accounts.AccountName) map[int64]string {
	names := make(map[int64]string, len(rows))
	for _, row := range rows {
		names[row.ID] = row.Name
	}
	return names
}

func nameForID(names map[int64]string, id *int64) string {
	if id == nil {
		return ""
	}
	return names[*id]
}

func namePtr(names map[int64]string, id *int64) *string {
	if id == nil {
		return nil
	}
	name, ok := names[*id]
	if !ok {
		return nil
	}
	return &name
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func idsFromSet(set map[int64]struct{}) []int64 {
	ids := make([]int64, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}
