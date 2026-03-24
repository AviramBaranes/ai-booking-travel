package reservation

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"encore.app/internal/api_errors"
	authpkg "encore.app/services/auth"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/mocks"
	"encore.dev/beta/auth"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

func authContext(userID int32) context.Context {
	uid := auth.UID(strconv.Itoa(int(userID)))
	return auth.WithContext(context.Background(), uid, &authpkg.AuthData{
		UserID:   userID,
		Role:     authpkg.UserRoleAgent,
		Username: "testuser",
	})
}

func TestListReservations_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  ListReservationsRequest
		wantErr error
	}{
		{
			name:    "rejects zero page",
			params:  ListReservationsRequest{Page: 0},
			wantErr: invalidValueErr("page"),
		},
		{
			name:    "rejects negative page",
			params:  ListReservationsRequest{Page: -1},
			wantErr: invalidValueErr("page"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api_errors.AssertApiError(t, tc.wantErr, tc.params.Validate())
		})
	}

	t.Run("accepts valid page", func(t *testing.T) {
		if err := (ListReservationsRequest{Page: 1}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestListReservations(t *testing.T) {
	const userID int32 = 42
	ctx := authContext(userID)

	t.Run("returns empty list when no reservations", func(t *testing.T) {
		resp, err := ListReservations(ctx, ListReservationsRequest{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) != 0 {
			t.Fatalf("expected 0 reservations, got %d", len(resp.Reservations))
		}
	})

	t.Run("returns reservations for user", func(t *testing.T) {
		s := &Service{query: testQuerier()}

		// Insert two reservations for our user.
		for i := 0; i < 2; i++ {
			p := validCreateReservationParams()
			p.UserID = userID
			p.BrokerReservationID = "LIST-" + strconv.Itoa(i)
			_, err := s.CreateReservation(ctx, *p)
			if err != nil {
				t.Fatalf("failed to insert reservation %d: %v", i, err)
			}
		}

		resp, err := ListReservations(ctx, ListReservationsRequest{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Reservations) != 2 {
			t.Fatalf("expected 2 reservations, got %d", len(resp.Reservations))
		}

		// Verify ordered by created_at DESC (most recent first).
		if resp.Reservations[0].CreatedAt < resp.Reservations[1].CreatedAt {
			t.Fatal("expected reservations ordered by created_at DESC")
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().ListReservationsByUser(gomock.Any(), db.ListReservationsByUserParams{
			UserID:      userID,
			QueryLimit:  10,
			QueryOffset: 0,
		}).Return(nil, errors.New("db error"))

		et.MockService[Interface]("reservation", &Service{query: q})

		_, err := ListReservations(ctx, ListReservationsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
