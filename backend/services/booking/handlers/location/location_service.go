package location

import (
	"encore.app/services/booking/db"
)

type LocationService struct {
	query db.Querier
}

func NewLocationService(query db.Querier) *LocationService {
	return &LocationService{query: query}
}
