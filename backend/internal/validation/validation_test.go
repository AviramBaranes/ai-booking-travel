package validation

import (
	"testing"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name     string `json:"name" validate:"required,notblank"`
		Username string `json:"username" validate:"required,username"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr error
	}{
		{"valid input", TestStruct{Name: "John Doe", Username: "john_doe"}, nil},
		{"empty name", TestStruct{Name: "", Username: "john_doe"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "name"})},
		{"name with only spaces", TestStruct{Name: "   ", Username: "john_doe"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "name"})},
		{"empty username", TestStruct{Name: "John Doe", Username: ""},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "username"})},
		{"username with invalid characters", TestStruct{Name: "John Doe", Username: "john@doe"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "username"})},
		{"username too long", TestStruct{Name: "John Doe", Username: "a_very_long_username_exceeding_31_characters"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "username"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateStruct() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateStruct() expected error, got nil")
			}
			api_errors.AssertApiError(t, tt.wantErr, err)
		})
	}
}
