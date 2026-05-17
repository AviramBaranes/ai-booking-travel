package translation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type GetPendingTranslationsResponse struct {
	Translations []db.BrokerTranslation `json:"translations"`
}

type GetPendingTranslationsParams struct {
	Token string `header:"X-Translation-Token" encore:"sensitive"`
}

func (s *TranslationService) GetPendingTranslations(ctx context.Context, p GetPendingTranslationsParams) (*GetPendingTranslationsResponse, error) {
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
