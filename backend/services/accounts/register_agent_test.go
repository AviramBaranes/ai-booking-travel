package accounts

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

func TestRegisterAgent(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation Tests", func(t *testing.T) {
		t.Run("Invalid email", func(t *testing.T) {
			officeID := int32(12345)
			phoneNumber := "0505050505"
			password := "StrongPassword123!"
			cases := []RegisterAgentParams{
				{Email: "", Password: password, OfficeID: officeID, PhoneNumber: phoneNumber},
				{Email: "invalid", Password: password, OfficeID: officeID, PhoneNumber: phoneNumber},
				{Email: "invalid@@email.com", Password: password, OfficeID: officeID, PhoneNumber: phoneNumber},
			}
			for _, p := range cases {
				err := p.Validate()
				expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
					Code:  api_errors.CodeInvalidValue,
					Field: "email",
				})

				t.Log("Testing with email:", p.Email)
				api_errors.AssertApiError(t, expectedErr, err)
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
					p := RegisterAgentParams{
						Email:    testEmail,
						Password: tt.password,
					}
					err := p.Validate()
					api_errors.AssertApiError(t, tt.error, err)
				})
			}
		})

		t.Run("Missing office id", func(t *testing.T) {
			p := RegisterAgentParams{
				Email:       testEmail,
				Password:    testPassword,
				PhoneNumber: "0505050505",
			}
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "officeId",
			})
			api_errors.AssertApiError(t, expectedErr, err)
		})

		t.Run("Missing office phone number", func(t *testing.T) {
			p := RegisterAgentParams{
				Email:    testEmail,
				Password: testPassword,
				OfficeID: int32(12345),
			}
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "phoneNumber",
			})
			api_errors.AssertApiError(t, expectedErr, err)
		})
	})

	t.Run("Integration Tests", func(t *testing.T) {
		t.Run("Successful registration", func(t *testing.T) {
			email := generateTestEmail()
			passwordStr := testPassword
			phoneNumber := randomIsraeliPhoneNumber()

			agent, cleanup, err := registerAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    passwordStr,
				PhoneNumber: phoneNumber,
			})
			if err != nil {
				t.Fatalf("Failed to register agent: %v", err)
			}
			defer cleanup()

			if agent.ID == 0 {
				t.Error("Expected non-zero ID")
			}

			// Verify in DB
			fetched, err := query.GetUserByEmail(ctx, email)
			if err != nil {
				t.Fatalf("Failed to get user by email: %v", err)
			}
			if fetched.Email != email {
				t.Errorf("Expected email %s, got %s", email, fetched.Email)
			}
			if fetched.Role != db.UserRoleAgent {
				t.Errorf("Expected role agent, got %s", fetched.Role)
			}
			if !password.ComparePassword(fetched.PasswordHash, passwordStr) {
				t.Error("Stored password hash does not match password")
			}
		})

		t.Run("User already exists", func(t *testing.T) {
			email := generateTestEmail()
			passwordStr := testPassword

			_, cleanup, err := registerAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    passwordStr,
				PhoneNumber: randomIsraeliPhoneNumber(),
			})
			if err != nil {
				t.Fatalf("Failed to register first agent: %v", err)
			}
			defer cleanup()

			// Try to register again with same email
			_, err = RegisterAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    passwordStr,
				OfficeID:    int32(12345),
				PhoneNumber: randomIsraeliPhoneNumber(),
			})
			api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
		})
	})

	t.Run("Unit Tests (Mocks)", func(t *testing.T) {
		t.Run("CheckUserExists fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			email := "test_agent@example.com"

			q.EXPECT().
				CheckUserExists(gomock.Any(), email).
				Return(int32(0), errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    testPassword,
				OfficeID:    int32(12345),
				PhoneNumber: randomIsraeliPhoneNumber(),
			})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})

		t.Run("User already exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			email := "existing_agent@example.com"

			q.EXPECT().
				CheckUserExists(gomock.Any(), email).
				Return(int32(123), nil) // User ID exists

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    testPassword,
				OfficeID:    int32(12345),
				PhoneNumber: randomIsraeliPhoneNumber(),
			})
			api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
		})

		t.Run("RegisterAgent DB fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			email := "test_agent@example.com"

			q.EXPECT().
				CheckUserExists(gomock.Any(), email).
				Return(int32(0), nil) // User does not exist

			q.EXPECT().
				RegisterAgent(gomock.Any(), gomock.Any()).
				Return(db.RegisterAgentRow{}, errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Email:       email,
				Password:    testPassword,
				OfficeID:    int32(12345),
				PhoneNumber: randomIsraeliPhoneNumber(),
			})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})
	})
}
