package translation

import (
	"encore.app/services/booking/db"
)

type TranslationService struct {
	query db.Querier
}

func NewTranslationService(query db.Querier) *TranslationService {
	return &TranslationService{query: query}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
