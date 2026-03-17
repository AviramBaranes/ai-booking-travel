package auth

import (
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

func newPasswordValidationError(code, msg string) error {
	return errs.B().
		Code(errs.InvalidArgument).
		Details(api_errors.ErrorDetails{
			Code:  code,
			Field: "password",
		}).
		Msg(msg).
		Err()
}

var (
	ErrPasswordTooShort = newPasswordValidationError(
		api_errors.CodePasswordTooShort,
		"Password must be at least 8 characters long",
	)
	ErrPasswordNoUpper = newPasswordValidationError(
		api_errors.CodePasswordNoUpper,
		"Password must contain at least one uppercase letter",
	)
	ErrPasswordNoLower = newPasswordValidationError(
		api_errors.CodePasswordNoLower,
		"Password must contain at least one lowercase letter",
	)
	ErrPasswordNoNumber = newPasswordValidationError(
		api_errors.CodePasswordNoNumber,
		"Password must contain at least one number",
	)
	ErrPasswordNoSymbol = newPasswordValidationError(
		api_errors.CodePasswordNoSymbol,
		"Password must contain at least one symbol",
	)
)

func validatePasswordForAPI(pass string) error {
	if err := password.ValidatePassword(pass); err != nil {
		switch {
		case errors.Is(err, password.ErrPasswordTooShort):
			return ErrPasswordTooShort
		case errors.Is(err, password.ErrPasswordNoUpper):
			return ErrPasswordNoUpper
		case errors.Is(err, password.ErrPasswordNoLower):
			return ErrPasswordNoLower
		case errors.Is(err, password.ErrPasswordNoNumber):
			return ErrPasswordNoNumber
		case errors.Is(err, password.ErrPasswordNoSymbol):
			return ErrPasswordNoSymbol
		default:
			rlog.Error("failed to validate password", "error", err)
			return api_errors.ErrInternalError
		}
	}
	return nil
}
