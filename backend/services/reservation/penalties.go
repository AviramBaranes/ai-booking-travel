package reservation

import (
	"context"

	"encore.app/services/reservation/handlers/penalties"
)

func (s *Service) newPenaltiesService() *penalties.PenaltiesService {
	return penalties.NewPenaltiesService(s.query, s.currencyCache)
}

// CreatePenalty records a cancellation fee or no-show fee the supplier charged on a canceled
// reservation, which we in turn charge the customer.
//
// encore:api auth method=POST path=/penalties tag:admin
func (s *Service) CreatePenalty(ctx context.Context, p penalties.CreatePenaltyParams) (*penalties.CreatePenaltyResponse, error) {
	return s.newPenaltiesService().CreatePenalty(ctx, p)
}
