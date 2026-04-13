package booking

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

var secrets struct {
	translationToken string
}

// GetPendingTranslationsResponse is the response for GetPendingTranslations endpoint
type GetPendingTranslationsResponse struct {
	Translations []db.BrokerTranslation `json:"translations"`
}

// GetPendingTranslationsRequest is the request for GetPendingTranslations endpoint
type GetPendingTranslationsRequest struct {
	Token string `header:"X-Translation-Token" encore:"sensitive"`
}

// GetPendingTranslations returns the list of pending translations for brokers. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/pending method=GET
func (s *Service) GetPendingTranslations(ctx context.Context, p *GetPendingTranslationsRequest) (*GetPendingTranslationsResponse, error) {
	if p.Token != secrets.translationToken {
		rlog.Warn("invalid translation token", "provided_token", p.Token)
		return nil, api_errors.ErrNotFound
	}

	ts, err := s.query.ListPendingTranslations(ctx)
	if err != nil {
		rlog.Error("failed to get pending translations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetPendingTranslationsResponse{
		Translations: ts,
	}, nil
}

// TranslateTranslationRequest is the request for TranslateTranslation endpoint
type TranslateTranslationRequest struct {
	Token      string `header:"X-Translation-Token" encore:"sensitive"`
	ID         int32  `json:"id"`
	TargetText string `json:"targetText"`
	Confidence int32  `json:"confidence"`
}

// TranslateTranslation translates a pending translation. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/translate method=PATCH
func (s *Service) TranslateTranslation(ctx context.Context, p *TranslateTranslationRequest) error {
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
