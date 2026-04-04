package accounts

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

func TestRegisterAdmin(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation Tests", func(t *testing.T) {
		t.Run("Invalid email", func(t *testing.T) {
			cases := []RegisterAdminParams{
				{Email: "", Password: "StrongPassword123!"},
				{Email: "invalid", Password: "StrongPassword123!"},
				{Email: "invalid@@email.com", Password: "StrongPassword123!"},
			}

			for _, p := range cases {
				err := p.Validate()
				if err == nil {
					t.Errorf("Expected error for email '%s', got nil", p.Email)
				}
			}
		})

		t.Run("Weak password", func(t *testing.T) {
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
					p := RegisterAdminParams{Email: generateTestEmail(), Password: tt.password}
					err := p.Validate()
					api_errors.AssertApiError(t, tt.error, err)
				})
			}
		})
	})

	t.Run("Integration Tests", func(t *testing.T) {
		t.Run("Successful registration", func(t *testing.T) {
			email := generateTestEmail()
			passwordStr := "StrongPassword123!"

			// Use helper to register (which uses real DB)
			admin, cleanup, err := registerAdmin(ctx, email, passwordStr)
			if err != nil {
				t.Fatalf("Failed to register admin: %v", err)
			}
			defer cleanup()

			if admin.ID == 0 {
				t.Error("Expected non-zero ID")
			}

			// Verify in DB
			fetched, err := query.GetUserById(ctx, admin.ID)
			if err != nil {
				t.Fatalf("Failed to get user by ID: %v", err)
			}
			if fetched.Email != email {
				t.Errorf("Expected email %s, got %s", email, fetched.Email)
			}
			if fetched.Role != db.UserRoleAdmin {
				t.Errorf("Expected role admin, got %s", fetched.Role)
			}
			if !password.ComparePassword(fetched.PasswordHash, passwordStr) {
				t.Error("Stored password hash does not match password")
			}
		})

		t.Run("User already exists", func(t *testing.T) {
			email := generateTestEmail()
			passwordStr := "StrongPassword123!"

			_, cleanup, err := registerAdmin(ctx, email, passwordStr)
			if err != nil {
				t.Fatalf("Failed to register first admin: %v", err)
			}
			defer cleanup()

			// Try to register again
			_, err = RegisterAdmin(ctx, RegisterAdminParams{Email: email, Password: passwordStr})
			api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
		})
	})

	t.Run("Unit Tests (Mocks)", func(t *testing.T) {
		t.Run("CheckUserExists fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			email := "test_user@example.com"

			q.EXPECT().
				CheckUserExists(gomock.Any(), email).
				Return(int32(0), errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAdmin(ctx, RegisterAdminParams{Email: email, Password: "StrongPassword123!"})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})

		t.Run("RegisterAdmin DB fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			email := "test_user@example.com"

			q.EXPECT().
				CheckUserExists(gomock.Any(), email).
				Return(int32(0), nil) // User does not exist

			q.EXPECT().
				RegisterAdmin(gomock.Any(), gomock.Any()).
				Return(db.RegisterAdminRow{}, errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAdmin(ctx, RegisterAdminParams{Email: email, Password: "StrongPassword123!"})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})
	})
}
