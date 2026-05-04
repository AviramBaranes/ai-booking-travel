package api_errors

const (
	CodePasswordTooShort = "password_too_short"
	CodePasswordNoUpper  = "password_no_upper"
	CodePasswordNoLower  = "password_no_lower"
	CodePasswordNoNumber = "password_no_number"
	CodePasswordNoSymbol = "password_no_symbol"

	CodeUserNotFound        = "user_not_found"
	CodeEmailAlreadyExists  = "email_already_exists"
	CodePhoneAlreadyExists  = "phone_already_exists"
	CodeInvalidCredentials  = "invalid_credentials"
	CodeInvalidRefreshToken = "invalid_refresh_token"
	CodeExpiredRefreshToken = "expired_refresh_token"
	CodeInvalidResetToken   = "invalid_reset_token"

	CodeOrganicOrgRequiresIcountClientID   = "organization_organic_requires_icount_client_id"
	CodeNonOrganicOrgForbidsIcountClientID = "organization_non_organic_forbids_icount_client_id"

	CodeOfficeOrganicForbidsIcountClientID    = "office_organic_forbids_icount_client_id"
	CodeOfficeNonOrganicRequiresIcountClientID = "office_non_organic_requires_icount_client_id"
)
