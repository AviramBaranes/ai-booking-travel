package accounts

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	ah "encore.app/services/accounts/handlers/auth"
	user "encore.app/services/accounts/handlers/user"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/auth"
	"go.uber.org/mock/gomock"
)

func TestRefreshTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid refresh token", func(t *testing.T) {
		cases := []string{"", "invalid.token", "invalid"}
		for _, tok := range cases {
			_, err := RefreshTokens(ctx, refreshParams(tok))
			api_errors.AssertApiError(t, ah.ErrInvalidRefreshToken, err)
		}
	})

	t.Run("Refresh token not found", func(t *testing.T) {
		// Sign a token but don't save it
		token, _, _, err := jwt.SignRefreshToken(123) // Random user ID
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}
		_, err = RefreshTokens(ctx, refreshParams(token))
		api_errors.AssertApiError(t, ah.ErrInvalidRefreshToken, err)
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		// Create a user first
		admin, del, err := registerAdmin(ctx, "expired_refresh_user@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer del()

		token, jti, _, err := jwt.SignRefreshToken(admin.ID)
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}

		// Save it as expired
		exp := time.Now().Add(-time.Hour)
		p := db.SaveRefreshTokenParams{
			Jti:       jti,
			UserID:    admin.ID,
			ExpiresAt: dbadapters.DBTime(exp),
		}
		if err := query.SaveRefreshToken(ctx, p); err != nil {
			t.Fatalf("failed to save expired token: %v", err)
		}

		defer query.DeleteRefreshToken(ctx, jti)

		_, err = RefreshTokens(ctx, refreshParams(token))
		api_errors.AssertApiError(t, ah.ErrExpiredRefreshToken, err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, del, err := registerAdmin(ctx, "missing_user_refresh@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		_, err = Login(ctx, ah.LoginParams{Email: "missing_user_refresh@example.com", Password: testPassword})
		if err != nil {
			del()
			t.Fatalf("failed to login: %v", err)
		}

		del()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		token, jti, _, err := jwt.SignRefreshToken(999999) // Non-existent user
		if err != nil {
			t.Fatalf("failed to sign: %v", err)
		}

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			GetRefreshToken(gomock.Any(), jti).
			Return(db.RefreshToken{
				Jti:       jti,
				UserID:    999999,
				ExpiresAt: dbadapters.DBTime(time.Now().Add(time.Hour)),
			}, nil)

		q.EXPECT().
			GetUserById(gomock.Any(), int64(999999)).
			Return(db.GetUserByIdRow{}, db.ErrNoRows)

		s := &Service{query: q}
		_, err = s.RefreshTokens(ctx, refreshParams(token))
		// The code returns ah.ErrInvalidRefreshToken if user not found (ErrNoRows)
		api_errors.AssertApiError(t, ah.ErrInvalidRefreshToken, err)
	})

	t.Run("Successful refresh", func(t *testing.T) {
		admin, del, err := registerAdmin(ctx, "refresh_success_user@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer del()

		loginResp, err := Login(ctx, ah.LoginParams{Email: "refresh_success_user@example.com", Password: testPassword})
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}

		origClaims, err := jwt.ValidateRefreshToken(refreshCookieToken(t, loginResp))
		if err != nil {
			t.Fatalf("failed to validate login refresh token: %v", err)
		}

		resp, err := RefreshTokens(ctx, refreshParams(refreshCookieToken(t, loginResp)))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.AccessToken == "" {
			t.Fatal("expected access token, got empty string")
		}
		newRefreshToken := refreshCookieToken(t, resp)
		if newRefreshToken == "" {
			t.Fatal("expected refresh token, got empty string")
		}

		accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
		if err != nil {
			t.Fatalf("failed to validate new access token: %v", err)
		}

		// Get user to compare
		user, err := query.GetUserById(ctx, admin.ID)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		assertAccessClaims(t, accessClaims, jwt.AccessTokenData{
			UserID: admin.ID,
			Role:   db.UserRoleAdmin,
		})
		if time.Until(accessClaims.ExpiresAt.Time) <= 0 {
			t.Error("access token already expired")
		}

		refreshClaims, err := jwt.ValidateRefreshToken(newRefreshToken)
		if err != nil {
			t.Fatalf("failed to validate new refresh token: %v", err)
		}
		assertRefreshClaims(t, refreshClaims, admin.ID)
		if time.Until(refreshClaims.ExpiresAt.Time) <= 0 {
			t.Error("refresh token already expired")
		}

		// Refresh tokens are single-use: the old jti must be consumed (rotated out).
		if _, err := query.GetRefreshToken(ctx, origClaims.ID); !errors.Is(err, db.ErrNoRows) {
			t.Errorf("expected old refresh token to be rotated out (deleted), got err=%v", err)
		}

		// Verify new refresh token is in DB
		rt, err := query.GetRefreshToken(ctx, refreshClaims.ID)
		if err != nil {
			t.Fatalf("failed to retrieve new refresh token: %v", err)
		}
		assertTimeAlmostEqual(t, rt.ExpiresAt.Time, refreshClaims.ExpiresAt.Time)
		if rt.UserID != user.ID {
			t.Errorf("expected token.UserID %d, got %d", user.ID, rt.UserID)
		}
		if rt.Jti != refreshClaims.ID {
			t.Errorf("expected token.JTI %s, got %s", refreshClaims.ID, rt.Jti)
		}
	})

	t.Run("Refresh preserves admin ref ID", func(t *testing.T) {
		admin, delAdmin, err := registerAdmin(ctx, generateTestEmail(), testPassword)
		if err != nil {
			t.Fatalf("failed to create admin: %v", err)
		}
		defer delAdmin()

		agentEmail := generateTestEmail()
		agent, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("failed to create agent: %v", err)
		}
		defer delAgent()

		// Login as agent (creates tokens with admin ref)
		adminCtx := auth.WithContext(ctx, auth.UID(strconv.Itoa(int(admin.ID))), &AuthData{
			UserID: admin.ID,
			Role:   UserRoleAdmin,
		})
		loginResp, err := LoginAsAgent(adminCtx, ah.LoginAsAgentParams{AgentID: agent.ID})
		if err != nil {
			t.Fatalf("failed to login as agent: %v", err)
		}

		// Refresh the tokens
		refreshResp, err := RefreshTokens(ctx, refreshParams(refreshCookieToken(t, loginResp)))
		if err != nil {
			t.Fatalf("expected no error on refresh, got %v", err)
		}

		// Verify new access token still has admin ref ID
		accessClaims, err := jwt.ValidateAccessToken(refreshResp.AccessToken)
		if err != nil {
			t.Fatalf("failed to validate refreshed access token: %v", err)
		}
		if accessClaims.AdminRefID == nil {
			t.Fatal("expected AdminRefID in refreshed access token")
		}
		if *accessClaims.AdminRefID != admin.ID {
			t.Errorf("expected AdminRefID %d, got %d", admin.ID, *accessClaims.AdminRefID)
		}

		// Verify new refresh token row also has admin ref
		refreshClaims, err := jwt.ValidateRefreshToken(refreshCookieToken(t, refreshResp))
		if err != nil {
			t.Fatalf("failed to validate refreshed refresh token: %v", err)
		}
		rt, err := query.GetRefreshToken(ctx, refreshClaims.ID)
		if err != nil {
			t.Fatalf("failed to get stored refresh token: %v", err)
		}
		if rt.AdminRefID == nil {
			t.Fatal("expected AdminRefID in stored refresh token after refresh")
		}
		if *rt.AdminRefID != admin.ID {
			t.Errorf("expected stored AdminRefID %d, got %d", admin.ID, *rt.AdminRefID)
		}
	})
}
