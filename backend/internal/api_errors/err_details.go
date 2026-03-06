package api_errors

import "encore.dev/beta/errs"

type ErrorDetails struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
}

func (_ ErrorDetails) ErrDetails() {}

const (
	CodeInternalError   = "internal_error"
	CodeInvalidValue    = "invalid_value"
	CodeUnauthorized    = "unauthorized"
	CodeUnauthenticated = "unauthenticated"
	CodeNotFound        = "not_found"
)

var (
	InternalErrorDetails   = ErrorDetails{Code: CodeInternalError}
	InvalidValueDetails    = ErrorDetails{Code: CodeInvalidValue}
	UnauthorizedDetails    = ErrorDetails{Code: CodeUnauthorized}
	UnauthenticatedDetails = ErrorDetails{Code: CodeUnauthenticated}
	NotFoundDetails        = ErrorDetails{Code: CodeNotFound}
)

// WithDetail returns a brand-new *errs.Error error cloned from err
// with its Details set to detail, leaving the original untouched.
func WithDetail(err error, detail ErrorDetails) error {
	orig, ok := err.(*errs.Error)
	if !ok {
		return nil
	}
	b := errs.B().
		Code(orig.Code).
		Msg(orig.Message)

	return b.
		Details(detail).
		Err()
}
