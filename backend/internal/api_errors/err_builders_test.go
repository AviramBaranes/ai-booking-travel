package api_errors

import (
	"testing"

	"encore.dev/beta/errs"
)

func TestErrBuilders(t *testing.T) {
	t.Run("NewError", func(t *testing.T) {
		err := NewError(errs.InvalidArgument, "Invalid input")
		e, ok := err.(*errs.Error)
		if !ok {
			t.Fatalf("Expected *errs.Error, got %T", err)
		}
		if e.Code != errs.InvalidArgument {
			t.Errorf("Expected code %v, got %v", errs.InvalidArgument, e.Code)
		}
		if e.Message != "Invalid input" {
			t.Errorf("Expected message 'Invalid input', got '%s'", e.Message)
		}
	})

	t.Run("NewErrorWithDetail", func(t *testing.T) {
		detail := ErrorDetails{Code: "custom_error", Field: "username"}
		err := NewErrorWithDetail(errs.InvalidArgument, "Invalid username", detail)
		e, ok := err.(*errs.Error)
		if !ok {
			t.Fatalf("Expected *errs.Error,	 got %T", err)
		}
		if e.Code != errs.InvalidArgument {
			t.Errorf("Expected code %v, got %v", errs.InvalidArgument, e.Code)
		}
		if e.Message != "Invalid username" {
			t.Errorf("Expected message 'Invalid username', got '%s'", e.Message)
		}

		errDetails := errs.Details(err)
		d := errDetails.(ErrorDetails)
		if d.Code != "custom_error" {
			t.Errorf("Expected detail code 'custom_error', got '%s'", d.Code)
		}
		if d.Field != "username" {
			t.Errorf("Expected detail field 'username', got '%s'", d.Field)
		}
	})

	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError("Invalid email")
		e, ok := err.(*errs.Error)
		if !ok {
			t.Fatalf("Expected *errs.Error, got %T", err)
		}
		if e.Code != errs.InvalidArgument {
			t.Errorf("Expected code %v, got %v", errs.InvalidArgument, e.Code)
		}
		if e.Message != "Invalid email" {
			t.Errorf("Expected message 'Invalid email', got '%s'", e.Message)
		}
	})
}
