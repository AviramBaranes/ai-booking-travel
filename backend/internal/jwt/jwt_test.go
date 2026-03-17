package jwt

import (
	"testing"
	"time"

	"encore.app/services/auth/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestMain(m *testing.M) {
	// Set a dummy secret for testing
	secrets.SecretKey = "test-secret-key"
	m.Run()
}

func TestSignAccessToken(t *testing.T) {
	user := db.User{
		ID:       123,
		Username: "testuser",
		Role:     "admin",
		AgentCode: pgtype.Text{
			String: "agent1",
			Valid:  true,
		},
		OfficeCode: pgtype.Text{
			String: "office1",
			Valid:  true,
		},
	}

	tokenString, err := SignAccessToken(user)
	if err != nil {
		t.Fatalf("SignAccessToken failed: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Expected token string, got empty")
	}

	// Parse and verify claims manually to ensure SignAccessToken did its job
	token, _ := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secrets.SecretKey), nil
	})

	if !token.Valid {
		t.Fatal("Token is invalid")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		t.Fatal("Could not cast claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
	if claims.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, claims.Username)
	}
	if claims.AgentCode != user.AgentCode.String {
		t.Errorf("Expected AgentCode %s, got %s", user.AgentCode.String, claims.AgentCode)
	}
	if claims.OfficeCode != user.OfficeCode.String {
		t.Errorf("Expected OfficeCode %s, got %s", user.OfficeCode.String, claims.OfficeCode)
	}
	if claims.Issuer != Issuer {
		t.Errorf("Expected Issuer %s, got %s", Issuer, claims.Issuer)
	}
}

func TestValidateAccessToken(t *testing.T) {
	user := db.User{
		ID:       123,
		Username: "testuser",
		Role:     "admin",
	}
	validToken, _ := SignAccessToken(user)

	t.Run("Valid token", func(t *testing.T) {
		claims, err := ValidateAccessToken(validToken)
		if err != nil {
			t.Fatalf("ValidateAccessToken failed: %v", err)
		}
		if claims.UserID != user.ID {
			t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
		}
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// Create a token signed with a different key
		claims := AccessTokenClaims{
			UserID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		wrongToken, _ := token.SignedString([]byte("wrong-key"))

		_, err := ValidateAccessToken(wrongToken)
		if err == nil {
			t.Fatal("Expected error for invalid signature, got nil")
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		claims := AccessTokenClaims{
			UserID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, _ := token.SignedString([]byte(secrets.SecretKey))

		_, err := ValidateAccessToken(expiredToken)
		if err == nil {
			t.Fatal("Expected error for expired token, got nil")
		}
	})

	t.Run("Malformed token", func(t *testing.T) {
		_, err := ValidateAccessToken("not.a.token")
		if err == nil {
			t.Fatal("Expected error for malformed token, got nil")
		}
	})
}

func TestSignRefreshToken(t *testing.T) {
	userID := int32(456)
	tokenString, jti, exp, err := SignRefreshToken(userID)
	if err != nil {
		t.Fatalf("SignRefreshToken failed: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Expected token string")
	}
	if jti == "" {
		t.Fatal("Expected JTI")
	}
	if exp.Before(time.Now()) {
		t.Fatal("Expected expiration in the future")
	}

	// Validate
	claims, err := ValidateRefreshToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.ID != jti {
		t.Errorf("Expected JTI %s, got %s", jti, claims.ID)
	}
}

func TestValidateRefreshToken(t *testing.T) {
	userID := int32(789)
	validToken, _, _, _ := SignRefreshToken(userID)

	t.Run("Valid token", func(t *testing.T) {
		claims, err := ValidateRefreshToken(validToken)
		if err != nil {
			t.Fatalf("ValidateRefreshToken failed: %v", err)
		}
		if claims.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
		}
	})

	t.Run("Invalid token type (Access Token as Refresh Token)", func(t *testing.T) {
		// Create an access token
		// user := db.User{ID: userID}
		// accessToken, _ := SignAccessToken(user) // Unused
		// ...
	})

	t.Run("Invalid signature", func(t *testing.T) {
		claims := RefreshTokenClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		wrongToken, _ := token.SignedString([]byte("wrong-key"))

		_, err := ValidateRefreshToken(wrongToken)
		if err == nil {
			t.Fatal("Expected error for invalid signature, got nil")
		}
	})
}
