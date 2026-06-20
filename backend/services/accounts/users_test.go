package accounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	oh "encore.app/services/accounts/handlers/office"
	organization "encore.app/services/accounts/handlers/organization"
	user "encore.app/services/accounts/handlers/user"
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

}

func TestGetUserCredit(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("it returns the credit of the office for a user of an inorganic organization", func(t *testing.T) {
		t.Parallel()
		icountID := int32(42)
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{Name: randomName(), IsOrganic: false})
		if err != nil {
			t.Fatalf("create org: %v", err)
		}
		off, err := s.CreateOffice(ctx, oh.CreateOfficeParams{Name: randomName(), OrganizationID: org.ID, IcountClientID: &icountID})
		if err != nil {
			t.Fatalf("create office: %v", err)
		}
		agent, err := s.CreateAgent(ctx, user.CreateAgentParams{
			FirstName: "Test", LastName: "Agent",
			Email: generateTestEmail(), Password: "Str0ng!Pass99",
			PhoneNumber: randomIsraeliPhoneNumber(), OfficeID: off.ID,
		})
		if err != nil {
			t.Fatalf("create agent: %v", err)
		}

		_ = s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: off.ID, BalanceChange: 150.0})

		credit, err := user.NewUserService(s.query).GetUserCredit(ctx, agent.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if credit.BalanceDue != 150.0 {
			t.Fatalf("expected balance_due 150.0 (office), got %v", credit.BalanceDue)
		}
	})

	t.Run("it returns the credit of the organization for a user of an organic organization", func(t *testing.T) {
		t.Parallel()
		icountID := int32(10)
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{Name: randomName(), IsOrganic: true, IcountClientID: &icountID})
		if err != nil {
			t.Fatalf("create org: %v", err)
		}
		off, err := s.CreateOffice(ctx, oh.CreateOfficeParams{Name: randomName(), OrganizationID: org.ID})
		if err != nil {
			t.Fatalf("create office: %v", err)
		}
		agent, err := s.CreateAgent(ctx, user.CreateAgentParams{
			FirstName: "Test", LastName: "Agent",
			Email: generateTestEmail(), Password: "Str0ng!Pass99",
			PhoneNumber: randomIsraeliPhoneNumber(), OfficeID: off.ID,
		})
		if err != nil {
			t.Fatalf("create agent: %v", err)
		}

		_ = s.UpdateOrganizationBalanceDue(ctx, organization.UpdateOrganizationBalanceDueParams{ID: org.ID, BalanceChange: 200.0})

		credit, err := user.NewUserService(s.query).GetUserCredit(ctx, agent.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if credit.BalanceDue != 200.0 {
			t.Fatalf("expected balance_due 200.0 (org), got %v", credit.BalanceDue)
		}
	})
}

func TestGetUserMarkupGross(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}
	strongPassword := "Str0ng!Pass99"

	t.Run("it returns the office markup gross if set", func(t *testing.T) {
		t.Parallel()
		officeMarkup := float64(12)
		icountID := int32(1)
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{Name: randomName(), IsOrganic: false})
		if err != nil {
			t.Fatalf("create org: %v", err)
		}
		off, err := s.CreateOffice(ctx, oh.CreateOfficeParams{Name: randomName(), OrganizationID: org.ID, IcountClientID: &icountID, GrossMarkup: &officeMarkup})
		if err != nil {
			t.Fatalf("create office: %v", err)
		}
		agent, err := s.CreateAgent(ctx, user.CreateAgentParams{
			FirstName: "Test", LastName: "Agent",
			Email: generateTestEmail(), Password: strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(), OfficeID: off.ID,
		})
		if err != nil {
			t.Fatalf("create agent: %v", err)
		}

		resp, err := user.NewUserService(s.query).GetUserMarkupGross(ctx, agent.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.GrossMarkup != officeMarkup {
			t.Fatalf("expected gross markup %.2f (office), got %.2f", officeMarkup, resp.GrossMarkup)
		}
	})

	t.Run("it returns the organization markup gross if set and office not set", func(t *testing.T) {
		t.Parallel()
		orgMarkup := float64(8)
		icountID := int32(2)
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{Name: randomName(), IsOrganic: false, GrossMarkup: orgMarkup})
		if err != nil {
			t.Fatalf("create org: %v", err)
		}
		// Office has no markup_gross (nil → 0 in DB)
		off, err := s.CreateOffice(ctx, oh.CreateOfficeParams{Name: randomName(), OrganizationID: org.ID, IcountClientID: &icountID})
		if err != nil {
			t.Fatalf("create office: %v", err)
		}
		agent, err := s.CreateAgent(ctx, user.CreateAgentParams{
			FirstName: "Test", LastName: "Agent",
			Email: generateTestEmail(), Password: strongPassword,
			PhoneNumber: randomIsraeliPhoneNumber(), OfficeID: off.ID,
		})
		if err != nil {
			t.Fatalf("create agent: %v", err)
		}

		resp, err := user.NewUserService(s.query).GetUserMarkupGross(ctx, agent.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.GrossMarkup != orgMarkup {
			t.Fatalf("expected gross markup %.2f (org), got %.2f", orgMarkup, resp.GrossMarkup)
		}
	})

	t.Run("it returns 0 for none agent user", func(t *testing.T) {
		t.Parallel()
		admin, cleanup, err := registerAdmin(ctx, generateTestEmail(), strongPassword)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		t.Cleanup(cleanup)

		resp, err := user.NewUserService(s.query).GetUserMarkupGross(ctx, admin.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.GrossMarkup != 0 {
			t.Fatalf("expected gross markup 0 for admin, got %.2f", resp.GrossMarkup)
		}
	})
}
