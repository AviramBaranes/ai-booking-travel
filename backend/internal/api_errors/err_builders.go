package api_errors

import "encore.dev/beta/errs"

// NewError creates a new *errs.Error with the given code and message.
func NewError(code errs.ErrCode, msg string) error {
	return errs.B().
		Code(code).
		Msg(msg).
		Err()
}

// NewErrorWithDetail creates a new *errs.Error with the given code, message, and detail.
func NewErrorWithDetail(code errs.ErrCode, msg string, detail ErrorDetails) error {
	b := errs.B().
		Code(code).
		Msg(msg)

	return b.
		Details(detail).
		Err()
}

// NewValidationError creates a new *errs.Error with the InvalidArgument code and the given message.
func NewValidationError(msg string) error {
	return NewError(errs.InvalidArgument, msg)
}
