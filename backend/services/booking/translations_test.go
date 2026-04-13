package booking

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	locations_mocks "encore.app/services/booking/mocks"
	"encore.dev/beta/errs"
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

//  UPDATE EP:

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

// LIST EP:

// --- Helpers ---

func translationInvalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

func strPtr(s string) *string { return &s }

// seedTranslation inserts a broker_translation row via the typed query interface.
// It registers a cleanup to delete the row when the test ends.
func seedTranslation(t *testing.T, q *db.Queries, source string, status db.BrokerTranslationStatus, confidence int32, target *string) int32 {
	t.Helper()
	ctx := context.Background()

	id, err := q.InsertBrokerTranslationFull(ctx, db.InsertBrokerTranslationFullParams{
		SourceText:      source,
		TargetText:      target,
		Status:          status,
		ConfidenceScore: &confidence,
	})
	if err != nil {
		t.Fatalf("failed to seed translation: %v", err)
	}

	t.Cleanup(func() {
		_ = q.DeleteBrokerTranslation(ctx, id)
	})

	return id
}

// --- Validation Tests ---

func TestListBrokerTranslationsValidation(t *testing.T) {
	t.Run("rejects page 0", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: 0}
		api_errors.AssertApiError(t, translationInvalidValueErr("page"), p.Validate())
	})

	t.Run("rejects negative page", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: -1}
		api_errors.AssertApiError(t, translationInvalidValueErr("page"), p.Validate())
	})

	t.Run("rejects invalid status", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: 1, Status: "bogus"}
		api_errors.AssertApiError(t, translationInvalidValueErr("status"), p.Validate())
	})

	t.Run("rejects invalid sort direction", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: 1, SortDir: "up"}
		api_errors.AssertApiError(t, translationInvalidValueErr("sortDir"), p.Validate())
	})

	t.Run("accepts valid params", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: 1, Status: "pending", SortDir: "asc", Search: "hello"}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts minimal params", func(t *testing.T) {
		p := ListBrokerTranslationsRequest{Page: 1}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts all valid statuses", func(t *testing.T) {
		for _, status := range []string{"pending", "translated", "verified"} {
			p := ListBrokerTranslationsRequest{Page: 1, Status: status}
			if err := p.Validate(); err != nil {
				t.Fatalf("expected no error for status %q, got %v", status, err)
			}
		}
	})
}

// --- Listing & Filter Tests ---

func TestListBrokerTranslations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Seed translations with distinct attributes for filter/sort tests.
	// Use a unique prefix to avoid collisions with other tests.
	seedTranslation(t, q, "bt-AAA", db.BrokerTranslationStatusPending, 1, nil)
	seedTranslation(t, q, "bt-BBB", db.BrokerTranslationStatusTranslated, 5, strPtr("translated-bbb"))
	seedTranslation(t, q, "bt-CCC", db.BrokerTranslationStatusVerified, 10, strPtr("verified-ccc"))

	t.Run("no filters returns all seeded translations", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		foundA, foundB, foundC := false, false, false
		for _, r := range resp.Translations {
			switch r.SourceText {
			case "bt-AAA":
				foundA = true
			case "bt-BBB":
				foundB = true
			case "bt-CCC":
				foundC = true
			}
		}
		if !foundA || !foundB || !foundC {
			t.Fatal("expected all three seeded translations in results")
		}
	})

	t.Run("filter by status returns only matching translations", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Status: "verified",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Translations {
			if r.Status != "verified" {
				t.Fatalf("expected only verified translations, got status %q", r.Status)
			}
		}

		found := false
		for _, r := range resp.Translations {
			if r.SourceText == "bt-CCC" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected bt-CCC in verified results")
		}
	})

	t.Run("search by source_text returns matching translation", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "bt-AAA",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
		if resp.Translations[0].SourceText != "bt-AAA" {
			t.Fatalf("expected source bt-AAA, got %q", resp.Translations[0].SourceText)
		}
	})

	t.Run("search by target_text returns matching translation", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "translated-bbb",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
		if resp.Translations[0].SourceText != "bt-BBB" {
			t.Fatalf("expected source bt-BBB, got %q", resp.Translations[0].SourceText)
		}
	})

	t.Run("search and status filter combine with AND", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "bt-", Status: "pending",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Translations {
			if r.Status != "pending" {
				t.Fatalf("expected only pending translations, got status %q", r.Status)
			}
		}

		if resp.Total != 1 {
			t.Fatalf("expected total 1 (only bt-AAA is pending), got %d", resp.Total)
		}
	})

	t.Run("non-matching search returns empty", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "zzzznonexistent",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) != 0 {
			t.Fatalf("expected 0 translations, got %d", len(resp.Translations))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})

	t.Run("total reflects filtered count", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Status: "pending",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != int64(len(resp.Translations)) {
			t.Fatalf("expected total %d to match translation count %d", resp.Total, len(resp.Translations))
		}
	})
}

// --- Sorting Tests ---

func TestListBrokerTranslationsSorting(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Seed with distinct confidence scores and a shared prefix for isolation.
	seedTranslation(t, q, "st-LOW", db.BrokerTranslationStatusPending, 1, nil)
	seedTranslation(t, q, "st-MID", db.BrokerTranslationStatusTranslated, 5, strPtr("mid"))
	seedTranslation(t, q, "st-HIGH", db.BrokerTranslationStatusVerified, 10, strPtr("high"))

	t.Run("sort asc returns lowest confidence first", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "st-", SortDir: "asc",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) < 3 {
			t.Fatalf("expected at least 3 results, got %d", len(resp.Translations))
		}

		for i := 1; i < len(resp.Translations); i++ {
			prev := resp.Translations[i-1].ConfidenceScore
			curr := resp.Translations[i].ConfidenceScore
			if prev != nil && curr != nil && *prev > *curr {
				t.Fatalf("expected ascending order, but got %d before %d", *prev, *curr)
			}
		}
	})

	t.Run("sort desc returns highest confidence first", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "st-", SortDir: "desc",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) < 3 {
			t.Fatalf("expected at least 3 results, got %d", len(resp.Translations))
		}

		for i := 1; i < len(resp.Translations); i++ {
			prev := resp.Translations[i-1].ConfidenceScore
			curr := resp.Translations[i].ConfidenceScore
			if prev != nil && curr != nil && *prev < *curr {
				t.Fatalf("expected descending order, but got %d before %d", *prev, *curr)
			}
		}
	})

	t.Run("default sort direction is asc", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{
			Page: 1, Search: "st-",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) < 3 {
			t.Fatalf("expected at least 3 results, got %d", len(resp.Translations))
		}

		for i := 1; i < len(resp.Translations); i++ {
			prev := resp.Translations[i-1].ConfidenceScore
			curr := resp.Translations[i].ConfidenceScore
			if prev != nil && curr != nil && *prev > *curr {
				t.Fatalf("expected ascending (default) order, but got %d before %d", *prev, *curr)
			}
		}
	})
}

// --- Pagination Tests ---

func TestListBrokerTranslationsPagination(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Seed 16 translations (TranslationsLimit=15), so page 1 has 15, page 2 has 1.
	prefix := "pg-"
	for i := 1; i <= 16; i++ {
		source := fmt.Sprintf("%s%02d", prefix, i)
		seedTranslation(t, q, source, db.BrokerTranslationStatusPending, int32(i%10), nil)
	}

	t.Run("page 1 returns exactly 15 results", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) != TranslationsLimit {
			t.Fatalf("expected %d translations on page 1, got %d", TranslationsLimit, len(resp.Translations))
		}
	})

	t.Run("total reflects all matching rows across pages", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 16 {
			t.Fatalf("expected total 16, got %d", resp.Total)
		}
	})

	t.Run("page 2 returns the remaining result", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 2, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) != 1 {
			t.Fatalf("expected 1 translation on page 2, got %d", len(resp.Translations))
		}
	})

	t.Run("page 2 does not repeat page 1 results", func(t *testing.T) {
		page1, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		page2, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 2, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		page1IDs := map[int32]bool{}
		for _, r := range page1.Translations {
			page1IDs[r.ID] = true
		}
		for _, r := range page2.Translations {
			if page1IDs[r.ID] {
				t.Fatalf("translation ID %d appeared on both page 1 and page 2", r.ID)
			}
		}
	})

	t.Run("page 3 returns empty", func(t *testing.T) {
		resp, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 3, Search: prefix})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Translations) != 0 {
			t.Fatalf("expected 0 translations on page 3, got %d", len(resp.Translations))
		}
	})
}

// --- DB Failure Tests (mocked) ---

func TestListBrokerTranslationsDBFailures(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error when count query fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		q := locations_mocks.NewMockQuerier(ctrl)
		s := &Service{query: q}

		q.EXPECT().CountAllTranslations(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("db error"))

		_, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when list query fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		q := locations_mocks.NewMockQuerier(ctrl)
		s := &Service{query: q}

		q.EXPECT().CountAllTranslations(gomock.Any(), gomock.Any()).
			Return(int64(5), nil)
		q.EXPECT().ListAllTranslations(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db error"))

		_, err := s.ListBrokerTranslations(ctx, ListBrokerTranslationsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

// DELETE EP:

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
