package accounts

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func invalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

const (
	allowedDelta = time.Second
	testEmailTpl = "loginuser_%d@example.com"
)

var (
	pgxdb = sqldb.Driver[*pgxpool.Pool](accountsDb)
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

func registerAdmin(ctx context.Context, email, password string) (*CreateAdminResponse, func(), error) {
	admin, err := CreateAdmin(ctx, CreateAdminRequest{
		FirstName: "Test",
		LastName:  "Admin",
		Email:     email,
		Password:  password,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("registering admin %w", err)
	}

	return admin, func() {
		query.DeleteUser(ctx, admin.ID)
	}, nil
}

var nameCounter atomic.Int64

func randomName() string {
	return fmt.Sprintf("name_%d_%d", time.Now().UnixNano(), nameCounter.Add(1))
}

func randomIsraeliPhoneNumber() string {
	return fmt.Sprintf("05%08d", time.Now().UnixNano()%100000000)
}

func createAgent(ctx context.Context, p CreateAgentRequest) (*CreateAgentResponse, func(), error) {
	org, err := query.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:      randomName(),
		IsOrganic: false,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("creating organization %w", err)
	}

	office, err := query.CreateOffice(ctx, db.CreateOfficeParams{
		Name:           randomName(),
		OrganizationID: org.ID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("creating office %w", err)
	}

	p.OfficeID = office.ID
	agent, err := CreateAgent(ctx, p)

	if err != nil {
		return nil, nil, err
	}

	return agent, func() {
		query.DeleteUser(ctx, agent.ID)
	}, nil
}
