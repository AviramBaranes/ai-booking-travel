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

func TestValidateStruct_Phone(t *testing.T) {
	type TestStruct struct {
		Phone string `json:"phone" validate:"required,phone"`
	}

	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"israeli local", "0521234567", true},
		{"israeli e164", "+972521234567", true},
		{"us e164", "+14155552671", true},
		{"uk e164", "+442071838750", true},
		{"empty", "", false},
		{"israeli local not starting with 05", "0321234567", false},
		{"israeli local too short", "052123456", false},
		{"international without plus", "972521234567", false},
		{"international with leading zero after plus", "+0521234567", false},
		{"international with dashes", "+972-52-1234567", false},
		{"international with spaces", "+972 52 1234567", false},
		{"international too short", "+97252", false},
		{"international too long", "+9725212345678901", false},
		{"non-digits", "+972abc4567", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(TestStruct{Phone: tt.phone})
			if tt.valid {
				if err != nil {
					t.Errorf("ValidateStruct() unexpected error for %q: %v", tt.phone, err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateStruct() expected error for %q, got nil", tt.phone)
			}
			api_errors.AssertApiError(t, api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg,
				api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue, Field: "phone"}), err)
		})
	}
}
