package reservation

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	authpkg "encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

func authContext(userID int32) context.Context {
	uid := auth.UID(strconv.Itoa(int(userID)))
	return auth.WithContext(context.Background(), uid, &authpkg.AuthData{
		UserID: userID,
		Role:   authpkg.UserRoleAgent,
	})
}

// seedReservation inserts a reservation for the given user and returns its ID.
func seedReservation(t *testing.T, ctx context.Context, s *Service, userID int32, modify func(p *CreateReservationRequest)) int64 {
	t.Helper()
	p := validCreateReservationParams()
	p.UserID = userID
	if modify != nil {
		modify(p)
	}
	resp, err := s.CreateReservation(ctx, *p)
	if err != nil {
		t.Fatalf("failed to seed reservation: %v", err)
	}
	return resp.ID
}

func TestListReservations_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  ListReservationsRequest
		wantErr error
	}{
		{
			name:    "rejects zero page",
			params:  ListReservationsRequest{SortBy: "created_at", Page: 0},
			wantErr: invalidValueErr("page"),
		},
		{
			name:    "rejects negative page",
			params:  ListReservationsRequest{SortBy: "created_at", Page: -1},
			wantErr: invalidValueErr("page"),
		},
		{
			name:    "rejects missing sortBy",
			params:  ListReservationsRequest{Page: 1},
			wantErr: invalidValueErr("sortBy"),
		},
		{
			name:    "rejects invalid sortBy",
			params:  ListReservationsRequest{SortBy: "invalid", Page: 1},
			wantErr: invalidValueErr("sortBy"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api_errors.AssertApiError(t, tc.wantErr, tc.params.Validate())
		})
	}

	t.Run("accepts valid params", func(t *testing.T) {
		if err := (ListReservationsRequest{SortBy: "created_at", Page: 1}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts pickup_date sortBy", func(t *testing.T) {
		if err := (ListReservationsRequest{SortBy: "pickup_date", Page: 1}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestListReservations(t *testing.T) {
	const userID int32 = 42
	ctx := authContext(userID)
	s := &Service{query: testQuerier()}

	// Seed reservations with distinct attributes for filtering tests.
	seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
		p.BrokerReservationID = "LIST-ALICE"
		p.DriverFirstName = "Alice"
		p.DriverLastName = "Smith"
		p.PickupDate = "2026-05-01"
	})
	// Small delay so created_at differs.
	time.Sleep(10 * time.Millisecond)
	seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
		p.BrokerReservationID = "LIST-BOB"
		p.DriverFirstName = "Bob"
		p.DriverLastName = "Jones"
		p.PickupDate = "2026-06-15"
	})

	t.Run("returns all reservations without filters", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) < 2 {
			t.Fatalf("expected at least 2 reservations, got %d", len(resp.Reservations))
		}
	})

	t.Run("total reflects result count", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total < 2 {
			t.Fatalf("expected total >= 2, got %d", resp.Total)
		}
	})

	t.Run("default sort by created_at DESC", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) < 2 {
			t.Fatalf("expected at least 2 reservations, got %d", len(resp.Reservations))
		}
		// Most recent first.
		if resp.Reservations[0].CreatedAt < resp.Reservations[1].CreatedAt {
			t.Fatal("expected reservations ordered by created_at DESC")
		}
	})

	t.Run("sort by pickup_date", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "pickup_date", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) < 2 {
			t.Fatalf("expected at least 2 reservations, got %d", len(resp.Reservations))
		}
		// Verify sorted by pickup_date ASC.
		if resp.Reservations[0].PickupDate > resp.Reservations[1].PickupDate {
			t.Fatalf("expected pickup_date ASC, got %s before %s",
				resp.Reservations[0].PickupDate, resp.Reservations[1].PickupDate)
		}
	})

	t.Run("filter by name matches first name", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Name: "Alice", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-BOB" {
				t.Fatal("Bob should not appear when filtering by Alice")
			}
		}
		found := false
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-ALICE" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected Alice reservation in results")
		}
	})

	t.Run("filter by name matches last name", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Name: "Jones", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-ALICE" {
				t.Fatal("Alice should not appear when filtering by Jones")
			}
		}
		found := false
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-BOB" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected Bob reservation in results")
		}
	})

	t.Run("filter by bookingId", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", BookingID: "LIST-ALICE", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) != 1 {
			t.Fatalf("expected 1 reservation, got %d", len(resp.Reservations))
		}
		if resp.Reservations[0].BrokerReservationID != "LIST-ALICE" {
			t.Fatalf("expected LIST-ALICE, got %s", resp.Reservations[0].BrokerReservationID)
		}
	})

	t.Run("filter by pickupDate", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", PickupDate: "2026-06-15", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-ALICE" {
				t.Fatal("Alice (2026-05-01) should not appear when filtering by 2026-06-15")
			}
		}
		found := false
		for _, r := range resp.Reservations {
			if r.BrokerReservationID == "LIST-BOB" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected Bob reservation with matching pickup date")
		}
	})

	t.Run("filter by status returns only matching", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Status: "booked", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, r := range resp.Reservations {
			if r.Status != "booked" {
				t.Fatalf("expected status 'booked', got %q", r.Status)
			}
		}
	})

	t.Run("non-matching filter returns empty with zero total", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Name: "zzzznonexistent", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) != 0 {
			t.Fatalf("expected 0 reservations, got %d", len(resp.Reservations))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})

	t.Run("total reflects filtered count", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", BookingID: "LIST-BOB", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
		if len(resp.Reservations) != 1 {
			t.Fatalf("expected 1 reservation, got %d", len(resp.Reservations))
		}
	})

	t.Run("returns error when list query fails", func(t *testing.T) {
		q, ms := mockService(t)
		q.EXPECT().ListReservationsByUser(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db error"))

		et.MockService[Interface]("reservation", ms)

		_, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count query fails", func(t *testing.T) {
		q, ms := mockService(t)
		q.EXPECT().ListReservationsByUser(gomock.Any(), gomock.Any()).
			Return([]db.ListReservationsByUserRow{}, nil)
		q.EXPECT().CountReservationsByUser(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("db error"))

		et.MockService[Interface]("reservation", ms)

		_, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestListReservations_EmptyUser(t *testing.T) {
	const emptyUserID int32 = 9999
	ctx := authContext(emptyUserID)

	t.Run("returns empty list when no reservations exist", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{SortBy: "created_at", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) != 0 {
			t.Fatalf("expected 0 reservations, got %d", len(resp.Reservations))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})
}
