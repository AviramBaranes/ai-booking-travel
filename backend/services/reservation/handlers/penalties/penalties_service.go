package penalties

import (
	"encore.app/internal/currency"
	"encore.app/services/reservation/db"
)

// PenaltiesService provides the operations on cancellation and no-show penalties.
type PenaltiesService struct {
	query         db.Querier
	currencyCache *currency.CurrenciesCache
}

// NewPenaltiesService creates a new PenaltiesService with the given dependencies.
func NewPenaltiesService(query db.Querier, currencyCache *currency.CurrenciesCache) *PenaltiesService {
	return &PenaltiesService{
		query:         query,
		currencyCache: currencyCache,
	}
}
