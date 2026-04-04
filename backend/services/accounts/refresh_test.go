package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

type refreshHybridQuerier struct {
	*mocks.MockQuerier
}

func (hq *refreshHybridQuerier) GetRefreshToken(ctx context.Context, id string) (db.RefreshToken, error) {
	return query.GetRefreshToken(ctx, id)
}

func (hq *refreshHybridQuerier) GetUserById(ctx context.Context, id int32) (db.User, error) {
	return query.GetUserById(ctx, id)
}

func TestRefreshTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid refresh token", func(t *testing.T) {
		cases := []string{"", "invalid.token", "invalid"}
		for _, tok := range cases {
			_, err := RefreshTokens(ctx, RefreshTokensParams{RefreshToken: tok})
			api_errors.AssertApiError(t, ErrInvalidRefreshToken, err)
		}
	})

	t.Run("Refresh token not found", func(t *testing.T) {
		// Sign a token but don't save it
		token, _, _, err := jwt.SignRefreshToken(123) // Random user ID
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}
		_, err = RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		api_errors.AssertApiError(t, ErrInvalidRefreshToken, err)
	})

	t.Run("Query refresh token failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		token, jti, _, err := jwt.SignRefreshToken(123)
		if err != nil {
			t.Fatalf("failed to sign refresh token: %v", err)
		}

		q.EXPECT().
			GetRefreshToken(gomock.Any(), jti).
			Return(db.RefreshToken{}, errors.New("db error"))

		s := &Service{query: q}
		_, err = s.RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
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
			ExpiresAt: db.DBTime(exp),
		}
		if err := query.SaveRefreshToken(ctx, p); err != nil {
			t.Fatalf("failed to save expired token: %v", err)
		}

		defer query.DeleteRefreshToken(ctx, jti)

		_, err = RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		api_errors.AssertApiError(t, ErrExpiredRefreshToken, err)
	})

	t.Run("Deleting refresh token failed", func(t *testing.T) {
		_, del, err := registerAdmin(ctx, "del_refresh_fail_user@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer del()

		loginResp, err := Login(ctx, LoginParams{Email: "del_refresh_fail_user@example.com", Password: testPassword})
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}

		claims, err := jwt.ValidateRefreshToken(loginResp.RefreshToken)
		if err != nil {
			t.Fatalf("failed to validate refresh token: %v", err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hq := &refreshHybridQuerier{
			MockQuerier: mocks.NewMockQuerier(ctrl),
		}

		hq.EXPECT().
			DeleteRefreshToken(gomock.Any(), claims.ID).
			Return(errors.New("db error"))

		s := &Service{query: hq}
		_, err = s.RefreshTokens(ctx, RefreshTokensParams{RefreshToken: loginResp.RefreshToken})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, del, err := registerAdmin(ctx, "missing_user_refresh@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		_, err = Login(ctx, LoginParams{Email: "missing_user_refresh@example.com", Password: testPassword})
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
				ExpiresAt: db.DBTime(time.Now().Add(time.Hour)),
			}, nil)

		// DeleteRefreshToken is called before GetUserById
		q.EXPECT().
			DeleteRefreshToken(gomock.Any(), jti).
			Return(nil)

		q.EXPECT().
			GetUserById(gomock.Any(), int32(999999)).
			Return(db.User{}, db.ErrNoRows)

		s := &Service{query: q}
		_, err = s.RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		// The code returns ErrInvalidRefreshToken if user not found (ErrNoRows)
		api_errors.AssertApiError(t, ErrInvalidRefreshToken, err)
	})

	t.Run("Query user failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		token, jti, _, err := jwt.SignRefreshToken(123)
		if err != nil {
			t.Fatal(err)
		}

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			GetRefreshToken(gomock.Any(), jti).
			Return(db.RefreshToken{
				Jti:       jti,
				UserID:    123,
				ExpiresAt: db.DBTime(time.Now().Add(time.Hour)),
			}, nil)

		q.EXPECT().DeleteRefreshToken(gomock.Any(), jti).Return(nil)

		q.EXPECT().
			GetUserById(gomock.Any(), int32(123)).
			Return(db.User{}, errors.New("db error"))

		s := &Service{query: q}
		_, err = s.RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Generating tokens failed", func(t *testing.T) {
		// Mock SaveRefreshToken failure inside generateTokens
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		token, jti, _, err := jwt.SignRefreshToken(123)
		if err != nil {
			t.Fatal(err)
		}

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			GetRefreshToken(gomock.Any(), jti).
			Return(db.RefreshToken{
				Jti:       jti,
				UserID:    123,
				ExpiresAt: db.DBTime(time.Now().Add(time.Hour)),
			}, nil)

		q.EXPECT().DeleteRefreshToken(gomock.Any(), jti).Return(nil)

		q.EXPECT().
			GetUserById(gomock.Any(), int32(123)).
			Return(db.User{ID: 123, Email: "test@example.com"}, nil)

		// generateTokens calls SaveRefreshToken
		q.EXPECT().
			SaveRefreshToken(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		s := &Service{query: q}
		_, err = s.RefreshTokens(ctx, RefreshTokensParams{RefreshToken: token})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful refresh", func(t *testing.T) {
		admin, del, err := registerAdmin(ctx, "refresh_success_user@example.com", testPassword)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer del()

		loginResp, err := Login(ctx, LoginParams{Email: "refresh_success_user@example.com", Password: testPassword})
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}

		origClaims, err := jwt.ValidateRefreshToken(loginResp.RefreshToken)
		if err != nil {
			t.Fatalf("failed to validate login refresh token: %v", err)
		}

		resp, err := RefreshTokens(ctx, RefreshTokensParams{RefreshToken: loginResp.RefreshToken})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.AccessToken == "" {
			t.Fatal("expected access token, got empty string")
		}
		if resp.RefreshToken == "" {
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

		assertAccessClaims(t, accessClaims, &user)
		if time.Until(accessClaims.ExpiresAt.Time) <= 0 {
			t.Error("access token already expired")
		}

		refreshClaims, err := jwt.ValidateRefreshToken(resp.RefreshToken)
		if err != nil {
			t.Fatalf("failed to validate new refresh token: %v", err)
		}
		assertRefreshClaims(t, refreshClaims, &user)
		if time.Until(refreshClaims.ExpiresAt.Time) <= 0 {
			t.Error("refresh token already expired")
		}

		// Verify old refresh token is deleted
		if _, err := query.GetRefreshToken(ctx, origClaims.ID); err == nil {
			t.Error("old refresh token still exists in DB")
		} else if !errors.Is(err, db.ErrNoRows) {
			t.Errorf("expected ErrNoRows for old token, got %v", err)
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
}
