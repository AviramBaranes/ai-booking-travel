package jwt

import (
	"fmt"
	"strconv"
	"time"

	"encore.app/services/accounts/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	Issuer          = "ai-booking-api"
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

var secrets struct {
	SecretKey string
}

// AccessTokenData represents the data stored in the JWT access token.
type AccessTokenData struct {
	UserID              int64
	Role                db.UserRole
	AdminRefID          *int64
	OrganizationContext *OrganizationContext
}

// OrganizationContext holds organization-related information to be included in the JWT claims.
type OrganizationContext struct {
	OfficeID       int64 `json:"officeId"`
	OrganizationID int64 `json:"organizationId"`
	IsOrganic      bool  `json:"isOrganic"`
}

// SignAccessToken generates a signed JWT access token for the given user ID and role.
func SignAccessToken(data AccessTokenData) (string, time.Time, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		UserID:              data.UserID,
		Role:                data.Role,
		AdminRefID:          data.AdminRefID,
		OrganizationContext: data.OrganizationContext,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   strconv.FormatInt(int64(data.UserID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	st, err := token.SignedString([]byte(secrets.SecretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token for user_id=%d: %w", data.UserID, err)
	}

	return st, claims.ExpiresAt.Time, nil
}

// SignRefreshToken generates a signed JWT refresh token for the given user ID.
func SignRefreshToken(userID int64) (string, string, time.Time, error) {
	now := time.Now()
	exp := now.Add(refreshTokenTTL)
	jti := uuid.NewString()

	claims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    Issuer,
			Subject:   strconv.FormatInt(int64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secrets.SecretKey))
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("sign refresh token for user_id=%d: %w", userID, err)
	}

	return signedToken, jti, exp, nil
}
