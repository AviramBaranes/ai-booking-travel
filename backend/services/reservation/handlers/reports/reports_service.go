package reports

import "encore.app/services/reservation/db"

type ReportsService struct {
	query db.Querier
}

func NewReportsService(query db.Querier) *ReportsService {
	return &ReportsService{query: query}
}
