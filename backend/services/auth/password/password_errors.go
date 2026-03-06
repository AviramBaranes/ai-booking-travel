package password

import (
	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

const (
	CodePasswordTooShort = "password_too_short"
	CodePasswordNoUpper  = "password_no_upper"
	CodePasswordNoLower  = "password_no_lower"
	CodePasswordNoNumber = "password_no_number"
	CodePasswordNoSymbol = "password_no_symbol"
)

var (
	passwordTooShortDetails = api_errors.ErrorDetails{
		Code:  CodePasswordTooShort,
		Field: "password",
	}
	passwordNoUpperDetails = api_errors.ErrorDetails{
		Code:  CodePasswordNoUpper,
		Field: "password",
	}
	passwordNoLowerDetails = api_errors.ErrorDetails{
		Code:  CodePasswordNoLower,
		Field: "password",
	}
	passwordNoNumberDetails = api_errors.ErrorDetails{
		Code:  CodePasswordNoNumber,
		Field: "password",
	}
	passwordNoSymbolDetails = api_errors.ErrorDetails{
		Code:  CodePasswordNoSymbol,
		Field: "password",
	}
)

var (
	errPasswordTooShort = errs.B().
				Code(errs.InvalidArgument).
				Details(passwordTooShortDetails).
				Msg("Password must be at least 8 characters long").
				Err()

	errPasswordNoUpper = errs.B().
				Code(errs.InvalidArgument).
				Details(passwordNoUpperDetails).
				Msg("Password must contain at least one uppercase letter").
				Err()

	errPasswordNoLower = errs.B().
				Code(errs.InvalidArgument).
				Details(passwordNoLowerDetails).
				Msg("Password must contain at least one lowercase letter").
				Err()

	errPasswordNoNumber = errs.B().
				Code(errs.InvalidArgument).
				Details(passwordNoNumberDetails).
				Msg("Password must contain at least one number").
				Err()

	errPasswordNoSymbol = errs.B().
				Code(errs.InvalidArgument).
				Details(passwordNoSymbolDetails).
				Msg("Password must contain at least one symbol").
				Err()
)
