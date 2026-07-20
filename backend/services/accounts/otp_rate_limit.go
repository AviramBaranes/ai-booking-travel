package accounts

import (
	"context"
	"time"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/storage/cache"
)

const (
	otpMaxSendAttempts     = 3
	otpMaxValidateAttempts = 5
)

// otpSendLimiter tracks how many times an OTP has been sent to a phone number.
// Resets after 10 minutes.
var otpSendLimiter = cache.NewIntKeyspace[string](GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "otp-send/:key",
	DefaultExpiry: cache.ExpireIn(10 * time.Minute),
})

// otpValidateLimiter tracks how many failed validation attempts have been made for a phone number.
// Resets after 10 minutes.
var otpValidateLimiter = cache.NewIntKeyspace[string](GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "otp-validate/:key",
	DefaultExpiry: cache.ExpireIn(10 * time.Minute),
})

// Errors specific to OTP rate limiting.
var (
	errOTPSendRateLimited     = api_errors.NewErrorWithDetail(errs.ResourceExhausted, "Too many OTP send requests", api_errors.ErrorDetails{Code: api_errors.CodeOTPSendRateLimit})
	errOTPValidateRateLimited = api_errors.NewErrorWithDetail(errs.ResourceExhausted, "Too many OTP validation attempts", api_errors.ErrorDetails{Code: api_errors.CodeOTPValidateRateLimit})
)

// OTPRateLimiter abstracts OTP rate limiting to avoid an import cycle with the accounts package.
type OTPRateLimiter interface {
	CheckAndIncrementSend(ctx context.Context, phone string) error
	// RecordFailedAttempt increments the wrong-attempt counter.
	// Returns (true, nil) when the lockout threshold is reached; (false, nil) otherwise.
	// Returns (false, err) on unexpected cache errors.
	RecordFailedAttempt(ctx context.Context, phone string) (limitReached bool, err error)
	ClearAttempts(ctx context.Context, phone string)
	// ValidateRateLimitErr returns the error to use when validate lockout is triggered.
	ValidateRateLimitErr() error
}

// otpRateLimiter is the concrete implementation backed by the global cache.
type otpRateLimiter struct{}

func newOTPRateLimiter() OTPRateLimiter {
	return &otpRateLimiter{}
}

// CheckAndIncrementSend increments the send counter for the given phone number.
// Returns errOTPSendRateLimited if the limit has been exceeded.
func (r *otpRateLimiter) CheckAndIncrementSend(ctx context.Context, phone string) error {
	count, err := otpSendLimiter.Increment(ctx, phone, 1)
	if err != nil {
		return api_errors.ErrInternalError
	}
	if count > otpMaxSendAttempts {
		return errOTPSendRateLimited
	}
	return nil
}

// RecordFailedAttempt increments the failed validation counter for the given phone number.
// Returns limitReached=true when the threshold is hit so the caller can clear the OTP.
func (r *otpRateLimiter) RecordFailedAttempt(ctx context.Context, phone string) (limitReached bool, err error) {
	count, err := otpValidateLimiter.Increment(ctx, phone, 1)
	if err != nil {
		return false, api_errors.ErrInternalError
	}
	return count >= otpMaxValidateAttempts, nil
}

// ClearAttempts deletes the failed-validation counter for the given phone number (call on success).
func (r *otpRateLimiter) ClearAttempts(ctx context.Context, phone string) {
	// Best-effort; a stale counter expires naturally.
	_, _ = otpValidateLimiter.Delete(ctx, phone)
}

// ValidateRateLimitErr returns the error to surface when the validate lockout is triggered.
func (r *otpRateLimiter) ValidateRateLimitErr() error {
	return errOTPValidateRateLimited
}
