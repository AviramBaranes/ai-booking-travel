package accounts

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	ah "encore.app/services/accounts/handlers/auth"
	user "encore.app/services/accounts/handlers/user"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

func invalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// refreshParams builds RefreshTokensParams carrying the given refresh token in
// the Cookie header string, mirroring how a browser sends it.
func refreshParams(token string) ah.RefreshTokensParams {
	return ah.RefreshTokensParams{CookieHeader: "refresh_token=" + token}
}

// refreshCookieToken extracts the refresh_token value from a login/refresh
// response's Set-Cookie headers.
func refreshCookieToken(t *testing.T, resp *ah.LoginResponse) string {
	t.Helper()
	for _, sc := range resp.SetCookies {
		c, err := http.ParseSetCookie(sc)
		if err != nil {
			continue
		}
		if c.Name == "refresh_token" {
			return c.Value
		}
	}
	t.Fatal("refresh_token cookie not found in response Set-Cookie headers")
	return ""
}

const (
	allowedDelta = time.Second
)

var (
	pgxdb = sqldb.Driver[*pgxpool.Pool](accountsDb)
	query = db.New(pgxdb)
)

func newService(newDb *sqldb.Database) *Service {
	return &Service{
		query: db.New(sqldb.Driver[*pgxpool.Pool](newDb)),
	}
}

func mockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

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
func assertRefreshClaims(t *testing.T, claims *jwt.RefreshTokenClaims, userID int64) {
	t.Helper()
	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Issuer != jwt.Issuer {
		t.Errorf("Expected Issuer %s, got %s", jwt.Issuer, claims.Issuer)
	}
	expectedSub := strconv.Itoa(int(userID))
	if claims.Subject != expectedSub {
		t.Errorf("Expected Subject %s, got %s", expectedSub, claims.Subject)
	}
	if claims.ID == "" {
		t.Error("Expected non-empty JTI")
	}
}

// generateTestEmail creates a unique email for each test run.
func generateTestEmail() string {
	return fmt.Sprintf("loginuser_%d_%d@example.com", time.Now().UnixNano(), emailCounter.Add(1))
}

// assertAccessClaims verifies core access token claims.
func assertAccessClaims(t *testing.T, claims *jwt.AccessTokenClaims, expectedData jwt.AccessTokenData) {
	t.Helper()
	if claims.UserID != expectedData.UserID {
		t.Errorf("Expected UserID %d, got %d", expectedData.UserID, claims.UserID)
	}
	if string(claims.Role) != string(expectedData.Role) {
		t.Errorf("Expected Role %s, got %s", expectedData.Role, claims.Role)
	}

	if claims.OrganizationContext != nil && expectedData.OrganizationContext != nil {
		if claims.OrganizationContext.OfficeID != expectedData.OrganizationContext.OfficeID {
			t.Errorf("Expected OfficeID %d, got %d", expectedData.OrganizationContext.OfficeID, claims.OrganizationContext.OfficeID)
		}
		if claims.OrganizationContext.OrganizationID != expectedData.OrganizationContext.OrganizationID {
			t.Errorf("Expected OrganizationID %d, got %d", expectedData.OrganizationContext.OrganizationID, claims.OrganizationContext.OrganizationID)
		}
		if claims.OrganizationContext.IsOrganic != expectedData.OrganizationContext.IsOrganic {
			t.Errorf("Expected IsOrganic %v, got %v", expectedData.OrganizationContext.IsOrganic, claims.OrganizationContext.IsOrganic)
		}
	} else if claims.OrganizationContext != nil || expectedData.OrganizationContext != nil {
		t.Error("Expected both OrganizationContext to be nil or non-nil")
	}

	if claims.Issuer != jwt.Issuer {
		t.Errorf("Expected Issuer %s, got %s", jwt.Issuer, claims.Issuer)
	}
	expectedSub := strconv.Itoa(int(expectedData.UserID))
	if claims.Subject != expectedSub {
		t.Errorf("Expected Subject %s, got %s", expectedSub, claims.Subject)
	}
}

func registerAdmin(ctx context.Context, email, password string) (*user.CreateAdminResponse, func(), error) {
	admin, err := CreateAdmin(ctx, user.CreateAdminParams{
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
var phoneCounter atomic.Int64
var emailCounter atomic.Int64

func randomName() string {
	return fmt.Sprintf("name_%d_%d", time.Now().UnixNano(), nameCounter.Add(1))
}

func randomIsraeliPhoneNumber() string {
	return fmt.Sprintf("05%08d", phoneCounter.Add(1)%100000000)
}

func createAgent(ctx context.Context, p user.CreateAgentParams) (*user.CreateAgentResponse, func(), error) {
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
