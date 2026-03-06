package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/auth/db"
	"encore.app/services/auth/jwt"
	"encore.app/services/auth/mocks"

	"go.uber.org/mock/gomock"
)

const (
	testPassword = "ValidPass123!"
	testUsername = "valid_username"
)

type hybridQuerier struct {
	*mocks.MockQuerier
}

func (hq *hybridQuerier) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return query.GetUserByUsername(ctx, username)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid username", func(t *testing.T) {
		cases := []LoginParams{
			{Username: "", Password: testPassword},
			{Username: "ab", Password: testPassword},
			{Username: strings.Repeat("a", 33), Password: testPassword},
			{Username: "לא חוקי", Password: testPassword},
			{Username: "endswith-", Password: testPassword},
		}

		for _, p := range cases {
			err := p.Validate()
			expectedErr := api_errors.WithDetail(err, api_errors.ErrorDetails{
				Field: "username",
				Code:  api_errors.CodeInvalidValue,
			})

			t.Log("Testing with username:", p.Username)
			api_errors.AssertApiError(t, expectedErr, err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		p := LoginParams{
			Username: testUsername,
		}
		err := p.Validate()
		expectedErr := api_errors.WithDetail(err, api_errors.ErrorDetails{
			Field: "password",
			Code:  api_errors.CodeInvalidValue,
		})

		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, err := Login(ctx, LoginParams{Username: testUsername, Password: testPassword})
		api_errors.AssertApiError(t, ErrInvalidCredentials, err)
	})

	t.Run("Search user by username fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			GetUserByUsername(gomock.Any(), testUsername).
			Return(db.User{}, errors.New(("db error")))

		s := &Service{query: q}
		_, err := s.Login(ctx, LoginParams{Username: testUsername, Password: testPassword})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		user, err := RegisterAdmin(ctx, RegisterAdminParams{
			Username: testUsername,
			Password: testPassword,
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		defer query.DeleteUser(ctx, user.ID)

		_, err = Login(ctx, LoginParams{Username: testUsername, Password: "WrongPass123!"})
		api_errors.AssertApiError(t, ErrInvalidCredentials, err)
	})

	t.Run("Store refresh token fails", func(t *testing.T) {
		user, err := RegisterAdmin(ctx, RegisterAdminParams{
			Username: testUsername,
			Password: testPassword,
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		ctrl := gomock.NewController(t)

		// we need to make sure to restore the mock before deleting the user
		defer func() {
			ctrl.Finish()
			query.DeleteUser(ctx, user.ID)
		}()

		// we don't need to mock the login logic
		hq := hybridQuerier{
			MockQuerier: mocks.NewMockQuerier(ctrl),
		}
		hq.EXPECT().
			SaveRefreshToken(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		s := &Service{query: &hq}
		_, err = s.Login(ctx, LoginParams{Username: testUsername, Password: testPassword})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful login", func(t *testing.T) {
		adminUsername := "admin_" + testUsername
		_, delAdmin, err := registerAdmin(ctx, adminUsername, testPassword)

		if err != nil {
			t.Fatalf("Failed to create test admin: %v", err)
		}

		defer delAdmin()

		agentUsername := "agent_" + testUsername
		_, delAgent, err := registerAgent(ctx, RegisterAgentParams{
			Username:   agentUsername,
			Password:   testPassword,
			OfficeCode: "12345",
			AgentCode:  "67890",
		})

		if err != nil {
			t.Fatalf("Failed to create test agent: %v", err)
		}

		defer delAgent()

		cases := []struct {
			name     string
			username string
		}{
			{name: "Admin user", username: adminUsername},
			{name: "Agent user", username: agentUsername},
		}

		for _, c := range cases {

			resp, err := Login(ctx, LoginParams{Username: c.username, Password: testPassword})
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.AccessToken == "" {
				t.Fatal("Expected access token, got empty string")
			}
			if resp.RefreshToken == "" {
				t.Fatal("Expected refresh token, got empty string")
			}

			accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
			if err != nil {
				t.Fatalf("Failed to validate access token: %v", err)
			}

			user, err := query.GetUserByUsername(ctx, c.username)
			if err != nil {
				t.Fatalf("Failed to query user: %v, user: %s", err, c.username)
			}

			assertAccessClaims(t, accessClaims, &user)
			if time.Until(accessClaims.ExpiresAt.Time) <= 0 {
				t.Error("Access token already expired")
			}

			refreshClaims, err := jwt.ValidateRefreshToken(resp.RefreshToken)
			if err != nil {
				t.Fatalf("Failed to validate refresh token: %v", err)
			}
			assertRefreshClaims(t, refreshClaims, &user)
			if time.Until(refreshClaims.ExpiresAt.Time) <= 0 {
				t.Error("Refresh token already expired")
			}

			// Verify stored refresh token in DB
			rt, err := query.GetRefreshToken(ctx, refreshClaims.ID)
			if err != nil {
				t.Fatalf("Failed to retrieve refresh token from DB: %v", err)
			}
			assertTimeAlmostEqual(t, rt.ExpiresAt.Time, refreshClaims.ExpiresAt.Time)
			if rt.UserID != user.ID {
				t.Errorf("Expected token.UserID %d, got %d", user.ID, rt.UserID)
			}
			if rt.Jti != refreshClaims.ID {
				t.Errorf("Expected token.JTI %s, got %s", refreshClaims.ID, rt.Jti)
			}
		}

	})
}
