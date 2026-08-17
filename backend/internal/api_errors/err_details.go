package api_errors

import "encore.dev/beta/errs"

type ErrorDetails struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
}

func (_ ErrorDetails) ErrDetails() {}

// DetailsOf extracts our ErrorDetails from err. Details survive an internal RPC while
// the error message does not carry the Code or the offending Field, so this is the only
// way a caller can report which field a callee rejected — relevant for private endpoints,
// whose request validation runs before the handler and so produces no span of its own.
// Returns EmptyDetails when err has no ErrorDetails attached.
func DetailsOf(err error) ErrorDetails {
	if d, ok := errs.Details(err).(ErrorDetails); ok {
		return d
	}
	return EmptyDetails
}

var (
	InternalErrorDetails   = ErrorDetails{Code: CodeInternalError}
	InvalidValueDetails    = ErrorDetails{Code: CodeInvalidValue}
	UnauthorizedDetails    = ErrorDetails{Code: CodeUnauthorized}
	UnauthenticatedDetails = ErrorDetails{Code: CodeUnauthenticated}
	NotFoundDetails        = ErrorDetails{Code: CodeNotFound}
	EmptyDetails           = ErrorDetails{}
)
