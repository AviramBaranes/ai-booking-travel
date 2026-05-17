package user_handlers

import (
	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
)

// UserService holds the handler logic for all user management endpoints.
type UserService struct {
	query db.Querier
}

// NewUserService creates a new UserService.
func NewUserService(query db.Querier) *UserService {
	return &UserService{query: query}
}

// GetAgentsResponse is the shared response type for agent ID lookups.
type GetAgentsResponse struct {
	IDs       []int64
	IsOrganic bool
}

var (
	ErrUserNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound, "User not found",
		api_errors.ErrorDetails{Code: api_errors.CodeUserNotFound},
	)

	ErrEmailAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Email already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeEmailAlreadyExists},
	)

	ErrPhoneAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Phone number already exists",
		api_errors.ErrorDetails{Code: api_errors.CodePhoneAlreadyExists},
	)
)
