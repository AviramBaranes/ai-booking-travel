package middleware

import (
	"encore.app/internal/api_errors"
	a "encore.app/services/accounts"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
)

// RequireRoleMiddleware is a middleware that checks if the user has the one of the required roles.
func RequireRoleMiddleware(allowedRoles []a.UserRole, req middleware.Request, next middleware.Next) middleware.Response {
	data, ok := auth.Data().(*a.AuthData)
	if !ok {
		return middleware.Response{
			Err: api_errors.ErrUnauthorized,
		}
	}

	for _, allowedRole := range allowedRoles {
		if data.Role == allowedRole {
			return next(req)
		}
	}

	return middleware.Response{
		Err: api_errors.ErrUnauthorized,
	}
}

// RequireAdminMiddleware is a middleware that checks if the user has the admin role.
// encore:middleware global target=tag:admin
func RequireAdminMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware([]a.UserRole{a.UserRoleAdmin}, req, next)
}

// RequireAgentMiddleware is a middleware that checks if the user has the agent role.
// encore:middleware global target=tag:agent
func RequireAgentMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware([]a.UserRole{a.UserRoleAgent}, req, next)
}

// RequireCustomerMiddleware is a middleware that checks if the user has the customer role.
// encore:middleware global target=tag:customer
func RequireCustomerMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware([]a.UserRole{a.UserRoleCustomer}, req, next)
}

// RequireAccountantMiddleware is a middleware that checks if the user has the accountant role.
// encore:middleware global target=tag:accountant
func RequireAccountantMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware([]a.UserRole{a.UserRoleAccountant, a.UserRoleAdmin}, req, next)
}

// RequireAgentOrCustomerMiddleware allows both agents and customers
// encore:middleware global target=tag:agent_customer
func RequireAgentOrCustomerMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	return RequireRoleMiddleware([]a.UserRole{a.UserRoleAgent, a.UserRoleCustomer}, req, next)
}
