package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

func generateTestUsername() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func TestRegisterAdmin(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation Tests", func(t *testing.T) {
		t.Run("Invalid username", func(t *testing.T) {
			cases := []RegisterAdminParams{
				{Username: "", Password: "StrongPassword123!"},
				{Username: "ab", Password: "StrongPassword123!"}, // Too short?
				{Username: "endswith-", Password: "StrongPassword123!"},
				{Username: "invalid@email.com", Password: "StrongPassword123!"}, // Assuming @ is not allowed in username
			}

			for _, p := range cases {
				err := p.Validate()
				if err == nil {
					t.Errorf("Expected error for username '%s', got nil", p.Username)
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
					p := RegisterAdminParams{Username: generateTestUsername(), Password: tt.password}
					err := p.Validate()
					api_errors.AssertApiError(t, tt.error, err)
				})
			}
		})

		t.Run("Missing office/agent code", func(t *testing.T) {
			code := "123"
			cases := []RegisterAdminParams{
				{Username: generateTestUsername(), Password: "StrongPassword123!", OfficeCode: &code, AgentCode: nil},
				{Username: generateTestUsername(), Password: "StrongPassword123!", OfficeCode: nil, AgentCode: &code},
			}

			for _, p := range cases {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				q := mocks.NewMockQuerier(ctrl)
				s := &Service{query: q}

				// CheckUserExists happens before the office/agent code check
				q.EXPECT().CheckUserExists(gomock.Any(), p.Username).Return(int32(0), nil)

				_, err := s.RegisterAdmin(ctx, p)
				if err == nil {
					t.Error("Expected error, got nil")
				}
			}
		})
	})

	t.Run("Integration Tests", func(t *testing.T) {
		t.Run("Successful registration", func(t *testing.T) {
			username := generateTestUsername()
			passwordStr := "StrongPassword123!"

			// Use helper to register (which uses real DB)
			admin, cleanup, err := registerAdmin(ctx, username, passwordStr)
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
			if fetched.Username != username {
				t.Errorf("Expected username %s, got %s", username, fetched.Username)
			}
			if fetched.Role != db.UserRoleAdmin {
				t.Errorf("Expected role admin, got %s", fetched.Role)
			}
			if !password.ComparePassword(fetched.PasswordHash, passwordStr) {
				t.Error("Stored password hash does not match password")
			}
		})

		t.Run("User already exists", func(t *testing.T) {
			username := generateTestUsername()
			passwordStr := "StrongPassword123!"

			_, cleanup, err := registerAdmin(ctx, username, passwordStr)
			if err != nil {
				t.Fatalf("Failed to register first admin: %v", err)
			}
			defer cleanup()

			// Try to register again
			_, err = RegisterAdmin(ctx, RegisterAdminParams{Username: username, Password: passwordStr})
			api_errors.AssertApiError(t, ErrUsernameAlreadyExists, err)
		})
	})

	t.Run("Unit Tests (Mocks)", func(t *testing.T) {
		t.Run("CheckUserExists fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			username := "test_user"

			q.EXPECT().
				CheckUserExists(gomock.Any(), username).
				Return(int32(0), errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAdmin(ctx, RegisterAdminParams{Username: username, Password: "StrongPassword123!"})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})

		t.Run("RegisterAdmin DB fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			username := "test_user"

			q.EXPECT().
				CheckUserExists(gomock.Any(), username).
				Return(int32(0), nil) // User does not exist

			q.EXPECT().
				RegisterAdmin(gomock.Any(), gomock.Any()).
				Return(db.RegisterAdminRow{}, errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAdmin(ctx, RegisterAdminParams{Username: username, Password: "StrongPassword123!"})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})
	})
}
