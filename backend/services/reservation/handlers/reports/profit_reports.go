package reports

import (
	"context"

	"encore.app/internal/pricing"
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

func (s *ReportsService) GetProfitReport(ctx context.Context, p ReportParams) (*ProfitReportResponse, error) {
	p.Status = "vouchered" // profit report only includes booked reservations
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
		price := pricing.NewWithReservation(r).Details()

		rows = append(rows, ProfitReportRow{
			BusinessReservationReportRow: businessRows[i],
			PurchasePrice:                price.TotalCost.Value,
			PurchasePriceInILS:           price.TotalCost.ILS,
			Profit:                       price.TotalProfit.Value,
			ProfitInILS:                  price.TotalProfit.ILS,
			ProfitPercentage:             profitPercentage(price),
		})
	}

	return rows, nil
}

// profitPercentage is the share of what the customer actually paid that we kept.
//
// It divides by the total price, the same denominator the report-level ProfitPercentage uses, so a
// row and the footer state the same measure. Dividing by the pre-discount sell price - as this did
// until the price breakdown was centralized - counted a coupon as revenue we never received.
func profitPercentage(price pricing.Details) float64 {
	if price.TotalPrice.Value <= 0 {
		return 0
	}
	return price.TotalProfit.Value / price.TotalPrice.Value * 100
}

// calculateProfit calculates total profit and profit percentage based on the report result
func calculateProfit(result *getReportResult) (totalProfit float64, profitPercentage float64) {
	totalProfit = result.TotalSales - result.TotalCarCost - result.TotalBrokerErpCost
	if result.TotalSales > 0 {
		profitPercentage = (totalProfit / result.TotalSales) * 100
	}

	return totalProfit, profitPercentage
}
