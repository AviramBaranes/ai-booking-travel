package validation

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

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

	// uppercase_only accepts A-Z and spaces only (e.g. "BEN DAVID"), and requires at
	// least one letter so a blank/whitespace-only value never passes.
	validator.RegisterValidation("uppercase_only", func(fl v.FieldLevel) bool {
		s := fl.Field().String()
		hasLetter := false
		for _, r := range s {
			switch {
			case r >= 'A' && r <= 'Z':
				hasLetter = true
			case r == ' ':
			default:
				return false
			}
		}

		return hasLetter
	})

	// person_name accepts letters in any script (Hebrew included) plus spaces, hyphens
	// and apostrophes (e.g. "בן דוד", "Ben-David"), and requires at least one letter.
	validator.RegisterValidation("person_name", func(fl v.FieldLevel) bool {
		s := fl.Field().String()
		hasLetter := false
		for _, r := range s {
			switch {
			case unicode.IsLetter(r):
				hasLetter = true
			case r == ' ' || r == '-' || r == '\'':
			default:
				return false
			}
		}

		return hasLetter
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

			msg := fmt.Sprintf("Validation failed on field '%s' with value '%s'", jsonField, ves[0].Value())
			rlog.Debug("validation error", "msg", msg)

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
