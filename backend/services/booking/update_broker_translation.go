package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type UpdateBrokerTranslationRequest struct {
	TargetText string `json:"target_text" validate:"required,notblank"`
}

func (p UpdateBrokerTranslationRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdateBrokerTranslation updates a broker translation target text by ID.
//
//encore:api auth method=PUT path=/broker-translations/:id tag:admin
func (s *Service) UpdateBrokerTranslation(ctx context.Context, id int32, params UpdateBrokerTranslationRequest) error {
	err := s.query.UpdateBrokerTranslation(ctx, db.UpdateBrokerTranslationParams{
		ID:         id,
		TargetText: &params.TargetText,
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
