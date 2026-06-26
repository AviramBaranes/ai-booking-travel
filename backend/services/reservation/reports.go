package reservation

import (
	"context"

	"encore.app/services/reservation/handlers/reports"
)

// encore:api auth tag:admin method=GET path=/reports/business
func (s *Service) GetBusinessReport(ctx context.Context, p reports.ReportParams) (*reports.BusinessReportResponse, error) {
	rs := reports.NewReportsService(s.query)
	return rs.GetBusinessReport(ctx, p)
}

// encore:api auth tag:admin method=GET path=/reports/profit
func (s *Service) GetProfitReport(ctx context.Context, p reports.ReportParams) (*reports.ProfitReportResponse, error) {
	rs := reports.NewReportsService(s.query)
	return rs.GetProfitReport(ctx, p)
}

// encore:api auth tag:admin method=GET path=/reports/businesses-balances
func (s *Service) GetBusinessesBalancesReport(ctx context.Context) (*reports.BusinessesBalancesReportResponse, error) {
	rs := reports.NewReportsService(s.query)
	return rs.GetBusinessesBalancesReport(ctx)
}
