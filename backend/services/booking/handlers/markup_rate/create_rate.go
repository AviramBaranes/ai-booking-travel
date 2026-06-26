package markup_rate

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type CreateMarkupRateParams struct {
	CountryCode string  `json:"countryCode" validate:"required"`
	Broker      string  `json:"broker" validate:"required,oneof=hertz flex"`
	MarkUpGross float64 `json:"markUpGross" validate:"required"`
	MarkUpNet   float64 `json:"markUpNet" validate:"required"`
}

func (p CreateMarkupRateParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *MarkupRateService) CreateMarkupRate(ctx context.Context, p CreateMarkupRateParams) (*MarkupRateResponse, error) {
	row, err := s.query.InsertMarkupRate(ctx, db.InsertMarkupRateParams{
		Broker:      db.Broker(p.Broker),
		CountryCode: p.CountryCode,
		MarkUpGross: p.MarkUpGross,
		MarkUpNet:   p.MarkUpNet,
	})
	if err != nil {
		rlog.Error("failed to create  markup rate", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toMarkupRateResponse(row)
	return &resp, nil
}
