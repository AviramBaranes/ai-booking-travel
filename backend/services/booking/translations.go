package booking

import (
	"context"

	"encore.app/services/booking/handlers/translation_handlers"
)

// GetPendingTranslations returns the list of pending translations for brokers. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/pending method=GET
func (s *Service) GetPendingTranslations(ctx context.Context, p translation_handlers.GetPendingTranslationsParams) (*translation_handlers.GetPendingTranslationsResponse, error) {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.GetPendingTranslations(ctx, p)
}

// TranslateTranslation translates a pending translation. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/translate method=PATCH
func (s *Service) TranslateTranslation(ctx context.Context, p translation_handlers.TranslateTranslationParams) error {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.TranslateTranslation(ctx, p)
}

//encore:api auth method=GET path=/broker-translations tag:admin
func (s *Service) ListBrokerTranslations(ctx context.Context, p translation_handlers.ListBrokerTranslationsParams) (*translation_handlers.ListBrokerTranslationsResponse, error) {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.ListBrokerTranslations(ctx, p)
}

// UpdateBrokerTranslation updates a broker translation target text by ID.
//
//encore:api auth method=PUT path=/broker-translations/:id tag:admin
func (s *Service) UpdateBrokerTranslation(ctx context.Context, id int64, p translation_handlers.UpdateBrokerTranslationParams) error {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.UpdateBrokerTranslation(ctx, id, p)
}

// VerifyBrokerTranslation marks a broker translation as verified by ID.
//
//encore:api auth method=PATCH path=/broker-translations/:id/verify tag:admin
func (s *Service) VerifyBrokerTranslation(ctx context.Context, id int64) error {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.VerifyBrokerTranslation(ctx, id)
}

// DeleteBrokerTranslation deletes a broker translation by ID.
//
//encore:api auth method=DELETE path=/broker-translations/:id tag:admin
func (s *Service) DeleteBrokerTranslation(ctx context.Context, id int64) error {
	ts := translation_handlers.NewTranslationService(s.query)
	return ts.DeleteBrokerTranslation(ctx, id)
}
