package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

// keyFuncHS256 is used to validate the signing method and return the secret key for HS256.
func keyFuncHS256(token *jwt.Token) (interface{}, error) {
	if token.Method != jwt.SigningMethodHS256 {
		return nil, ErrInvalidSigningMethod
	}

	return []byte(secrets.SecretKey), nil
}

// ValidateAccessToken validates the given JWT access token and returns the claims if valid.
func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, keyFuncHS256)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAccessToken, err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidAccessToken
}

// ValidateRefreshToken validates the given JWT refresh token and returns the claims if valid.
func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, keyFuncHS256)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidRefreshToken
}
