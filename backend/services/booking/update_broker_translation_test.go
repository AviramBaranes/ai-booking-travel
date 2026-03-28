package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"go.uber.org/mock/gomock"
)

// --- Validation Tests ---

func TestUpdateBrokerTranslationValidation(t *testing.T) {
	t.Run("rejects empty target_text", func(t *testing.T) {
		p := UpdateBrokerTranslationRequest{TargetText: ""}
		api_errors.AssertApiError(t, translationInvalidValueErr("target_text"), p.Validate())
	})

	t.Run("rejects blank target_text", func(t *testing.T) {
		p := UpdateBrokerTranslationRequest{TargetText: "   "}
		api_errors.AssertApiError(t, translationInvalidValueErr("target_text"), p.Validate())
	})

	t.Run("accepts valid target_text", func(t *testing.T) {
		p := UpdateBrokerTranslationRequest{TargetText: "hello world"}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

// --- Endpoint Tests (real DB) ---

func TestUpdateBrokerTranslation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	t.Run("updates target text successfully", func(t *testing.T) {
		id := seedTranslation(t, q, "upd-ok", db.BrokerTranslationStatusPending, 1, nil)

		err := s.UpdateBrokerTranslation(ctx, id, UpdateBrokerTranslationRequest{
			TargetText: "new translation",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify the row was actually updated by listing it back.
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "upd-ok",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 1 {
			t.Fatalf("expected 1 result, got %d", resp.Total)
		}
		row := resp.Translations[0]
		if row.TargetText == nil || *row.TargetText != "new translation" {
			t.Fatalf("expected target_text 'new translation', got %v", row.TargetText)
		}
		if row.Status != string(db.BrokerTranslationStatusVerified) {
			t.Fatalf("expected status 'verified', got %q", row.Status)
		}
		if row.ConfidenceScore == nil || *row.ConfidenceScore != 10 {
			t.Fatalf("expected confidence_score 10, got %v", row.ConfidenceScore)
		}
	})

	t.Run("no error for non-existent id", func(t *testing.T) {
		err := s.UpdateBrokerTranslation(ctx, 999999, UpdateBrokerTranslationRequest{
			TargetText: "does not matter",
		})
		// The :exec query won't trigger ErrNoRows (UPDATE affects 0 rows silently),
		// so this should succeed without error.
		if err != nil {
			t.Fatalf("expected no error for non-existent id, got %v", err)
		}
	})
}

// --- DB Failure Tests (mocked) ---

func TestUpdateBrokerTranslationDBFailures(t *testing.T) {
	ctx := context.Background()

	t.Run("returns internal error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateBrokerTranslation(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		err := s.UpdateBrokerTranslation(ctx, 1, UpdateBrokerTranslationRequest{
			TargetText: "fail",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns not found when db returns ErrNoRows", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateBrokerTranslation(gomock.Any(), gomock.Any()).
			Return(db.ErrNoRows)

		err := s.UpdateBrokerTranslation(ctx, 1, UpdateBrokerTranslationRequest{
			TargetText: "missing",
		})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})
}
