package accounts

import (
	"context"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	ah "encore.app/services/accounts/handlers/auth"
)

// resolveParams builds ResolveSessionParams carrying the given refresh token in
// the Cookie header string, mirroring how a browser sends it.
func resolveParams(token string) ah.ResolveSessionParams {
	return ah.ResolveSessionParams{CookieHeader: "refresh_token=" + token}
}

func TestResolveSession(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid or missing refresh token", func(t *testing.T) {
		cases := []string{"", "invalid.token", "invalid"}
		for _, tok := range cases {
			_, err := ResolveSession(ctx, resolveParams(tok))
			api_errors.AssertApiError(t, ah.ErrInvalidRefreshToken, err)
		}
	})

	t.Run("Refresh token not found in DB", func(t *testing.T) {
		token, _, _, err := jwt.SignRefreshToken(123) // signed but never saved
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}
		_, err = ResolveSession(ctx, resolveParams(token))
		api_errors.AssertApiError(t, ah.ErrInvalidRefreshToken, err)
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		admin, del, err := registerAdmin(ctx, generateTestEmail(), testPassword)
		if err != nil {
			t.Fatalf("failed to register admin: %v", err)
		}
		defer del()

		token, jti, _, err := jwt.SignRefreshToken(admin.ID)
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}

		if err := query.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
			Jti:       jti,
			UserID:    admin.ID,
			ExpiresAt: dbadapters.DBTime(time.Now().Add(-time.Hour)),
		}); err != nil {
			t.Fatalf("failed to save expired token: %v", err)
		}
		defer query.DeleteRefreshToken(ctx, jti)

		_, err = ResolveSession(ctx, resolveParams(token))
		api_errors.AssertApiError(t, ah.ErrExpiredRefreshToken, err)
	})

	t.Run("Resolves user without rotating the token", func(t *testing.T) {
		email := generateTestEmail()
		admin, del, err := registerAdmin(ctx, email, testPassword)
		if err != nil {
			t.Fatalf("failed to register admin: %v", err)
		}
		defer del()

		loginResp, err := Login(ctx, ah.LoginParams{Email: email, Password: testPassword})
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}
		token := refreshCookieToken(t, loginResp)

		claims, err := jwt.ValidateRefreshToken(token)
		if err != nil {
			t.Fatalf("failed to validate refresh token: %v", err)
		}

		resp, err := ResolveSession(ctx, resolveParams(token))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID != admin.ID {
			t.Errorf("expected user ID %d, got %d", admin.ID, resp.ID)
		}
		if resp.Role != db.UserRoleAdmin {
			t.Errorf("expected role %s, got %s", db.UserRoleAdmin, resp.Role)
		}

		// The token must NOT be rotated/consumed: it still exists in the DB and
		// a second call succeeds with the same token.
		if _, err := query.GetRefreshToken(ctx, claims.ID); err != nil {
			t.Fatalf("refresh token was unexpectedly removed: %v", err)
		}
		if _, err := ResolveSession(ctx, resolveParams(token)); err != nil {
			t.Fatalf("expected repeated resolve to succeed, got %v", err)
		}
	})
}
