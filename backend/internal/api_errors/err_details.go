package api_errors

type ErrorDetails struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
}

func (_ ErrorDetails) ErrDetails() {}

var (
	InternalErrorDetails   = ErrorDetails{Code: CodeInternalError}
	InvalidValueDetails    = ErrorDetails{Code: CodeInvalidValue}
	UnauthorizedDetails    = ErrorDetails{Code: CodeUnauthorized}
	UnauthenticatedDetails = ErrorDetails{Code: CodeUnauthenticated}
	NotFoundDetails        = ErrorDetails{Code: CodeNotFound}
)
