package validation

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	v "github.com/go-playground/validator/v10"
)

var validator = v.New()
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._-]{1,30}[a-zA-Z0-9])?$`)

func init() {
	validator.RegisterValidation("notblank", func(fl v.FieldLevel) bool {
		s := fl.Field().String()
		return strings.TrimSpace(s) != ""
	})

	validator.RegisterValidation("username", func(fl v.FieldLevel) bool {
		s := fl.Field().String()

		if len(s) > 31 {
			return false
		}

		return usernameRegex.MatchString(s)
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

// ValidateStruct validates the provided struct using the validator package.
func ValidateStruct(p any) error {
	if err := validator.Struct(p); err != nil {
		var ves v.ValidationErrors
		if errors.As(err, &ves) {
			jsonField := getFieldName(p, ves[0].StructField())

			detail := api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: jsonField,
			}

			return errs.B().
				Code(errs.InvalidArgument).
				Details(detail).
				Msg("Invalid value provided").
				Err()
		}
		// non‐validation error: use your generic one
		return api_errors.ErrInvalidValue
	}
	return nil
}

func NewValidationError(msg string) error {
	return errs.B().
		Code(errs.InvalidArgument).
		Msg(msg).
		Err()
}
