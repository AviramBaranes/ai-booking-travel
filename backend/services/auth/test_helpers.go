package auth

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"encore.app/services/auth/db"
	"encore.app/services/auth/jwt"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	allowedDelta = time.Second
	testEmailTpl = "loginuser_%d@example.com"
)

var (
	pgxdb = sqldb.Driver[*pgxpool.Pool](usersDB)
	query = db.New(pgxdb)
)

// assertTimeAlmostEqual checks if two time.Time values are within an acceptable delta.
func assertTimeAlmostEqual(t *testing.T, got, want time.Time) {
	t.Helper()
	diff := got.Sub(want)
	if diff > allowedDelta || diff < -allowedDelta {
		t.Errorf("Times differ too much: got %v; want %v (±%v), diff=%v",
			got, want, allowedDelta, diff)
	}
}

// assertRefreshClaims verifies core refresh token claims.
func assertRefreshClaims(t *testing.T, claims *jwt.RefreshTokenClaims, user *db.User) {
	t.Helper()
	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if claims.Issuer != jwt.Issuer {
		t.Errorf("Expected Issuer %s, got %s", jwt.Issuer, claims.Issuer)
	}
	expectedSub := strconv.Itoa(int(user.ID))
	if claims.Subject != expectedSub {
		t.Errorf("Expected Subject %s, got %s", expectedSub, claims.Subject)
	}
	if claims.ID == "" {
		t.Error("Expected non-empty JTI")
	}
}

// generateTestEmail creates a unique email for each test run.
func generateTestEmail() string {
	return fmt.Sprintf(testEmailTpl, time.Now().UnixNano())
}

// assertAccessClaims verifies core access token claims.
func assertAccessClaims(t *testing.T, claims *jwt.AccessTokenClaims, user *db.User) {
	t.Helper()
	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if string(claims.Role) != string(user.Role) {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
	if claims.Issuer != jwt.Issuer {
		t.Errorf("Expected Issuer %s, got %s", jwt.Issuer, claims.Issuer)
	}
	expectedSub := strconv.Itoa(int(user.ID))
	if claims.Subject != expectedSub {
		t.Errorf("Expected Subject %s, got %s", expectedSub, claims.Subject)
	}
}

func registerAdmin(ctx context.Context, username, password string) (*RegisterAdminResponse, func(), error) {
	admin, err := RegisterAdmin(ctx, RegisterAdminParams{
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, nil, err
	}

	return admin, func() {
		query.DeleteUser(ctx, admin.ID)
	}, nil
}

func registerAgent(ctx context.Context, p RegisterAgentParams) (*RegisterAgentResponse, func(), error) {
	agent, err := RegisterAgent(ctx, p)

	if err != nil {
		return nil, nil, err
	}

	return agent, func() {
		query.DeleteUser(ctx, agent.ID)
	}, nil
}
