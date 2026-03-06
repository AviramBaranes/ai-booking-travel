package auth

import (
	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
)

var (
	UserNotFoundDetails          = api_errors.ErrorDetails{Code: "user_not_found"}
	UsernameAlreadyExistsDetails = api_errors.ErrorDetails{Code: "username_already_exists"}
	InvalidCredentialsDetails    = api_errors.ErrorDetails{Code: "invalid_credentials"}
	InvalidRefreshTokenDetails   = api_errors.ErrorDetails{Code: "invalid_refresh_token"}
	ExpiredRefreshTokenDetails   = api_errors.ErrorDetails{Code: "expired_refresh_token"}
	InvalidResetTokenDetails     = api_errors.ErrorDetails{Code: "invalid_reset_token"}
)

var (
	ErrUserNotFound = errs.B().
			Code(errs.NotFound).
			Details(UserNotFoundDetails).
			Msg("User not found").
			Err()

	ErrUsernameAlreadyExists = errs.B().
					Code(errs.AlreadyExists).
					Details(UsernameAlreadyExistsDetails).
					Msg("Username already exists").
					Err()

	ErrInvalidCredentials = errs.B().
				Code(errs.Unauthenticated).
				Details(InvalidCredentialsDetails).
				Msg("Invalid credentials").
				Err()

	ErrInvalidRefreshToken = errs.B().
				Code(errs.Unauthenticated).
				Details(InvalidRefreshTokenDetails).
				Msg("Invalid refresh token").
				Err()

	ErrExpiredRefreshToken = errs.B().
				Code(errs.Unauthenticated).
				Details(ExpiredRefreshTokenDetails).
				Msg("Expired refresh token").
				Err()

	ErrInvalidResetToken = errs.B().
				Code(errs.InvalidArgument).
				Details(InvalidResetTokenDetails).
				Msg("Invalid reset token").
				Err()
)
