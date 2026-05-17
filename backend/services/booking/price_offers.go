package booking

import (
	"context"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/booking/db"
	availability "encore.app/services/booking/handlers/availability_handlers"
	poh "encore.app/services/booking/handlers/price_offer_handlers"
)

// CreatePriceOffer creates a new price offer based on the provided parameters.
//
// encore:api auth method=POST path=/booking/price-offers tag:agent
func (s *Service) CreatePriceOffer(ctx context.Context, p poh.CreatePriceOfferParams) (*poh.PriceOfferResponse, error) {
	pos := poh.NewPriceOfferService(s.query, nil)
	return pos.CreatePriceOffer(ctx, p)
}

// RenewPriceOffer refreshes the stored pricing details for a price offer if the original plan is still available.
//
// encore:api auth method=POST path=/booking/price-offers/:id/renew tag:agent
func (s *Service) RenewPriceOffer(ctx context.Context, id int64) (*poh.RenewPriceOfferResponse, error) {
	searchFn := func(ctx context.Context, pickupLocID, dropoffLocID int64, pickupDate, dropoffDate, pickupTime, dropoffTime string, driverAge int) (int64, error) {
		resp, err := SearchAvailability(ctx, availability.SearchAvailabilityParams{
			PickupLocationID:  pickupLocID,
			DropoffLocationID: dropoffLocID,
			PickupDate:        pickupDate,
			DropoffDate:       dropoffDate,
			PickupTime:        pickupTime,
			DropoffTime:       dropoffTime,
			DriverAge:         driverAge,
		})
		if err != nil || resp == nil {
			return 0, err
		}
		return resp.SnapshotID, nil
	}
	pos := poh.NewPriceOfferService(s.query, searchFn)
	return pos.RenewPriceOffer(ctx, id)
}

// GetClientPriceOffer retrieves the details of a price offer by token (public, no internal pricing).
//
// encore:api public method=GET path=/booking/price-offers/client/:token
func (s *Service) GetClientPriceOffer(ctx context.Context, token string) (*poh.GetPriceOfferResponse, error) {
	pos := poh.NewPriceOfferService(s.query, nil)
	return pos.GetClientPriceOffer(ctx, token)
}

// GetAgentPriceOffer retrieves the details of a price offer for the authenticated agent, including internal pricing.
//
// encore:api auth method=GET path=/booking/price-offers/agent/:id tag:agent
func (s *Service) GetAgentPriceOffer(ctx context.Context, id int64) (*poh.GetAgentPriceOfferResponse, error) {
	pos := poh.NewPriceOfferService(s.query, nil)
	return pos.GetAgentPriceOffer(ctx, id)
}

// ListPriceOffers returns a paginated list of the authenticated agent's price offers.
//
// encore:api auth method=GET path=/booking/price-offers tag:agent
func (s *Service) ListPriceOffers(ctx context.Context, p poh.ListPriceOffersRequest) (*poh.ListPriceOffersResponse, error) {
	pos := poh.NewPriceOfferService(s.query, nil)
	return pos.ListPriceOffers(ctx, p)
}

// UpdatePriceOffer updates a price offer's mutable fields for the authenticated agent.
//
// encore:api auth method=PATCH path=/booking/price-offers/:id tag:agent
func (s *Service) UpdatePriceOffer(ctx context.Context, id int64, p poh.UpdatePriceOfferParams) error {
	pos := poh.NewPriceOfferService(s.query, nil)
	return pos.UpdatePriceOffer(ctx, id, p)
}

func isPriceOfferErpIncluded(offer db.GetPriceOfferByIdRow) bool {
	return offer.BtErpPrice != 0 || dbadapters.NumericToFloat64(offer.BrokerErpPrice) != 0
}
