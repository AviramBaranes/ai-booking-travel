package api_errors

import "encore.dev/beta/errs"

var (
	ErrInternalError = errs.B().
				Code(errs.Internal).
				Msg("Internal server error").
				Details(InternalErrorDetails).
				Err()

	ErrInvalidValue = errs.B().
			Code(errs.InvalidArgument).
			Msg("Invalid value provided").
			Details(InvalidValueDetails).
			Err()

	ErrUnauthorized = errs.B().
			Code(errs.PermissionDenied).
			Msg("Unauthorized access").
			Details(UnauthorizedDetails).
			Err()

	ErrUnauthenticated = errs.B().
				Code(errs.Unauthenticated).
				Msg("Unauthenticated request").
				Details(UnauthenticatedDetails).
				Err()

	ErrNotFound = errs.B().
			Code(errs.NotFound).
			Msg("Resource not found").
			Details(NotFoundDetails).
			Err()
)
