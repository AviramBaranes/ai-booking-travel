package reservation

import (
	"context"

	queries "encore.app/services/reservation/handlers/queries"
)

// encore:api auth method=GET path=/reservations/:id
func (s *Service) GetReservation(ctx context.Context, id int64) (*queries.GetReservationResponse, error) {
	return queries.NewQueryService(s.query).GetReservation(ctx, id)
}

// encore:api auth method=GET path=/reservations
func (s *Service) ListReservations(ctx context.Context, p queries.ListReservationsParams) (*queries.ListReservationsResponse, error) {
	return queries.NewQueryService(s.query).ListReservations(ctx, p)
}

// encore:api private
func (s *Service) GetOpenReservations(ctx context.Context) (*queries.GetOpenReservationsResponse, error) {
	return queries.NewQueryService(s.query).GetOpenReservations(ctx)
}

//encore:api auth method=GET path=/reservations-for-billing tag:accountant
func (s *Service) ListOpenReservationsByBillingEntity(ctx context.Context, p *queries.ListOpenReservationsByBillingEntityParams) (*queries.ListOpenReservationsByBillingEntityResponse, error) {
	return queries.NewQueryService(s.query).ListOpenReservationsByBillingEntity(ctx, p)
}

// encore:api auth method=GET path=/full-reservations/:id tag:admin
func (s *Service) GetFullReservation(ctx context.Context, id int64) (*queries.GetFullReservationResponse, error) {
	return queries.NewQueryService(s.query).GetFullReservation(ctx, id)
}
