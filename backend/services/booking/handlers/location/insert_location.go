package location

import (
	"context"
	"slices"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type InsertLocationParams struct {
	Broker      broker.Name `json:"broker"`
	ID          string      `json:"id"`
	Name        string      `json:"name" validate:"required"`
	Country     string      `json:"country" validate:"required"`
	CountryCode string      `json:"country_code" validate:"required,len=2"`
	City        string      `json:"city" validate:"omitempty"`
	Iata        string      `json:"iata" validate:"omitempty,len=3"`
	// Brands lists the brands the station serves, required for Hertz GR and rejected for brokers
	// without brands.
	Brands []string `json:"brands,omitempty" encore:"optional"`
}

var ErrInvalidBroker = api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
	Code:  api_errors.CodeInvalidValue,
	Field: "broker",
})

var ErrInvalidBrands = api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
	Code:  api_errors.CodeInvalidValue,
	Field: "brands",
})

func (p InsertLocationParams) Validate() error {
	_, err := toDbBroker(p.Broker)
	if err != nil {
		return ErrInvalidBroker
	}
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}
	return validateBrands(p.Broker, p.Brands)
}

// validateBrands checks a station's brands against its broker: a Hertz GR station must name the
// brands it serves, and a station of any other broker must name none.
func validateBrands(b broker.Name, brands []string) error {
	if b != broker.BrokerHertzGR {
		if len(brands) > 0 {
			return ErrInvalidBrands
		}
		return nil
	}

	if len(brands) == 0 {
		return ErrInvalidBrands
	}
	for i, brand := range brands {
		if !slices.Contains(broker.HertzGRBrands, brand) || slices.Contains(brands[:i], brand) {
			return ErrInvalidBrands
		}
	}
	return nil
}

func (s *LocationService) InsertLocation(ctx context.Context, p InsertLocationParams) error {
	loc := broker.Location{
		ID:          p.ID,
		Name:        p.Name,
		Country:     p.Country,
		CountryCode: p.CountryCode,
		City:        p.City,
		Iata:        p.Iata,
		Brands:      p.Brands,
	}

	err := insertBatch(ctx, s.query, []broker.Location{loc}, p.Broker)
	if err != nil {
		rlog.Error("failed to insert location", "broker", p.Broker, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
