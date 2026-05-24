package queries

import "encore.app/services/reservation/db"

type QueryService struct {
	query db.Querier
}

func NewQueryService(query db.Querier) *QueryService {
	return &QueryService{query: query}
}
