package accounts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	emailevents "encore.app/internal/email_events"
	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	ah "encore.app/services/accounts/handlers/auth"
	user "encore.app/services/accounts/handlers/user"
	"encore.dev/et"
)

func TestRequestPasswordReset(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid email", func(t *testing.T) {
		err := (ah.SendPasswordResetTokenParams{Email: "not-an-email"}).Validate()
		api_errors.AssertApiError(t, invalidValueErr("email"), err)
	})

	t.Run("user not found does not reveal account existence", func(t *testing.T) {
		publishedBefore := len(et.Topic(emailevents.EmailRequestedTopic).PublishedMessages())

		err := RequestPasswordReset(ctx, ah.SendPasswordResetTokenParams{Email: generateTestEmail()})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		publishedAfter := len(et.Topic(emailevents.EmailRequestedTopic).PublishedMessages())
		if publishedAfter != publishedBefore {
			t.Fatalf("expected no email event, before=%d after=%d", publishedBefore, publishedAfter)
		}
	})

	t.Run("invalid user role", func(t *testing.T) {
		phone := randomIsraeliPhoneNumber()
		email := generateTestEmail()
		created, err := query.CreateCustomer(ctx, db.CreateCustomerParams{
			FirstName:    "Test",
			LastName:     "Customer",
			Email:        email,
			PhoneNumber:  &phone,
			PasswordHash: "test-password-hash",
		})
		if err != nil {
			t.Fatalf("failed to create customer: %v", err)
		}
		defer query.DeleteUser(ctx, created.ID)

		err = RequestPasswordReset(ctx, ah.SendPasswordResetTokenParams{Email: email})
		api_errors.AssertApiError(t, ah.ErrCustomerPasswordResetNotAllowed, err)
	})

	t.Run("success saves token and publishes email event", func(t *testing.T) {
		email := generateTestEmail()
		admin, err := CreateAdmin(ctx, user.CreateAdminParams{
			FirstName: "Test",
			LastName:  "Admin",
			Email:     email,
			Password:  testPassword,
		})
		if err != nil {
			t.Fatalf("failed to create admin: %v", err)
		}
		defer query.DeleteUser(ctx, admin.ID)

		publishedBefore := len(et.Topic(emailevents.EmailRequestedTopic).PublishedMessages())

		err = RequestPasswordReset(ctx, ah.SendPasswordResetTokenParams{Email: email})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		published := et.Topic(emailevents.EmailRequestedTopic).PublishedMessages()
		if len(published) != publishedBefore+1 {
			t.Fatalf("expected one email event, before=%d after=%d", publishedBefore, len(published))
		}

		last := published[len(published)-1]
		if last.Type != emailevents.EmailEventTypePasswordReset {
			t.Fatalf("expected password reset event, got %q", last.Type)
		}

		var payload emailevents.PasswordResetEmailPayload
		if err := json.Unmarshal(last.Payload, &payload); err != nil {
			t.Fatalf("failed to unmarshal payload: %v", err)
		}
		if payload.Email != email {
			t.Fatalf("expected email %q, got %q", email, payload.Email)
		}

		tokenHash := testPasswordResetHash(payload.TokenHash)
		saved, err := query.GetPasswordResetTokenByHash(ctx, tokenHash)
		if err != nil {
			t.Fatalf("failed to get saved reset token: %v", err)
		}
		if saved.UserID != admin.ID {
			t.Fatalf("expected token user ID %d, got %d", admin.ID, saved.UserID)
		}
	})
}

func TestResetPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("password and token validation", func(t *testing.T) {
		api_errors.AssertApiError(t, user.ErrPasswordTooShort, ah.ResetPasswordParams{
			Token:       "token",
			NewPassword: "short",
		}.Validate())

		api_errors.AssertApiError(t, invalidValueErr("token"), ah.ResetPasswordParams{
			Token:       "",
			NewPassword: testPassword,
		}.Validate())
	})

	t.Run("token not found", func(t *testing.T) {
		err := ResetPassword(ctx, ah.ResetPasswordParams{
			Token:       "missing-token",
			NewPassword: "NewPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("invalid expired token", func(t *testing.T) {
		admin, err := CreateAdmin(ctx, user.CreateAdminParams{
			FirstName: "Expired",
			LastName:  "Admin",
			Email:     generateTestEmail(),
			Password:  testPassword,
		})
		if err != nil {
			t.Fatalf("failed to create admin: %v", err)
		}
		defer query.DeleteUser(ctx, admin.ID)

		rawToken := "expired-reset-token"
		tokenHash := testPasswordResetHash(rawToken)
		created, err := query.InsertPasswordResetToken(ctx, db.InsertPasswordResetTokenParams{
			UserID:    admin.ID,
			TokenHash: tokenHash,
			ExpiresAt: dbadapters.DBTime(time.Now().Add(-time.Hour)),
		})
		if err != nil {
			t.Fatalf("failed to insert expired token: %v", err)
		}
		defer query.DeletePasswordResetTokenByID(ctx, created.ID)

		err = ResetPassword(ctx, ah.ResetPasswordParams{
			Token:       rawToken,
			NewPassword: "NewPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("success updates user password and deletes token", func(t *testing.T) {
		admin, err := CreateAdmin(ctx, user.CreateAdminParams{
			FirstName: "Reset",
			LastName:  "Admin",
			Email:     generateTestEmail(),
			Password:  testPassword,
		})
		if err != nil {
			t.Fatalf("failed to create admin: %v", err)
		}
		defer query.DeleteUser(ctx, admin.ID)

		rawToken := "valid-reset-token"
		tokenHash := testPasswordResetHash(rawToken)
		_, err = query.InsertPasswordResetToken(ctx, db.InsertPasswordResetTokenParams{
			UserID:    admin.ID,
			TokenHash: tokenHash,
			ExpiresAt: dbadapters.DBTime(time.Now().Add(time.Hour)),
		})
		if err != nil {
			t.Fatalf("failed to insert reset token: %v", err)
		}

		newPassword := "NewPass123!"
		err = ResetPassword(ctx, ah.ResetPasswordParams{
			Token:       rawToken,
			NewPassword: newPassword,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updated, err := query.GetUserById(ctx, admin.ID)
		if err != nil {
			t.Fatalf("failed to get updated user: %v", err)
		}
		if !password.ComparePassword(updated.PasswordHash, newPassword) {
			t.Fatal("expected password hash to match new password")
		}
		if password.ComparePassword(updated.PasswordHash, testPassword) {
			t.Fatal("expected password hash not to match old password")
		}

		_, err = query.GetPasswordResetTokenByHash(ctx, tokenHash)
		if !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("expected reset token to be deleted, got %v", err)
		}
	})
}

func testPasswordResetHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
