package booking

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgtype"
)

// --- Request / Response types ---

var allowedSortFields = map[string]bool{
	"country":                 true,
	"brand":                   true,
	"car_group":               true,
	"pickup_date_from":        true,
	"num_of_rental_days_from": true,
}

type HertzMarkupRateResponse struct {
	ID                  int64   `json:"id"`
	Country             string  `json:"country"`
	Brand               string  `json:"brand"`
	PickupDateFrom      string  `json:"pickupDateFrom"`
	PickupDateTo        string  `json:"pickupDateTo"`
	CarGroup            string  `json:"carGroup"`
	NumOfRentalDaysFrom int     `json:"numOfRentalDaysFrom"`
	NumOfRentalDaysTo   int     `json:"numOfRentalDaysTo"`
	MarkUpGross         float64 `json:"markUpGross"`
	MarkUpNet           float64 `json:"markUpNet"`
}

type ListHertzMarkupRatesRequest struct {
	Country  string `query:"country" validate:"omitempty"`
	Brand    string `query:"brand" validate:"omitempty"`
	CarGroup string `query:"carGroup" validate:"omitempty"`
	SortBy   string `query:"sortBy" validate:"omitempty"`
	SortDir  string `query:"sortDir" validate:"omitempty,oneof=asc desc"`
	Page     int32  `query:"page" validate:"required,gte=1"`
	Limit    int32  `query:"limit" validate:"required,gte=1,lte=100"`
}

func (p ListHertzMarkupRatesRequest) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}
	if p.SortBy != "" && !allowedSortFields[p.SortBy] {
		return api_errors.ErrInvalidValue
	}
	return nil
}

type ListHertzMarkupRatesResponse struct {
	Rates []HertzMarkupRateResponse `json:"rates"`
}

type CreateHertzMarkupRateRequest struct {
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

func (p CreateHertzMarkupRateRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateHertzMarkupRateRequest struct {
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

func (p UpdateHertzMarkupRateRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

func parseDate(s string) pgtype.Date {
	t, _ := time.Parse("2006-01-02", s)
	return pgtype.Date{Time: t, Valid: true}
}

func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func formatDate(d pgtype.Date) string {
	return d.Time.Format("2006-01-02")
}

func toHertzMarkupRateResponse(r db.HertzMarkupRate) HertzMarkupRateResponse {
	return HertzMarkupRateResponse{
		ID:                  r.ID,
		Country:             r.Country,
		Brand:               r.Brand,
		PickupDateFrom:      formatDate(r.PickupDateFrom),
		PickupDateTo:        formatDate(r.PickupDateTo),
		CarGroup:            r.CarGroup,
		NumOfRentalDaysFrom: int(r.NumOfRentalDaysFrom),
		NumOfRentalDaysTo:   int(r.NumOfRentalDaysTo),
		MarkUpGross:         r.MarkUpGross,
		MarkUpNet:           r.MarkUpNet,
	}
}

// --- Endpoints ---

// ListHertzMarkupRates lists hertz markup rates with pagination, optional filtering, and sorting.
//
//encore:api auth method=GET path=/hertz-markup-rates tag:admin
func (s *Service) ListHertzMarkupRates(ctx context.Context, params ListHertzMarkupRatesRequest) (*ListHertzMarkupRatesResponse, error) {
	offset := (params.Page - 1) * params.Limit

	sortField := params.SortBy
	sortDir := params.SortDir
	if sortField == "" {
		sortField = "country"
		sortDir = "asc"
	}
	if sortDir == "" {
		sortDir = "asc"
	}

	rows, err := s.query.ListHertzMarkupRates(ctx, db.ListHertzMarkupRatesParams{
		Country:     toStringPtr(params.Country),
		Brand:       toStringPtr(params.Brand),
		CarGroup:    toStringPtr(params.CarGroup),
		SortField:   sortField,
		SortDir:     sortDir,
		QueryOffset: offset,
		QueryLimit:  params.Limit,
	})
	if err != nil {
		rlog.Error("failed to list hertz markup rates", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rates := make([]HertzMarkupRateResponse, 0, len(rows))
	for _, r := range rows {
		rates = append(rates, toHertzMarkupRateResponse(r))
	}

	return &ListHertzMarkupRatesResponse{Rates: rates}, nil
}

// CreateHertzMarkupRate creates a new hertz markup rate.
//
//encore:api auth method=POST path=/hertz-markup-rates tag:admin
func (s *Service) CreateHertzMarkupRate(ctx context.Context, params CreateHertzMarkupRateRequest) (*HertzMarkupRateResponse, error) {
	row, err := s.query.InsertHertzMarkupRate(ctx, db.InsertHertzMarkupRateParams{
		Country:             params.Country,
		Brand:               params.Brand,
		PickupDateFrom:      parseDate(params.PickupDateFrom),
		PickupDateTo:        parseDate(params.PickupDateTo),
		CarGroup:            params.CarGroup,
		NumOfRentalDaysFrom: params.NumOfRentalDaysFrom,
		NumOfRentalDaysTo:   params.NumOfRentalDaysTo,
		MarkUpGross:         params.MarkUpGross,
		MarkUpNet:           params.MarkUpNet,
	})
	if err != nil {
		rlog.Error("failed to create hertz markup rate", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toHertzMarkupRateResponse(row)
	return &resp, nil
}

// UpdateHertzMarkupRate updates an existing hertz markup rate.
//
//encore:api auth method=PUT path=/hertz-markup-rates/:id tag:admin
func (s *Service) UpdateHertzMarkupRate(ctx context.Context, id int64, params UpdateHertzMarkupRateRequest) (*HertzMarkupRateResponse, error) {
	row, err := s.query.UpdateHertzMarkupRate(ctx, db.UpdateHertzMarkupRateParams{
		ID:                  id,
		Country:             params.Country,
		Brand:               params.Brand,
		PickupDateFrom:      parseDate(params.PickupDateFrom),
		PickupDateTo:        parseDate(params.PickupDateTo),
		CarGroup:            params.CarGroup,
		NumOfRentalDaysFrom: params.NumOfRentalDaysFrom,
		NumOfRentalDaysTo:   params.NumOfRentalDaysTo,
		MarkUpGross:         params.MarkUpGross,
		MarkUpNet:           params.MarkUpNet,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update hertz markup rate", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toHertzMarkupRateResponse(row)
	return &resp, nil
}

// DeleteHertzMarkupRate deletes a hertz markup rate by its ID.
//
//encore:api auth method=DELETE path=/hertz-markup-rates/:id tag:admin
func (s *Service) DeleteHertzMarkupRate(ctx context.Context, id int64) error {
	_, err := s.query.DeleteHertzMarkupRate(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete hertz markup rate", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
