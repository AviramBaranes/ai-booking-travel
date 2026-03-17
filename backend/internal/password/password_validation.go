package password

import (
	"errors"
	"unicode"
)

const (
	pwMinLength = 8
	pwMinUpper  = 1
	pwMinLower  = 1
	pwMinNumber = 1
	pwMinSymbol = 1
)

var (
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordNoUpper  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower  = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber = errors.New("password must contain at least one number")
	ErrPasswordNoSymbol = errors.New("password must contain at least one symbol")
)

// ValidatePassword checks if the provided password meets the security requirements.
func ValidatePassword(pw string) error {
	if len(pw) < pwMinLength {
		return ErrPasswordTooShort
	}
	var upp, low, num, sym int
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			upp++
		case unicode.IsLower(r):
			low++
		case unicode.IsNumber(r):
			num++
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			sym++
		}
	}
	if upp < pwMinUpper {
		return ErrPasswordNoUpper
	}
	if low < pwMinLower {
		return ErrPasswordNoLower
	}
	if num < pwMinNumber {
		return ErrPasswordNoNumber
	}
	if sym < pwMinSymbol {
		return ErrPasswordNoSymbol
	}
	return nil
}
