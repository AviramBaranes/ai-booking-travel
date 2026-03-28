package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"go.uber.org/mock/gomock"
)

// --- Endpoint Tests (real DB) ---

func TestVerifyBrokerTranslation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	t.Run("verifies translation successfully", func(t *testing.T) {
		id := seedTranslation(t, q, "vfy-ok", db.BrokerTranslationStatusPending, 1, strPtr("some text"))

		err := s.VerifyBrokerTranslation(ctx, id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "vfy-ok",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 1 {
			t.Fatalf("expected 1 result, got %d", resp.Total)
		}
		row := resp.Translations[0]
		if row.Status != string(db.BrokerTranslationStatusVerified) {
			t.Fatalf("expected status 'verified', got %q", row.Status)
		}
		if row.ConfidenceScore == nil || *row.ConfidenceScore != 10 {
			t.Fatalf("expected confidence_score 10, got %v", row.ConfidenceScore)
		}
	})

	t.Run("no error for non-existent id", func(t *testing.T) {
		err := s.VerifyBrokerTranslation(ctx, 999999)
		if err != nil {
			t.Fatalf("expected no error for non-existent id, got %v", err)
		}
	})
}

// --- DB Failure Tests (mocked) ---

func TestVerifyBrokerTranslationDBFailures(t *testing.T) {
	ctx := context.Background()

	t.Run("returns internal error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().VerifyBrokerTranslation(gomock.Any(), int32(1)).
			Return(errors.New("db error"))

		err := s.VerifyBrokerTranslation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns not found when db returns ErrNoRows", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().VerifyBrokerTranslation(gomock.Any(), int32(1)).
			Return(db.ErrNoRows)

		err := s.VerifyBrokerTranslation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})
}
