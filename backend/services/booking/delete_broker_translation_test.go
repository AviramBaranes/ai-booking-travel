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

func TestDeleteBrokerTranslation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	t.Run("deletes translation successfully", func(t *testing.T) {
		id := seedTranslation(t, q, "del-ok", db.BrokerTranslationStatusPending, 1, nil)

		err := s.DeleteBrokerTranslation(ctx, id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's gone.
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "del-ok",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 0 {
			t.Fatalf("expected 0 results after delete, got %d", resp.Total)
		}
	})

	t.Run("no error for non-existent id", func(t *testing.T) {
		err := s.DeleteBrokerTranslation(ctx, 999999)
		if err != nil {
			t.Fatalf("expected no error for non-existent id, got %v", err)
		}
	})
}

// --- DB Failure Tests (mocked) ---

func TestDeleteBrokerTranslationDBFailures(t *testing.T) {
	ctx := context.Background()

	t.Run("returns internal error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteBrokerTranslation(gomock.Any(), int32(1)).
			Return(errors.New("db error"))

		err := s.DeleteBrokerTranslation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns not found when db returns ErrNoRows", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteBrokerTranslation(gomock.Any(), int32(1)).
			Return(db.ErrNoRows)

		err := s.DeleteBrokerTranslation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})
}
