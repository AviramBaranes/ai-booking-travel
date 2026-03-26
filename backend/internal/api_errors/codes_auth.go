package api_errors

const (
	CodePasswordTooShort = "password_too_short"
	CodePasswordNoUpper  = "password_no_upper"
	CodePasswordNoLower  = "password_no_lower"
	CodePasswordNoNumber = "password_no_number"
	CodePasswordNoSymbol = "password_no_symbol"

	CodeUserNotFound          = "user_not_found"
	CodeUsernameAlreadyExists = "username_already_exists"
	CodeInvalidCredentials    = "invalid_credentials"
	CodeInvalidRefreshToken   = "invalid_refresh_token"
	CodeExpiredRefreshToken   = "expired_refresh_token"
	CodeInvalidResetToken     = "invalid_reset_token"
)
