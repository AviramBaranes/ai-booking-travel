package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// KeyFuncHS256 is used to validate the signing method and return the secret key for HS256.
func KeyFuncHS256(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("invalid signing method")
	}

	return []byte(secrets.SecretKey), nil
}

// ValidateAccessToken validates the given JWT access token and returns the claims if valid.
func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, KeyFuncHS256)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid access token")
}

// ValidateRefreshToken validates the given JWT refresh token and returns the claims if valid.
func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, KeyFuncHS256)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}
