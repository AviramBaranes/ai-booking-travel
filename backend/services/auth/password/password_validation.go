package password

import (
	"unicode"
)

const (
	pwMinLength = 8
	pwMinUpper  = 1
	pwMinLower  = 1
	pwMinNumber = 1
	pwMinSymbol = 1
)

// ValidatePassword checks if the provided password meets the security requirements.
func ValidatePassword(pw string) error {
	if len(pw) < pwMinLength {
		return errPasswordTooShort
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
		return errPasswordNoUpper
	}
	if low < pwMinLower {
		return errPasswordNoLower
	}
	if num < pwMinNumber {
		return errPasswordNoNumber
	}
	if sym < pwMinSymbol {
		return errPasswordNoSymbol
	}
	return nil
}
