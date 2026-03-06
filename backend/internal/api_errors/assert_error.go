package api_errors

import (
	"errors"
	"reflect"
	"testing"

	"encore.dev/beta/errs"
)

// AssertApiError compares two Encore *errs.Error values.
func AssertApiError(t *testing.T, want, got error) {
	t.Helper()

	var we, ge *errs.Error
	if !errors.As(want, &we) {
		t.Fatalf("want is not *errs.Error: %T", want)
	}
	if !errors.As(got, &ge) {
		t.Fatalf("got is not *errs.Error: %v (%T)", got, got)
	}

	if ge.Code != we.Code {
		t.Fatalf("error code mismatch: want %q, got %q", we.Code, ge.Code)
	}
	if ge.Message != we.Message {
		t.Fatalf("message mismatch:\n want: %q\n  got: %q", we.Message, ge.Message)
	}

	if !reflect.DeepEqual(ge.Details, we.Details) {
		t.Fatalf("details mismatch:\n want: %v\n  got: %v", we.Details, ge.Details)
	}
}
