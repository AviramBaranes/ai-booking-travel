package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type ListLocationsResponse struct {
	Locations []LocationRow `json:"locations"`
	Total     int64         `json:"total"`
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
	CountryCode string `query:"country_code"`
	Broker      string `query:"broker"`
	Name        string `query:"name"`
	Iata        string `query:"iata"`
	Enabled     string `query:"enabled"`
	Page        int    `query:"page" validate:"required,min=1"`
}

func (l ListLocationsRequest) Validate() error {
	if l.Enabled != "" && l.Enabled != "true" && l.Enabled != "false" {
		return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code: api_errors.CodeInvalidValue, Field: "enabled",
		})
	}
	return validation.ValidateStruct(l)
}

const LocationsLimit = 15

//encore:api auth method=GET path=/locations tag:admin
func (s *Service) ListLocations(ctx context.Context, p ListLocationsRequest) (*ListLocationsResponse, error) {
	var enabled *bool
	if p.Enabled != "" {
		v := p.Enabled == "true"
		enabled = &v
	}

	filterParams := db.CountLocationBrokerCodesWithLocationParams{
		CountryCode: nilIfEmpty(p.CountryCode),
		Broker:      nilIfEmpty(p.Broker),
		Name:        nilIfEmpty(p.Name),
		Iata:        nilIfEmpty(p.Iata),
		Enabled:     enabled,
	}

	total, err := s.query.CountLocationBrokerCodesWithLocation(ctx, filterParams)
	if err != nil {
		rlog.Error("failed to count locations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	offset := (p.Page - 1) * LocationsLimit
	rows, err := s.query.ListLocationBrokerCodesWithLocation(ctx, db.ListLocationBrokerCodesWithLocationParams{
		Limit:       LocationsLimit,
		Offset:      int32(offset),
		CountryCode: filterParams.CountryCode,
		Broker:      filterParams.Broker,
		Name:        filterParams.Name,
		Iata:        filterParams.Iata,
		Enabled:     filterParams.Enabled,
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
		Total:     total,
	}, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
