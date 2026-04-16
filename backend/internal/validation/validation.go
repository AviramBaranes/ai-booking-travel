package validation

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	v "github.com/go-playground/validator/v10"
)

var validator = v.New()
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._-]{1,30}[a-zA-Z0-9])?$`)
var israeliPhoneRegex = regexp.MustCompile(`^05\d{8}$`)

func init() {
	validator.RegisterValidation("notblank", func(fl v.FieldLevel) bool {
		s := fl.Field().String()
		return strings.TrimSpace(s) != ""
	})

	// israeli_phone accepts digits only, must start with 05 followed by exactly 8 digits (e.g. 0521234567).
	validator.RegisterValidation("israeli_phone", func(fl v.FieldLevel) bool {
		return israeliPhoneRegex.MatchString(fl.Field().String())
	})

	validator.RegisterValidation("uppercase_only", func(fl v.FieldLevel) bool {
		s := fl.Field().String()
		for _, r := range s {
			if r < 'A' || r > 'Z' {
				return false
			}
		}

		return len(s) > 0
	})
}

// getFieldName returns the first part of the `json:"..."` tag for the given struct field.
// If there is no json tag or it’s "-", it falls back to the `query` tag.
// If that’s also missing or “-”, it finally falls back to the Go field name.
func getFieldName(p any, goField string) string {
	t := reflect.TypeOf(p)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if f, ok := t.FieldByName(goField); ok {
		// try json tag first
		if tag := strings.Split(f.Tag.Get("json"), ",")[0]; tag != "" && tag != "-" {
			return tag
		}
		// fallback to query tag
		if tag := strings.Split(f.Tag.Get("query"), ",")[0]; tag != "" && tag != "-" {
			return tag
		}
	}
	return goField
}

const (
	// InvalidValueMsg is the default error message for invalid values in validation errors. Used for tests assertions.
	InvalidValueMsg = "Invalid value provided"
)

// ValidateStruct validates the provided struct using the validator package.
func ValidateStruct(p any) error {
	if err := validator.Struct(p); err != nil {
		var ves v.ValidationErrors
		if errors.As(err, &ves) {
			jsonField := getFieldName(p, ves[0].StructField())

			return api_errors.NewErrorWithDetail(errs.InvalidArgument, InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: jsonField,
			})
		}

		rlog.Error("Validation error", "error", err.Error())
		return api_errors.ErrInvalidValue
	}
	return nil
}
