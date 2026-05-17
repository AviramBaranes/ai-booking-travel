package booking

import (
	"context"
	"net/http"

	"encore.app/services/booking/handlers/location"
)

// ListLocations lists location broker codes with optional filters.
//
//encore:api auth method=GET path=/locations tag:admin
func (s *Service) ListLocations(ctx context.Context, p location.ListLocationsParams) (*location.ListLocationsResponse, error) {
	ls := location.NewLocationService(s.query)
	return ls.ListLocations(ctx, p)
}

// SearchLocations searches for locations matching the given query.
//
// encore:api public method=GET path=/locations/search
func (s *Service) SearchLocations(ctx context.Context, p location.SearchLocationParams) (*location.SearchLocationResponse, error) {
	ls := location.NewLocationService(s.query)
	return ls.SearchLocations(ctx, p)
}

// InsertLocation inserts a single location broker code into the database.
//
// encore:api auth method=POST path=/locations tag:admin
func (s *Service) InsertLocation(ctx context.Context, p location.InsertLocationParams) error {
	ls := location.NewLocationService(s.query)
	return ls.InsertLocation(ctx, p)
}

// InsertFlexLocations fetches all Flex locations from the broker and upserts them.
//
// encore:api auth method=POST path=/locations/flex tag:admin
func (s *Service) InsertFlexLocations(ctx context.Context) error {
	ls := location.NewLocationService(s.query)
	return ls.InsertFlexLocations(ctx)
}

// InsertHertzLocations reads a CSV file upload and upserts Hertz locations.
//
//encore:api auth method=POST path=/locations/hertz tag:admin raw
func (s *Service) InsertHertzLocations(w http.ResponseWriter, req *http.Request) {
	ls := location.NewLocationService(s.query)
	ls.InsertHertzLocations(w, req)
}

// ToggleLocation enables or disables a location broker code by ID.
//
//encore:api auth method=PATCH path=/locations/:id tag:admin
func (s *Service) ToggleLocation(ctx context.Context, id int64, p location.ToggleLocationParams) error {
	ls := location.NewLocationService(s.query)
	return ls.ToggleLocation(ctx, id, p)
}

// BulkToggleLocations enables or disables multiple location broker codes.
//
//encore:api auth method=PATCH path=/location-bulk-toggle tag:admin
func (s *Service) BulkToggleLocations(ctx context.Context, p location.BulkToggleLocationsParams) error {
	ls := location.NewLocationService(s.query)
	return ls.BulkToggleLocations(ctx, p)
}

// DeleteLocation deletes a location broker code by ID. If no other broker codes
// reference the same location, the location is also deleted.
//
//encore:api auth method=DELETE path=/locations/:id tag:admin
func (s *Service) DeleteLocation(ctx context.Context, id int64) error {
	ls := location.NewLocationService(s.query)
	return ls.DeleteLocation(ctx, id)
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
