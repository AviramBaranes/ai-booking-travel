package accounts

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"
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
				expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
					Code:  api_errors.CodeInvalidValue,
					Field: "id",
				})

				api_errors.AssertApiError(t, expectedErr, err)
			}
		})

		t.Run("Invalid password", func(t *testing.T) {
			tests := []struct {
				password string
				error    error
			}{
				{"Short1!", ErrPasswordTooShort},
				{"missing_capital1", ErrPasswordNoUpper},
				{"MISSING_LOWER1", ErrPasswordNoLower},
				{"missingNumber!", ErrPasswordNoNumber},
				{"MissingSymbol1", ErrPasswordNoSymbol},
			}
			for _, tt := range tests {
				t.Run(tt.password, func(t *testing.T) {
					p := RegisterAdminParams{Username: generateTestUsername(), Password: tt.password}
					err := p.Validate()
					api_errors.AssertApiError(t, tt.error, err)
				})
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
