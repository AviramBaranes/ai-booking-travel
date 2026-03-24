package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type GetLocationByBrokerIDParams struct {
	BrokerLocationID string `query:"broker_location_id" validate:"required,notblank"`
}

func (p GetLocationByBrokerIDParams) Validate() error {
	return validation.ValidateStruct(p)
}

type GetLocationResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        *string `json:"city"`
	Iata        *string `json:"iata"`
}

// encore:api private method=GET path=/locations/by-broker-id
func (s *Service) GetLocationByBrokerLocationID(ctx context.Context, params GetLocationByBrokerIDParams) (*GetLocationResponse, error) {
	loc, err := s.query.GetLocationByBrokerLocationID(ctx, params.BrokerLocationID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get location by broker location ID", "brokerLocationID", params.BrokerLocationID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetLocationResponse{
		ID:          loc.ID,
		Name:        loc.Name,
		Country:     loc.Country,
		CountryCode: loc.CountryCode,
		City:        loc.City,
		Iata:        loc.Iata,
	}, nil
}
