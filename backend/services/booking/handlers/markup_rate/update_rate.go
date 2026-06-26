package markup_rate

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type UpdateMarkupRateParams struct {
	CountryCode string  `json:"countryCode" validate:"required"`
	Broker      string  `json:"broker" validate:"required,oneof=hertz flex"`
	MarkUpGross float64 `json:"markUpGross" validate:"required"`
	MarkUpNet   float64 `json:"markUpNet" validate:"required"`
}

func (p UpdateMarkupRateParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *MarkupRateService) UpdateMarkupRate(ctx context.Context, id int64, p UpdateMarkupRateParams) (*MarkupRateResponse, error) {
	row, err := s.query.UpdateMarkupRate(ctx, db.UpdateMarkupRateParams{
		ID:          id,
		CountryCode: p.CountryCode,
		Broker:      db.Broker(p.Broker),
		MarkUpGross: p.MarkUpGross,
		MarkUpNet:   p.MarkUpNet,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update markup rate", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toMarkupRateResponse(row)
	return &resp, nil
}
