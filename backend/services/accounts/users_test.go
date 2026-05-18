package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	user "encore.app/services/accounts/handlers/user"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func ptrInt32(v int32) *int32 { return &v }

// --- Tests ---

func TestUpdateUser(t *testing.T) {
	s := &Service{query: query}
	ctx := context.Background()
	strongPassword := "Str0ng!Pass99"

	// ── Integration: update email ──

	t.Run("updates email successfully", func(t *testing.T) {
		t.Parallel()
		agent, cleanup, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		newEmail := fmt.Sprintf("updated_%d@test.com", time.Now().UnixNano())
		resp, err := s.UpdateUser(ctx, agent.ID, user.UpdateUserParams{
			Email: ptrStr(newEmail),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Email != newEmail {
			t.Fatalf("expected email %s, got %s", newEmail, resp.Email)
		}
	})

	// ── Integration: update phone ──

	t.Run("updates phone successfully", func(t *testing.T) {
		t.Parallel()
		agent, cleanup, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		newPhone := randomIsraeliPhoneNumber()
		resp, err := s.UpdateUser(ctx, agent.ID, user.UpdateUserParams{
			PhoneNumber: ptrStr(newPhone),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.PhoneNumber == nil || *resp.PhoneNumber != newPhone {
			t.Fatalf("expected phone %s, got %v", newPhone, resp.PhoneNumber)
		}
	})

	// ── Integration: update officeId for agent ──

	t.Run("updates officeId for agent", func(t *testing.T) {
		t.Parallel()
		_, officeA := seedOrgAndOffice(t)
		_, officeB := seedOrgAndOffice(t)

		agent, cleanup, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
			OfficeID:    officeA,
		})
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		resp, err := s.UpdateUser(ctx, agent.ID, user.UpdateUserParams{
			OfficeID: &officeB,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.OfficeID == nil || *resp.OfficeID != officeB {
			t.Fatalf("expected officeId %d, got %v", officeB, resp.OfficeID)
		}
	})

	// ── Integration: setting officeId on admin fails at DB layer ──

	t.Run("setting officeId on admin fails with invalid officeId", func(t *testing.T) {
		t.Parallel()
		admin, cleanup, err := registerAdmin(ctx, generateTestEmail(), strongPassword)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		officeID := int64(999999)
		// Use a non-existent office ID – FK constraint should reject this.
		_, err = s.UpdateUser(ctx, admin.ID, user.UpdateUserParams{
			OfficeID: &officeID,
		})
		if err == nil {
			t.Fatal("expected error when setting invalid officeId on admin, got nil")
		}
	})

	// ── Integration: setting officeId on admin fails (check constraint) ──

	t.Run("setting officeId on admin fails even with valid office", func(t *testing.T) {
		t.Parallel()
		admin, cleanup, err := registerAdmin(ctx, generateTestEmail(), strongPassword)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		_, validOffice := seedOrgAndOffice(t)

		_, err = s.UpdateUser(ctx, admin.ID, user.UpdateUserParams{
			OfficeID: &validOffice,
		})
		if err == nil {
			t.Fatal("expected error when setting officeId on admin, got nil")
		}
	})

	// ── Integration: duplicate email ──

	t.Run("returns error on duplicate email", func(t *testing.T) {
		t.Parallel()
		emailA := generateTestEmail()
		agentA, cleanupA, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "AgentA",
			Email:       emailA,
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup agentA: %v", err)
		}
		_ = agentA
		t.Cleanup(cleanupA)

		agentB, cleanupB, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "AgentB",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup agentB: %v", err)
		}
		t.Cleanup(cleanupB)

		// Try to set agentB's email to agentA's email
		_, err = s.UpdateUser(ctx, agentB.ID, user.UpdateUserParams{
			Email: ptrStr(emailA),
		})
		api_errors.AssertApiError(t, user.ErrEmailAlreadyExists, err)
	})

	// ── Integration: same email for same user is fine ──

	t.Run("allows setting same email on same user", func(t *testing.T) {
		t.Parallel()
		email := generateTestEmail()
		agent, cleanup, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       email,
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		resp, err := s.UpdateUser(ctx, agent.ID, user.UpdateUserParams{
			Email: ptrStr(email),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Email != email {
			t.Fatalf("expected email %s, got %s", email, resp.Email)
		}
	})

	// ── Integration: duplicate phone ──

	t.Run("returns error on duplicate phone", func(t *testing.T) {
		t.Parallel()
		phone := randomIsraeliPhoneNumber()
		_, cleanupA, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "AgentA",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: phone,
		})
		if err != nil {
			t.Fatalf("setup agentA: %v", err)
		}
		t.Cleanup(cleanupA)

		agentB, cleanupB, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "AgentB",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("setup agentB: %v", err)
		}
		t.Cleanup(cleanupB)

		_, err = s.UpdateUser(ctx, agentB.ID, user.UpdateUserParams{
			PhoneNumber: ptrStr(phone),
		})
		api_errors.AssertApiError(t, user.ErrPhoneAlreadyExists, err)
	})

	// ── Integration: same phone for same user is fine ──

	t.Run("allows setting same phone on same user", func(t *testing.T) {
		t.Parallel()
		phone := randomIsraeliPhoneNumber()
		agent, cleanup, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       generateTestEmail(),
			Password:    strongPassword,
			PhoneNumber: phone,
		})
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		resp, err := s.UpdateUser(ctx, agent.ID, user.UpdateUserParams{
			PhoneNumber: ptrStr(phone),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.PhoneNumber == nil || *resp.PhoneNumber != phone {
			t.Fatalf("expected phone %s, got %v", phone, resp.PhoneNumber)
		}
	})

	// ── Integration: user not found ──

	t.Run("returns not found for nonexistent user", func(t *testing.T) {
		t.Parallel()
		_, err := s.UpdateUser(ctx, 999999, user.UpdateUserParams{
			Email: ptrStr("nobody@test.com"),
		})
		api_errors.AssertApiError(t, user.ErrUserNotFound, err)
	})

	// ── Validation ──

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		err := user.UpdateUserParams{Email: ptrStr("not-an-email")}.Validate()
		api_errors.AssertApiError(t, invalidValueErr("email"), err)
	})

	t.Run("validation rejects officeId 0", func(t *testing.T) {
		t.Parallel()
		officeID := int64(0)
		err := user.UpdateUserParams{OfficeID: &officeID}.Validate()
		api_errors.AssertApiError(t, invalidValueErr("officeId"), err)
	})

	t.Run("validation rejects weak password", func(t *testing.T) {
		t.Parallel()
		err := user.UpdateUserParams{Password: ptrStr("weak")}.Validate()
		if err == nil {
			t.Fatal("expected validation error for weak password, got nil")
		}
	})

	// ── Mock: check email uniqueness DB failure ──

	t.Run("returns error when check email db fails", func(t *testing.T) {
		t.Parallel()
		q, ms := mockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("db error"))

		_, err := ms.UpdateUser(ctx, 1, user.UpdateUserParams{
			Email: ptrStr("fail@test.com"),
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	// ── Mock: check phone uniqueness DB failure ──

	t.Run("returns error when check phone db fails", func(t *testing.T) {
		t.Parallel()
		q, ms := mockService(t)
		q.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
			Return(db.User{}, errors.New("db error"))

		_, err := ms.UpdateUser(ctx, 1, user.UpdateUserParams{
			PhoneNumber: ptrStr("0501234567"),
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	// ── Mock: UpdateUser DB failure ──

	t.Run("returns error when update db fails", func(t *testing.T) {
		t.Parallel()
		q, ms := mockService(t)
		q.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
			Return(db.UpdateUserRow{}, errors.New("db error"))

		_, err := ms.UpdateUser(ctx, 1, user.UpdateUserParams{})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
