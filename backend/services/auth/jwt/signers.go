package jwt

import (
	"strconv"
	"time"

	"encore.app/services/auth/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	Issuer          = "global-rental-api"
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

var secrets struct {
	SecretKey string
}

// SignAccessToken generates a signed JWT access token for the given user ID and role.
func SignAccessToken(user db.User) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		UserID:     user.ID,
		Role:       user.Role,
		Username:   user.Username,
		AgentCode:  user.AgentCode.String,
		OfficeCode: user.OfficeCode.String,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   strconv.FormatInt(int64(user.ID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secrets.SecretKey))
}

// SignRefreshToken generates a signed JWT refresh token for the given user ID.
func SignRefreshToken(userID int32) (string, string, time.Time, error) {
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
	return signedToken, jti, exp, err
}
