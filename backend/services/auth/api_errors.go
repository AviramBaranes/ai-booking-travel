package auth

import (
	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

var (
	ErrUserNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound, "User not found",
		api_errors.ErrorDetails{Code: api_errors.CodeUserNotFound},
	)

	ErrUsernameAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Username already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeUsernameAlreadyExists},
	)

	ErrInvalidCredentials = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Invalid credentials",
		api_errors.ErrorDetails{Code: api_errors.CodeInvalidCredentials},
	)

	ErrInvalidRefreshToken = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Invalid refresh token",
		api_errors.ErrorDetails{Code: api_errors.CodeInvalidRefreshToken},
	)

	ErrExpiredRefreshToken = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Expired refresh token",
		api_errors.ErrorDetails{Code: api_errors.CodeExpiredRefreshToken},
	)

	ErrInvalidResetToken = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Invalid reset token",
		api_errors.ErrorDetails{Code: api_errors.CodeInvalidResetToken},
	)
)
