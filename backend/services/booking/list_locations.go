package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ListLocationsResponse struct {
	Locations []LocationRow `json:"locations"`
}

type LocationRow struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	CountryCode      string  `json:"country_code"`
	Country          string  `json:"country"`
	City             *string `json:"city"`
	Iata             *string `json:"iata"`
	Enabled          bool    `json:"enabled"`
	BrokerLocationID string  `json:"broker_location_id"`
}

type ListLocationsRequest struct {
	Search string `json:"search"`
	Page   int    `json:"page" validate:"required,min=1"`
}

func (l ListLocationsRequest) Validate() error {
	return validation.ValidateStruct(l)
}

const LocationsLimit = 15

//encore:api auth method=GET path=/locations tag:admin
func (s *Service) ListLocations(ctx context.Context, p ListLocationsRequest) (*ListLocationsResponse, error) {
	offset := (p.Page - 1) * LocationsLimit
	rows, err := s.query.ListLocationBrokerCodesWithLocation(ctx, db.ListLocationBrokerCodesWithLocationParams{
		Limit:  LocationsLimit,
		Offset: int32(offset),
		Search: &p.Search,
	})

	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListLocationsResponse{
				Locations: []LocationRow{},
			}, nil
		}
		rlog.Error("failed to list locations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	locations := make([]LocationRow, len(rows))
	for i, row := range rows {
		locations[i] = LocationRow{
			ID:               row.ID,
			Name:             row.LocationName,
			CountryCode:      row.LocationCountryCode,
			Country:          row.LocationCountry,
			City:             row.LocationCity,
			Iata:             row.LocationIata,
			Enabled:          row.Enabled,
			BrokerLocationID: row.BrokerLocationID,
		}
	}

	return &ListLocationsResponse{
		Locations: locations,
	}, nil
}
