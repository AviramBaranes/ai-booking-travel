package markup_rate

import (
	"time"

	"encore.app/services/booking/db"
	"github.com/jackc/pgx/v5/pgtype"
)

var allowedSortFields = map[string]bool{
	"country":       true,
	"broker":        true,
	"mark_up_gross": true,
	"mark_up_net":   true,
	"created_at":    true,
	"updated_at":    true,
}

const limit = 15

type MarkupRateService struct {
	query db.Querier
}

func NewMarkupRateService(query db.Querier) *MarkupRateService {
	return &MarkupRateService{query: query}
}

type MarkupRateResponse struct {
	ID          int64   `json:"id"`
	CountryCode string  `json:"countryCode"`
	Broker      string  `json:"broker"`
	MarkUpGross float64 `json:"markUpGross"`
	MarkUpNet   float64 `json:"markUpNet"`
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

func toMarkupRateResponse(r db.MarkupRate) MarkupRateResponse {
	return MarkupRateResponse{
		ID:          r.ID,
		CountryCode: r.CountryCode,
		Broker:      string(r.Broker),
		MarkUpGross: r.MarkUpGross,
		MarkUpNet:   r.MarkUpNet,
	}
}
