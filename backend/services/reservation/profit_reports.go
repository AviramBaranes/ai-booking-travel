package reservation

import (
	"context"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type ProfitReportRow struct {
	BusinessReservationReportRow
	PurchasePrice      float64 `json:"purchasePrice"`
	PurchasePriceInILS float64 `json:"purchasePriceInILS"`
	Profit             float64 `json:"profit"`
	ProfitInILS        float64 `json:"profitInILS"`
}

type ProfitReportResponse struct {
	Reservations []ProfitReportRow `json:"reservations"`
	Total        int64             `json:"total"`
}

// encore:api auth tag:admin method=GET path=/reports/profit
func (s *Service) GetProfitReport(ctx context.Context, p ReportParams) (*ProfitReportResponse, error) {
	rows, total, err := s.getReports(ctx, p, true)
	if err != nil {
		return nil, err
	}

	accountsSet := buildAccountsSet(rows)
	accountsLookup, err := accounts.GetAccountsLookup(ctx, accounts.GetAccountsLookupParams{
		OrganizationIDs: idsFromSet(accountsSet.organizationIDs),
		OfficeIDs:       idsFromSet(accountsSet.officeIDs),
		UserIDs:         idsFromSet(accountsSet.userIDs),
	})
	if err != nil {
		rlog.Error("failed to get accounts lookup for profit report", "error", err)
		return nil, err
	}

	reservations, err := buildProfitReportRows(rows, accountsLookup)
	if err != nil {
		return nil, err
	}

	return &ProfitReportResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}

func buildProfitReportRows(reservations []db.Reservation, accountsLookup *accounts.GetAccountsLookupResponse) ([]ProfitReportRow, error) {
	businessRows, err := buildBusinessReportRows(reservations, accountsLookup)
	if err != nil {
		return nil, err
	}

	rows := make([]ProfitReportRow, 0, len(reservations))
	for i, r := range reservations {
		currencyRate := dbadapters.NumericToFloat64(r.CurrencyRate)
		purchasePrice := calculateCarPurchasePriceWithBrokerERP(r)
		profit := businessRows[i].CarSellPriceWithBrokerERP - purchasePrice

		rows = append(rows, ProfitReportRow{
			BusinessReservationReportRow: businessRows[i],
			PurchasePrice:                purchasePrice,
			PurchasePriceInILS:           purchasePrice * currencyRate,
			Profit:                       profit,
			ProfitInILS:                  profit * currencyRate,
		})
	}

	return rows, nil
}

func calculateCarPurchasePriceWithBrokerERP(reservation db.Reservation) float64 {
	return dbadapters.NumericToFloat64(reservation.PurchasePrice) + dbadapters.NumericToFloat64(reservation.BrokerErpPrice)
}
