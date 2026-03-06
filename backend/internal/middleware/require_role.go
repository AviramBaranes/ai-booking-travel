package middleware

import (
	"encore.app/internal/api_errors"
	a "encore.app/services/auth"
	"encore.app/services/auth/db"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
)

// RequireRoleMiddleware is a middleware that checks if the user has the required role.
func RequireRoleMiddleware(role db.UserRole, req middleware.Request, next middleware.Next) middleware.Response {
	data, ok := auth.Data().(*a.AuthData)
	if !ok {
		return middleware.Response{
			Err: api_errors.ErrUnauthorized,
		}
	}

	if data.Role != role {
		return middleware.Response{
			Err: api_errors.ErrUnauthorized,
		}
	}

	return next(req)
}

// RequireAdminMiddleware is a middleware that checks if the user has the admin role.
// encore:middleware global target=tag:admin
func RequireAdminMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware(db.UserRoleAdmin, req, next)
}

// RequireAgentMiddleware is a middleware that checks if the user has the agent role.
// encore:middleware global target=tag:agent
func RequireAgentMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware(db.UserRoleAgent, req, next)
}

// RequireCustomerMiddleware is a middleware that checks if the user has the customer role.
// encore:middleware global target=tag:customer
func RequireCustomerMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware(db.UserRoleCustomer, req, next)
}
