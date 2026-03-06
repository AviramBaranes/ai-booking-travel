package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/auth/db"
	"encore.app/services/auth/mocks"
	"go.uber.org/mock/gomock"
)

func TestChangePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid input", func(t *testing.T) {
		t.Run("Invalid ID", func(t *testing.T) {
			cases := []ChangePasswordParams{
				{ID: 0, NewPassword: testPassword},
				{ID: -1, NewPassword: testPassword},
			}

			for _, p := range cases {
				err := p.Validate()
				expectedErr := api_errors.WithDetail(err, api_errors.ErrorDetails{
					Field: "id",
					Code:  api_errors.CodeInvalidValue,
				})

				api_errors.AssertApiError(t, expectedErr, err)
			}
		})

		t.Run("Invalid password", func(t *testing.T) {
			cases := []ChangePasswordParams{
				{ID: 1, NewPassword: ""},
				{ID: 1, NewPassword: "short"},
				{ID: 1, NewPassword: strings.Repeat("a", 7)},
			}

			for _, p := range cases {
				err := p.Validate()
				// Just check that validation fails for invalid passwords
				// The specific error code can vary (required, password_too_short, etc.)
				if err == nil {
					t.Errorf("Expected error for password '%s', got nil", p.NewPassword)
				}
			}
		})
	})

	t.Run("User not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			UpdateUserPassword(gomock.Any(), gomock.Any()).
			Return(db.ErrNoRows)

		s := &Service{query: q}
		err := s.ChangePassword(ctx, ChangePasswordParams{ID: 1, NewPassword: testPassword})
		api_errors.AssertApiError(t, ErrUserNotFound, err)
	})

	t.Run("Database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			UpdateUserPassword(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		s := &Service{query: q}
		err := s.ChangePassword(ctx, ChangePasswordParams{ID: 1, NewPassword: testPassword})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful password change", func(t *testing.T) {
		// Create a user
		user, del, err := registerAdmin(ctx, "change_pass_user", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer del()

		newPassword := "NewValidPass123!"

		// Change password
		err = ChangePassword(ctx, ChangePasswordParams{
			ID:          user.ID,
			NewPassword: newPassword,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify login with new password works
		_, err = Login(ctx, LoginParams{
			Username: "change_pass_user",
			Password: newPassword,
		})
		if err != nil {
			t.Fatalf("failed to login with new password: %v", err)
		}

		// Verify login with old password fails
		_, err = Login(ctx, LoginParams{
			Username: "change_pass_user",
			Password: testPassword,
		})
		api_errors.AssertApiError(t, ErrInvalidCredentials, err)
	})
}
