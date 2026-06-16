package booking

import (
	"context"

	"encore.app/services/booking/handlers/translation"
)

// GetPendingTranslations returns the list of pending translations for brokers. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/pending method=GET tag:service_client
func (s *Service) GetPendingTranslations(ctx context.Context, p translation.GetPendingTranslationsParams) (*translation.GetPendingTranslationsResponse, error) {
	ts := translation.NewTranslationService(s.query)
	return ts.GetPendingTranslations(ctx, p)
}

// TranslateTranslation translates a pending translation. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/translate method=PATCH tag:service_client
func (s *Service) TranslateTranslation(ctx context.Context, p translation.TranslateTranslationParams) error {
	ts := translation.NewTranslationService(s.query)
	return ts.TranslateTranslation(ctx, p)
}

//encore:api auth method=GET path=/broker-translations tag:admin
func (s *Service) ListBrokerTranslations(ctx context.Context, p translation.ListBrokerTranslationsParams) (*translation.ListBrokerTranslationsResponse, error) {
	ts := translation.NewTranslationService(s.query)
	return ts.ListBrokerTranslations(ctx, p)
}

// UpdateBrokerTranslation updates a broker translation target text by ID.
//
//encore:api auth method=PUT path=/broker-translations/:id tag:admin
func (s *Service) UpdateBrokerTranslation(ctx context.Context, id int64, p translation.UpdateBrokerTranslationParams) error {
	ts := translation.NewTranslationService(s.query)
	return ts.UpdateBrokerTranslation(ctx, id, p)
}

// VerifyBrokerTranslation marks a broker translation as verified by ID.
//
//encore:api auth method=PATCH path=/broker-translations/:id/verify tag:admin
func (s *Service) VerifyBrokerTranslation(ctx context.Context, id int64) error {
	ts := translation.NewTranslationService(s.query)
	return ts.VerifyBrokerTranslation(ctx, id)
}

// DeleteBrokerTranslation deletes a broker translation by ID.
//
//encore:api auth method=DELETE path=/broker-translations/:id tag:admin
func (s *Service) DeleteBrokerTranslation(ctx context.Context, id int64) error {
	ts := translation.NewTranslationService(s.query)
	return ts.DeleteBrokerTranslation(ctx, id)
}
