package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type SearchLocationParams struct {
	Search string `query:"search" validate:"required,min=3"`
}

func (p SearchLocationParams) Validate() error {
	return validation.ValidateStruct(p)
}

type SearchLocationResponse struct {
	Locations []LocationResult `json:"locations"`
}

type LocationResult struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        *string `json:"city,omitempty"`
	Iata        *string `json:"iata,omitempty"`
}

// SearchLocations searches for locations in the database that match the given search query
// encore:api public method=GET path=/locations/search
func (s *Service) SearchLocations(ctx context.Context, params SearchLocationParams) (*SearchLocationResponse, error) {
	rlog.Info("searching for locations matching query", "search", params.Search)
	locs, err := s.query.SearchLocations(ctx, params.Search)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &SearchLocationResponse{Locations: []LocationResult{}}, nil
		}
		rlog.Error("failed to search locations in database", "error", err)
		return nil, api_errors.ErrInternalError
	}

	results := make([]LocationResult, 0, len(locs))

	for _, loc := range locs {
		results = append(results, LocationResult{
			ID:          loc.ID,
			Name:        loc.Name,
			Country:     loc.Country,
			CountryCode: loc.CountryCode,
			City:        loc.City,
			Iata:        loc.Iata,
		})
	}

	return &SearchLocationResponse{
		Locations: results,
	}, nil
}
