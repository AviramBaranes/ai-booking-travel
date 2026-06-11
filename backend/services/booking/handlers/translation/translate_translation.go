package translation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type TranslateTranslationParams struct {
	ID         int64  `json:"id"`
	TargetText string `json:"targetText"`
	Confidence int32  `json:"confidence"`
}

func (s *TranslationService) TranslateTranslation(ctx context.Context, p TranslateTranslationParams) error {
	err := s.query.TranslatePendingTranslation(ctx, db.TranslatePendingTranslationParams{
		ID:              p.ID,
		TargetText:      &p.TargetText,
		ConfidenceScore: &p.Confidence,
	})

	if err != nil {
		rlog.Error("failed to translate translation", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
