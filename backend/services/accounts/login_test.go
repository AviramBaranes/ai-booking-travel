package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"

	"go.uber.org/mock/gomock"
)

const (
	testPassword = "ValidPass123!"
	testEmail    = "valid_email@example.com"
)

type hybridQuerier struct {
	*mocks.MockQuerier
}

func (hq *hybridQuerier) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return query.GetUserByEmail(ctx, email)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid email", func(t *testing.T) {
		cases := []LoginParams{
			{Email: "", Password: testPassword},
			{Email: "ab", Password: testPassword},
			{Email: "xsxs@@dd.com", Password: testPassword},
		}

		for _, p := range cases {
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "email",
			})

			api_errors.AssertApiError(t, expectedErr, err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		p := LoginParams{
			Email: testEmail,
		}
		err := p.Validate()
		expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code:  api_errors.CodeInvalidValue,
			Field: "password",
		})

		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, err := Login(ctx, LoginParams{Email: testEmail, Password: testPassword})
		api_errors.AssertApiError(t, ErrInvalidCredentials, err)
	})

	t.Run("Search user by email fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			GetUserByEmail(gomock.Any(), testEmail).
			Return(db.User{}, errors.New(("db error")))

		s := &Service{query: q}
		_, err := s.Login(ctx, LoginParams{Email: testEmail, Password: testPassword})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		user, err := CreateAdmin(ctx, CreateAdminRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		defer query.DeleteUser(ctx, user.ID)

		_, err = Login(ctx, LoginParams{Email: testEmail, Password: "WrongPass123!"})
		api_errors.AssertApiError(t, ErrInvalidCredentials, err)
	})

	t.Run("Store refresh token fails", func(t *testing.T) {
		user, err := CreateAdmin(ctx, CreateAdminRequest{
			Email:    testEmail,
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
		_, err = s.Login(ctx, LoginParams{Email: testEmail, Password: testPassword})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful login", func(t *testing.T) {
		adminEmail := "admin_" + testEmail
		_, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)

		if err != nil {
			t.Fatalf("Failed to create test admin: %v", err)
		}

		defer delAdmin()

		agentEmail := "agent_" + testEmail
		_, delAgent, err := createAgent(ctx, CreateAgentRequest{
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: "0505050505",
		})

		if err != nil {
			t.Fatalf("Failed to create test agent: %v", err)
		}

		defer delAgent()

		cases := []struct {
			name  string
			email string
		}{
			{name: "Admin user", email: adminEmail},
			{name: "Agent user", email: agentEmail},
		}

		for _, c := range cases {

			resp, err := Login(ctx, LoginParams{Email: c.email, Password: testPassword})
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

			user, err := query.GetUserByEmail(ctx, c.email)
			if err != nil {
				t.Fatalf("Failed to query user: %v, user: %s", err, c.email)
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
