package markup_rate

import (
	"time"

	"encore.app/services/booking/db"
	"github.com/jackc/pgx/v5/pgtype"
)

var allowedSortFields = map[string]bool{
	"country":                 true,
	"brand":                   true,
	"car_group":               true,
	"pickup_date_from":        true,
	"num_of_rental_days_from": true,
}

const limit = 15

type HertzMarkupRateService struct {
	query db.Querier
}

func NewHertzMarkupRateService(query db.Querier) *HertzMarkupRateService {
	return &HertzMarkupRateService{query: query}
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
