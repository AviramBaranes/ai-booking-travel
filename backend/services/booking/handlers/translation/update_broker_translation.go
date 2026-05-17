package translation

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type UpdateBrokerTranslationParams struct {
	TargetText string `json:"target_text" validate:"required,notblank"`
}

func (p UpdateBrokerTranslationParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *TranslationService) UpdateBrokerTranslation(ctx context.Context, id int64, p UpdateBrokerTranslationParams) error {
	err := s.query.UpdateBrokerTranslation(ctx, db.UpdateBrokerTranslationParams{
		ID:         id,
		TargetText: &p.TargetText,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to update broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
