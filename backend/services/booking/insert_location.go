package booking

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// InsertLocationParams defines the parameters for inserting a single location.
type InsertLocationParams struct {
	Broker      broker.Name `json:"broker"`
	ID          string      `json:"id"`
	Name        string      `json:"name" validate:"required"`
	Country     string      `json:"country" validate:"required"`
	CountryCode string      `json:"country_code" validate:"required,len=2"`
	City        string      `json:"city" validate:"omitempty"`
	Iata        string      `json:"iata" validate:"omitempty,len=3"`
}

var (
	ErrInvalidBroker = api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code:  api_errors.CodeInvalidValue,
		Field: "broker",
	})
)

// Validate validates the InsertLocationParams struct.
func (p InsertLocationParams) Validate() error {
	_, err := toDbBroker(p.Broker)
	if err != nil {
		return ErrInvalidBroker
	}
	return validation.ValidateStruct(p)
}

// InsertLocations handles the HTTP request to insert a single location
// encore:api auth method=POST path=/locations tag:admin
func (s *Service) InsertLocation(ctx context.Context, p InsertLocationParams) error {
	loc := broker.Location{
		ID:          p.ID,
		Name:        p.Name,
		Country:     p.Country,
		CountryCode: p.CountryCode,
		City:        p.City,
		Iata:        p.Iata,
	}

	err := insertBatch(ctx, s.query, []broker.Location{loc}, p.Broker)
	if err != nil {
		rlog.Error("failed to insert location", "broker", p.Broker, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
