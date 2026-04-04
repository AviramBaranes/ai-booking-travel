package validation

import (
	"testing"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" validate:"required,notblank"`
		Phone string `json:"phone" validate:"required,israeli_phone"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr error
	}{
		{"valid input", TestStruct{Name: "John Doe", Phone: "0521234567"}, nil},
		{"empty name", TestStruct{Name: "", Phone: "0521234567"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "name"})},
		{"name with only spaces", TestStruct{Name: "   ", Phone: "0521234567"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "name"})},
		{"empty phone", TestStruct{Name: "John Doe", Phone: ""},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"})},
		{"phone not starting with 05", TestStruct{Name: "John Doe", Phone: "0321234567"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"})},
		{"phone too short", TestStruct{Name: "John Doe", Phone: "052123456"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"})},
		{"phone too long", TestStruct{Name: "John Doe", Phone: "05212345678"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"})},
		{"phone with non-digits", TestStruct{Name: "John Doe", Phone: "052-123456"},
			api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"})},
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
