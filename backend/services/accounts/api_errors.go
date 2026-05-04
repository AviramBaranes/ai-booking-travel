package accounts

import (
	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

var (
	ErrUserNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound, "User not found",
		api_errors.ErrorDetails{Code: api_errors.CodeUserNotFound},
	)

	ErrEmailAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Email already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeEmailAlreadyExists},
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

	ErrNameAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Name already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeNameAlreadyExists},
	)

	ErrPhoneAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Phone number already exists",
		api_errors.ErrorDetails{Code: api_errors.CodePhoneAlreadyExists},
	)

	ErrOrganizationOrganicRequiresIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Organic organization must have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeOrganicOrgRequiresIcountClientID, Field: "icountClientId"},
	)

	ErrOrganizationNonOrganicForbidsIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Non-organic organization must not have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeNonOrganicOrgForbidsIcountClientID, Field: "icountClientId"},
	)
)
