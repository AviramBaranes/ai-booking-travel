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
	ProfitPercentage   float64 `json:"profitPercentage"`
}

type ProfitReportResponse struct {
	Reservations     []ProfitReportRow `json:"reservations"`
	Count            int64             `json:"count"`
	TotalSales       float64           `json:"totalSales"`
	TotalProfit      float64           `json:"totalProfit"`
	ProfitPercentage float64           `json:"profitPercentage"`
}

// encore:api auth tag:admin method=GET path=/reports/profit
func (s *Service) GetProfitReport(ctx context.Context, p ReportParams) (*ProfitReportResponse, error) {
	p.Status = "vouchered" // profit report only includes booked reservations
	result, err := s.getReports(ctx, p, true)
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
		rlog.Error("failed to get accounts lookup for profit report", "error", err)
		return nil, err
	}

	reservations, err := buildProfitReportRows(result.Reservations, accountsLookup)
	if err != nil {
		return nil, err
	}

	totalProfit, profitPercentage := calculateProfit(result)

	return &ProfitReportResponse{
		Reservations:     reservations,
		Count:            result.Count,
		TotalSales:       result.TotalSales,
		TotalProfit:      totalProfit,
		ProfitPercentage: profitPercentage,
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
		btErpPrice := dbadapters.NumericToFloat64(r.BtErpPrice)
		purchasePrice := calculateCarPurchasePriceWithBrokerERP(r)
		profit := businessRows[i].CarSellPriceWithBrokerERP - purchasePrice + btErpPrice

		rows = append(rows, ProfitReportRow{
			BusinessReservationReportRow: businessRows[i],
			PurchasePrice:                purchasePrice,
			PurchasePriceInILS:           purchasePrice * currencyRate,
			Profit:                       profit,
			ProfitInILS:                  profit * currencyRate,
			ProfitPercentage:             (profit / businessRows[i].CarSellPriceWithBrokerERP) * 100,
		})
	}

	return rows, nil
}

func calculateCarPurchasePriceWithBrokerERP(reservation db.Reservation) float64 {
	return dbadapters.NumericToFloat64(reservation.PurchasePrice) + dbadapters.NumericToFloat64(reservation.BrokerErpPrice)
}

// calculateProfit calculates total profit and profit percentage based on the report result
func calculateProfit(result *getReportResult) (totalProfit float64, profitPercentage float64) {
	totalProfit = result.TotalSales - result.TotalCarCost - result.TotalBrokerErpCost
	if result.TotalSales > 0 {
		profitPercentage = (totalProfit / result.TotalSales) * 100
	}

	return totalProfit, profitPercentage
}
