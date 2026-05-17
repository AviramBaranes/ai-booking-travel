package markup_rate

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type CreateHertzMarkupRateParams struct {
	Country             string  `json:"country" validate:"required"`
	Brand               string  `json:"brand" validate:"required"`
	PickupDateFrom      string  `json:"pickupDateFrom" validate:"required,datetime=2006-01-02"`
	PickupDateTo        string  `json:"pickupDateTo" validate:"required,datetime=2006-01-02"`
	CarGroup            string  `json:"carGroup" validate:"required"`
	NumOfRentalDaysFrom int32   `json:"numOfRentalDaysFrom" validate:"required,gte=1"`
	NumOfRentalDaysTo   int32   `json:"numOfRentalDaysTo" validate:"required,gte=1,gtefield=NumOfRentalDaysFrom"`
	MarkUpGross         float64 `json:"markUpGross" validate:"required"`
	MarkUpNet           float64 `json:"markUpNet" validate:"required"`
}

func (p CreateHertzMarkupRateParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *HertzMarkupRateService) CreateHertzMarkupRate(ctx context.Context, p CreateHertzMarkupRateParams) (*HertzMarkupRateResponse, error) {
	row, err := s.query.InsertHertzMarkupRate(ctx, db.InsertHertzMarkupRateParams{
		Country:             p.Country,
		Brand:               p.Brand,
		PickupDateFrom:      parseDate(p.PickupDateFrom),
		PickupDateTo:        parseDate(p.PickupDateTo),
		CarGroup:            p.CarGroup,
		NumOfRentalDaysFrom: p.NumOfRentalDaysFrom,
		NumOfRentalDaysTo:   p.NumOfRentalDaysTo,
		MarkUpGross:         p.MarkUpGross,
		MarkUpNet:           p.MarkUpNet,
	})
	if err != nil {
		rlog.Error("failed to create hertz markup rate", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toHertzMarkupRateResponse(row)
	return &resp, nil
}
