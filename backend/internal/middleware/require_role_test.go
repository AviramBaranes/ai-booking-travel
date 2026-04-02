package middleware

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	encore "encore.dev"
	"encore.dev/et"
	mw "encore.dev/middleware"
)

func runRequireRoleTests(t *testing.T, mwFunc func(mw.Request, mw.Next) mw.Response, role accounts.UserRole) {
	ctx := context.Background()
	req := mw.NewRequest(ctx, &encore.Request{})

	t.Run("No Auth", func(t *testing.T) {
		nextCalled := false
		resp := mwFunc(req, func(r mw.Request) mw.Response {
			nextCalled = true
			return mw.Response{}
		})
		if nextCalled {
			t.Error("next should not be called")
		}
		api_errors.AssertApiError(t, api_errors.ErrUnauthorized, resp.Err)
	})

	t.Run("Wrong Role", func(t *testing.T) {
		nextCalled := false
		otherRole := accounts.UserRoleAdmin
		if role == accounts.UserRoleAdmin {
			otherRole = accounts.UserRoleCustomer
		}
		et.OverrideAuthInfo("1", &accounts.AuthData{UserID: 1, Role: otherRole})
		defer et.OverrideAuthInfo("", nil)
		resp := mwFunc(req, func(r mw.Request) mw.Response {
			nextCalled = true
			return mw.Response{}
		})
		if nextCalled {
			t.Error("next should not be called")
		}
		api_errors.AssertApiError(t, api_errors.ErrUnauthorized, resp.Err)
	})

	t.Run("Correct Role", func(t *testing.T) {
		nextCalled := false
		et.OverrideAuthInfo("2", &accounts.AuthData{UserID: 2, Role: role})
		defer et.OverrideAuthInfo("", nil)
		resp := mwFunc(req, func(r mw.Request) mw.Response {
			nextCalled = true
			return mw.Response{}
		})
		if resp.Err != nil {
			t.Fatalf("unexpected error: %v", resp.Err)
		}
		if !nextCalled {
			t.Error("next was not called")
		}
	})
}

func TestRequireAdmin(t *testing.T) {
	runRequireRoleTests(t, RequireAdminMiddleware, accounts.UserRoleAdmin)
}

func TestRequireCustomer(t *testing.T) {
	runRequireRoleTests(t, RequireCustomerMiddleware, accounts.UserRoleCustomer)
}

func TestRequireAgent(t *testing.T) {
	runRequireRoleTests(t, RequireAgentMiddleware, accounts.UserRoleAgent)
}
