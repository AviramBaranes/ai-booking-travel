package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/auth/db"
	"encore.app/services/auth/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

func TestRegisterAgent(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation Tests", func(t *testing.T) {
		t.Run("Invalid username", func(t *testing.T) {
			cases := []RegisterAgentParams{
				{Username: "", Password: testPassword, OfficeCode: "12345", AgentCode: "67890"},
				{Username: "ab", Password: testPassword, OfficeCode: "12345", AgentCode: "67890"},
				{Username: strings.Repeat("a", 33), Password: testPassword, OfficeCode: "12345", AgentCode: "67890"},
				{Username: "לא חוקי", Password: testPassword, OfficeCode: "12345", AgentCode: "67890"},
				{Username: "endswith-", Password: testPassword, OfficeCode: "12345", AgentCode: "67890"},
			}

			for _, p := range cases {
				err := p.Validate()
				expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
					Code:  api_errors.CodeInvalidValue,
					Field: "username",
				})

				t.Log("Testing with username:", p.Username)
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
						Username:   testUsername,
						Password:   tt.password,
						OfficeCode: "12345",
						AgentCode:  "67890",
					}
					err := p.Validate()
					api_errors.AssertApiError(t, tt.error, err)
				})
			}
		})

		t.Run("Missing office code", func(t *testing.T) {
			p := RegisterAgentParams{
				Username:   testUsername,
				Password:   testPassword,
				OfficeCode: "",
				AgentCode:  "67890",
			}
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "office_code",
			})
			api_errors.AssertApiError(t, expectedErr, err)
		})

		t.Run("Missing agent code", func(t *testing.T) {
			p := RegisterAgentParams{
				Username:   testUsername,
				Password:   testPassword,
				OfficeCode: "12345",
				AgentCode:  "",
			}
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "agent_code",
			})
			api_errors.AssertApiError(t, expectedErr, err)
		})
	})

	t.Run("Integration Tests", func(t *testing.T) {
		t.Run("Successful registration", func(t *testing.T) {
			username := generateTestUsername()
			passwordStr := testPassword
			officeCode := "OFF123"
			agentCode := "AGT456"

			agent, cleanup, err := registerAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   passwordStr,
				OfficeCode: officeCode,
				AgentCode:  agentCode,
			})
			if err != nil {
				t.Fatalf("Failed to register agent: %v", err)
			}
			defer cleanup()

			if agent.ID == 0 {
				t.Error("Expected non-zero ID")
			}

			// Verify in DB
			fetched, err := query.GetUserByUsername(ctx, username)
			if err != nil {
				t.Fatalf("Failed to get user by username: %v", err)
			}
			if fetched.Username != username {
				t.Errorf("Expected username %s, got %s", username, fetched.Username)
			}
			if fetched.Role != db.UserRoleAgent {
				t.Errorf("Expected role agent, got %s", fetched.Role)
			}
			if !password.ComparePassword(fetched.PasswordHash, passwordStr) {
				t.Error("Stored password hash does not match password")
			}
			if !fetched.OfficeCode.Valid || fetched.OfficeCode.String != officeCode {
				t.Errorf("Expected office code %s, got %v", officeCode, fetched.OfficeCode)
			}
			if !fetched.AgentCode.Valid || fetched.AgentCode.String != agentCode {
				t.Errorf("Expected agent code %s, got %v", agentCode, fetched.AgentCode)
			}
		})

		t.Run("User already exists", func(t *testing.T) {
			username := generateTestUsername()
			passwordStr := testPassword

			_, cleanup, err := registerAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   passwordStr,
				OfficeCode: "12345",
				AgentCode:  "67890",
			})
			if err != nil {
				t.Fatalf("Failed to register first agent: %v", err)
			}
			defer cleanup()

			// Try to register again with same username
			_, err = RegisterAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   passwordStr,
				OfficeCode: "12345",
				AgentCode:  "67890",
			})
			api_errors.AssertApiError(t, ErrUsernameAlreadyExists, err)
		})
	})

	t.Run("Unit Tests (Mocks)", func(t *testing.T) {
		t.Run("CheckUserExists fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			username := "test_agent"

			q.EXPECT().
				CheckUserExists(gomock.Any(), username).
				Return(int32(0), errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   testPassword,
				OfficeCode: "12345",
				AgentCode:  "67890",
			})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})

		t.Run("User already exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			username := "existing_agent"

			q.EXPECT().
				CheckUserExists(gomock.Any(), username).
				Return(int32(123), nil) // User ID exists

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   testPassword,
				OfficeCode: "12345",
				AgentCode:  "67890",
			})
			api_errors.AssertApiError(t, ErrUsernameAlreadyExists, err)
		})

		t.Run("RegisterAgent DB fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mocks.NewMockQuerier(ctrl)
			username := "test_agent"

			q.EXPECT().
				CheckUserExists(gomock.Any(), username).
				Return(int32(0), nil) // User does not exist

			q.EXPECT().
				RegisterAgent(gomock.Any(), gomock.Any()).
				Return(db.RegisterAgentRow{}, errors.New("db error"))

			s := &Service{query: q}
			_, err := s.RegisterAgent(ctx, RegisterAgentParams{
				Username:   username,
				Password:   testPassword,
				OfficeCode: "12345",
				AgentCode:  "67890",
			})
			api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
		})
	})
}
