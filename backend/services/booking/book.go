package booking

import (
	"context"

	"encore.app/services/booking/handlers/booking_handlers"
)

//encore:api auth method=POST path=/booking tag:agent
func (s *Service) Book(ctx context.Context, p booking_handlers.BookParams) (*booking_handlers.BookResponse, error) {
	bs := booking_handlers.NewBookingService(s.query)
	return bs.Book(ctx, p)
}

//encore:api auth method=POST path=/offers/book tag:agent
func (s *Service) BookPriceOffer(ctx context.Context, p booking_handlers.BookPriceOfferParams) (*booking_handlers.BookResponse, error) {
	bs := booking_handlers.NewBookingService(s.query)
	return bs.BookPriceOffer(ctx, p)
}
