package translation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type TranslateTranslationParams struct {
	Token      string `header:"X-Translation-Token" encore:"sensitive"`
	ID         int64  `json:"id"`
	TargetText string `json:"targetText"`
	Confidence int32  `json:"confidence"`
}

func (s *TranslationService) TranslateTranslation(ctx context.Context, p TranslateTranslationParams) error {
	if p.Token != secrets.translationToken {
		rlog.Warn("invalid translation token", "provided_token", p.Token)
		return api_errors.ErrNotFound
	}

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
