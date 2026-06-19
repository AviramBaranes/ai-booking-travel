package accounts

import (
	"context"
	"fmt"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	oh "encore.app/services/accounts/handlers/office"
	organization "encore.app/services/accounts/handlers/organization"
	user "encore.app/services/accounts/handlers/user"
	"encore.dev/et"
)

// --- Helpers ---

func validCreateOfficeParams() oh.CreateOfficeParams {
	phone := "0521234567"
	address := "123 Office St"
	return oh.CreateOfficeParams{
		Name:           "Test Office",
		OrganizationID: 1, // will be overridden in tests with a real org ID
		Phone:          &phone,
		Address:        &address,
	}
}

func validUpdateOfficeParams() oh.UpdateOfficeParams {
	name := "Updated Office"
	orgID := int64(1)
	phone := "0529876543"
	address := "456 Updated St"
	return oh.UpdateOfficeParams{
		Name:           name,
		OrganizationID: orgID,
		Phone:          &phone,
		Address:        &address,
	}
}

func createTestOffice(t *testing.T, s *Service, orgID int64, name string) *oh.OfficeResponse {
	t.Helper()
	p := validCreateOfficeParams()
	p.Name = name
	p.OrganizationID = orgID
	resp, err := s.CreateOffice(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to seed office %s: %v", name, err)
	}
	return resp
}

// --- Tests grouped by endpoint ---

func TestListOffices(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("pagination returns max 15 per page", func(t *testing.T) {
		t.Parallel()
		// Create an org to attach offices to
		org := createTestOrg(t, s, "PagOfficeOrg")

		for i := 1; i <= 18; i++ {
			createTestOffice(t, s, org.ID, fmt.Sprintf("PagOffice%02d Branch", i))
		}

		page1, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "PagOffice",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page1.Offices) != 15 {
			t.Fatalf("expected 15 offices on page 1, got %d", len(page1.Offices))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "PagOffice",
			Page:   2,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page2.Offices) != 3 {
			t.Fatalf("expected 3 offices on page 2, got %d", len(page2.Offices))
		}
		if page2.Total != 18 {
			t.Fatalf("expected total 18, got %d", page2.Total)
		}

		// No overlap between pages
		page1IDs := make(map[int64]bool)
		for _, o := range page1.Offices {
			page1IDs[o.ID] = true
		}
		for _, o := range page2.Offices {
			if page1IDs[o.ID] {
				t.Fatalf("office %d (%s) appeared on both pages", o.ID, o.Name)
			}
		}
	})

	t.Run("empty page returns no results", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "EmptyPageOfficeOrg")
		for i := 1; i <= 18; i++ {
			createTestOffice(t, s, org.ID, fmt.Sprintf("EmptyPageOffice%02d", i))
		}

		resp, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "EmptyPageOffice",
			Page:   3,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Offices) != 0 {
			t.Fatalf("expected 0 offices on page 3, got %d", len(resp.Offices))
		}
		if resp.Total != 18 {
			t.Fatalf("expected total 18 (unchanged), got %d", resp.Total)
		}
	})

	t.Run("filters by search", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "SearchOfficeOrg")
		createTestOffice(t, s, org.ID, "Searchable UniqueOFC123 Branch")

		resp, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "UniqueOFC123",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Offices) != 1 {
			t.Fatalf("expected 1 result, got %d", len(resp.Offices))
		}
		if resp.Offices[0].Name != "Searchable UniqueOFC123 Branch" {
			t.Fatalf("unexpected office: %s", resp.Offices[0].Name)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("filters by orgId", func(t *testing.T) {
		t.Parallel()
		orgA := createTestOrg(t, s, "OrgFilterA Offices")
		orgB := createTestOrg(t, s, "OrgFilterB Offices")

		createTestOffice(t, s, orgA.ID, "OrgFilterOffice A1")
		createTestOffice(t, s, orgA.ID, "OrgFilterOffice A2")
		createTestOffice(t, s, orgB.ID, "OrgFilterOffice B1")

		resp, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "OrgFilterOffice",
			OrgID:  orgA.ID,
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Offices) != 2 {
			t.Fatalf("expected 2 offices for orgA, got %d", len(resp.Offices))
		}
		if resp.Total != 2 {
			t.Fatalf("expected total 2, got %d", resp.Total)
		}

		// Org B should only see its own office
		resp, err = s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "OrgFilterOffice",
			OrgID:  orgB.ID,
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Offices) != 1 {
			t.Fatalf("expected 1 office for orgB, got %d", len(resp.Offices))
		}
		if resp.Offices[0].Name != "OrgFilterOffice B1" {
			t.Fatalf("expected office B1, got %s", resp.Offices[0].Name)
		}
	})

	t.Run("returns correct contactCount and agentCount", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "OfficeCountsOrg")
		office := createTestOffice(t, s, org.ID, "OfficeCountsTarget")

		// Create 2 contacts on this office
		_, err := query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "OffC", LastName: "One", Role: "manager",
			Cellphone: "0501110001", Email: "offc_count1@test.com",
			OfficeID: &office.ID,
		})
		if err != nil {
			t.Fatalf("failed to create contact 1: %v", err)
		}
		_, err = query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "OffC", LastName: "Two", Role: "sales",
			Cellphone: "0501110002", Email: "offc_count2@test.com",
			OfficeID: &office.ID,
		})
		if err != nil {
			t.Fatalf("failed to create contact 2: %v", err)
		}

		// Create 1 agent on this office
		_, err = query.CreateAgent(ctx, db.CreateAgentParams{
			Email: "agent_count1@offcounts.com", PasswordHash: "hash",
			OfficeID: &office.ID,
		})
		if err != nil {
			t.Fatalf("failed to create agent: %v", err)
		}

		resp, err := s.ListOffices(ctx, oh.ListOfficesParams{
			Search: "OfficeCountsTarget",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Offices) != 1 {
			t.Fatalf("expected 1 office, got %d", len(resp.Offices))
		}
		o := resp.Offices[0]
		if o.ContactCount != 2 {
			t.Fatalf("expected 2 contacts, got %d", o.ContactCount)
		}
		if o.AgentCount != 1 {
			t.Fatalf("expected 1 agent, got %d", o.AgentCount)
		}
	})

	t.Run("validation rejects page 0", func(t *testing.T) {
		t.Parallel()
		p := oh.ListOfficesParams{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

}

func TestCreateOffice(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates office with all fields", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "CreateOfficeFullOrg")

		p := validCreateOfficeParams()
		p.Name = "Create Full Office"
		p.OrganizationID = org.ID
		resp, err := s.CreateOffice(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.Name != "Create Full Office" {
			t.Fatalf("expected name 'Create Full Office', got %q", resp.Name)
		}
		if resp.OrganizationID != org.ID {
			t.Fatalf("expected organizationId %d, got %d", org.ID, resp.OrganizationID)
		}
		if resp.Phone == nil || *resp.Phone != *p.Phone {
			t.Fatalf("expected phone %q, got %v", *p.Phone, resp.Phone)
		}
		if resp.Address == nil || *resp.Address != *p.Address {
			t.Fatalf("expected address %q, got %v", *p.Address, resp.Address)
		}
	})

	t.Run("creates office with nil optional fields", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "CreateOfficeMinOrg")

		resp, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name:           "Minimal Office",
			OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Phone != nil {
			t.Fatalf("expected nil phone, got %v", resp.Phone)
		}
		if resp.Address != nil {
			t.Fatalf("expected nil address, got %v", resp.Address)
		}
	})

	t.Run("validation rejects icount_client_id under organic org", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "CreateOfficeOrganicIcountOrg")
		icountID := int32(55)
		_, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: "Should Fail Office", OrganizationID: org.ID, IcountClientID: &icountID,
		})
		api_errors.AssertApiError(t, oh.ErrOfficeOrganicForbidsIcountClientID, err)
	})

	t.Run("validation rejects missing icount_client_id under non-organic org", func(t *testing.T) {
		t.Parallel()
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "CreateOfficeNonOrganicNoIcount", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create org: %v", err)
		}
		_, err = s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: "Should Fail Office", OrganizationID: org.ID,
		})
		api_errors.AssertApiError(t, oh.ErrOfficeNonOrganicRequiresIcountClientID, err)
	})

	t.Run("creates office with icount_client_id under non-organic org", func(t *testing.T) {
		t.Parallel()
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "CreateOfficeNonOrganicWithIcount", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create org: %v", err)
		}
		icountID := int32(77)
		resp, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: "NonOrganic Office With Icount", OrganizationID: org.ID, IcountClientID: &icountID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.IcountClientID == nil || *resp.IcountClientID != 77 {
			t.Fatalf("expected icountClientId 77, got %v", resp.IcountClientID)
		}
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		t.Parallel()
		p := validCreateOfficeParams()
		p.Name = "   "
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects organizationId 0", func(t *testing.T) {
		t.Parallel()
		p := validCreateOfficeParams()
		p.OrganizationID = 0
		api_errors.AssertApiError(t, invalidValueErr("organizationId"), p.Validate())
	})

	t.Run("returns error on duplicate name", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "DupOfficeOrg")
		createTestOffice(t, s, org.ID, "Duplicate Office Name")

		p := validCreateOfficeParams()
		p.Name = "Duplicate Office Name"
		p.OrganizationID = org.ID
		_, err := s.CreateOffice(ctx, p)
		api_errors.AssertApiError(t, oh.ErrNameAlreadyExists, err)
	})

}

func TestUpdateOffice(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("validation rejects setting icount_client_id under organic org", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "UpdateOfficeOrganicIcountOrg")
		office := createTestOffice(t, s, org.ID, "UpdateOfficeOrganicIcountTarget")
		icountID := int32(55)
		p := validUpdateOfficeParams()
		p.Name = office.Name
		p.OrganizationID = org.ID
		p.IcountClientID = &icountID
		_, err := s.UpdateOffice(ctx, office.ID, p)
		api_errors.AssertApiError(t, oh.ErrOfficeOrganicForbidsIcountClientID, err)
	})

	t.Run("updates icount_client_id under non-organic org", func(t *testing.T) {
		t.Parallel()
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "UpdateOfficeNonOrganicIcountOrg", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create org: %v", err)
		}
		icountID := int32(42)
		office, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: "UpdateOfficeNonOrganicTarget", OrganizationID: org.ID, IcountClientID: &icountID,
		})
		if err != nil {
			t.Fatalf("failed to create office: %v", err)
		}
		newIcount := int32(99)
		p := validUpdateOfficeParams()
		p.Name = office.Name
		p.OrganizationID = org.ID
		p.IcountClientID = &newIcount
		resp, err := s.UpdateOffice(ctx, office.ID, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.IcountClientID == nil || *resp.IcountClientID != 99 {
			t.Fatalf("expected icountClientId 99, got %v", resp.IcountClientID)
		}
	})

	t.Run("validation rejects icount when moving office to organic org", func(t *testing.T) {
		t.Parallel()
		nonOrg, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "UpdateOfficeMoveToOrganicSrc", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create non-organic org: %v", err)
		}
		icountID := int32(42)
		office, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: "UpdateOfficeMoveToOrganicTarget", OrganizationID: nonOrg.ID, IcountClientID: &icountID,
		})
		if err != nil {
			t.Fatalf("failed to create office: %v", err)
		}
		orgOrg := createTestOrg(t, s, "UpdateOfficeMoveToOrganicDst")
		newIcount := int32(55)

		p := validUpdateOfficeParams()
		p.Name = office.Name
		p.OrganizationID = orgOrg.ID
		p.IcountClientID = &newIcount

		_, err = s.UpdateOffice(ctx, office.ID, p)
		api_errors.AssertApiError(t, oh.ErrOfficeOrganicForbidsIcountClientID, err)
	})

	t.Run("updates all fields", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "UpdateOfficeFullOrg")
		created := createTestOffice(t, s, org.ID, "Update Full Office")

		params := validUpdateOfficeParams()
		params.OrganizationID = org.ID // keep same org
		resp, err := s.UpdateOffice(ctx, created.ID, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != params.Name {
			t.Fatalf("expected name %q, got %q", params.Name, resp.Name)
		}
		if resp.Phone == nil || *resp.Phone != *params.Phone {
			t.Fatalf("expected phone %q, got %v", *params.Phone, resp.Phone)
		}
		if resp.Address == nil || *resp.Address != *params.Address {
			t.Fatalf("expected address %q, got %v", *params.Address, resp.Address)
		}
	})

	t.Run("returns not found when office does not exist", func(t *testing.T) {
		t.Parallel()
		_, err := s.UpdateOffice(ctx, 999999, validUpdateOfficeParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error on duplicate name", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "DupUpdateOfficeOrg")
		createTestOffice(t, s, org.ID, "Dup Update Office A")
		officeB := createTestOffice(t, s, org.ID, "Dup Update Office B")

		dupName := "Dup Update Office A"
		p := validUpdateOfficeParams()
		p.Name = dupName
		p.OrganizationID = org.ID
		_, err := s.UpdateOffice(ctx, officeB.ID, p)
		api_errors.AssertApiError(t, oh.ErrNameAlreadyExists, err)
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		t.Parallel()
		p := validUpdateOfficeParams()
		p.Name = "   "
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

}

func TestListInorganicOffices(t *testing.T) {
	ctx := context.Background()
	t.Run("it returns all inorganic offices", func(t *testing.T) {
		newDb, err := et.NewTestDatabase(ctx, "accounts")
		if err != nil {
			t.Fatalf("failed to create test database: %v", err)
		}
		s := newService(newDb)

		org1 := createTestOrg(t, s, randomName())
		org2, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{Name: randomName(), IsOrganic: false})
		if err != nil {
			t.Fatalf("create inorganic org: %v", err)
		}

		createTestOffice(t, s, org1.ID, randomName())
		createTestOffice(t, s, org1.ID, randomName())

		icountID1, icountID2 := int32(42), int32(43)
		inorganicOffice1, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: randomName(), OrganizationID: org2.ID, IcountClientID: &icountID1,
		})
		if err != nil {
			t.Fatalf("create inorganic office 1: %v", err)
		}
		inorganicOffice2, err := s.CreateOffice(ctx, oh.CreateOfficeParams{
			Name: randomName(), OrganizationID: org2.ID, IcountClientID: &icountID2,
		})
		if err != nil {
			t.Fatalf("create inorganic office 2: %v", err)
		}

		resp, err := s.ListInorganicOffices(ctx)

		if err != nil {
			t.Fatalf("failed to list inorganic offices: %v", err)
		}

		if len(resp.Offices) != 2 {
			t.Fatalf("expected 2 inorganic offices, got %d", len(resp.Offices))
		}

		expectedResults := []oh.InorganicOffice{
			{ID: inorganicOffice1.ID, Name: inorganicOffice1.Name},
			{ID: inorganicOffice2.ID, Name: inorganicOffice2.Name},
		}

		for i, office := range resp.Offices {
			if office.ID != expectedResults[i].ID || office.Name != expectedResults[i].Name {
				t.Fatalf("expected office %v, got %v", expectedResults[i], office)
			}
		}
	})
}

func TestUpdateOfficeBalanceDue(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	setup := func(t *testing.T) (officeID int64, agentID int64) {
		t.Helper()
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
		return off.ID, agent.ID
	}

	getBalance := func(t *testing.T, agentID int64) float64 {
		t.Helper()
		credit, err := user.NewUserService(s.query).GetUserCredit(ctx, agentID)
		if err != nil {
			t.Fatalf("GetUserCredit: %v", err)
		}
		return credit.BalanceDue
	}

	t.Run("positive delta increases balance", func(t *testing.T) {
		t.Parallel()
		officeID, agentID := setup(t)
		if err := s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: officeID, BalanceChange: 100.0}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got := getBalance(t, agentID); got != 100.0 {
			t.Fatalf("expected balance 100.0, got %v", got)
		}
	})

	t.Run("negative delta decreases balance (payment resolved)", func(t *testing.T) {
		t.Parallel()
		officeID, agentID := setup(t)
		_ = s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: officeID, BalanceChange: 100.0})
		if err := s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: officeID, BalanceChange: -60.0}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got := getBalance(t, agentID); got != 40.0 {
			t.Fatalf("expected balance 40.0, got %v", got)
		}
	})

	t.Run("balance never goes below 0", func(t *testing.T) {
		t.Parallel()
		officeID, agentID := setup(t)
		_ = s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: officeID, BalanceChange: 50.0})
		if err := s.UpdateOfficeBalanceDue(ctx, oh.UpdateOfficeBalanceDueParams{ID: officeID, BalanceChange: -200.0}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got := getBalance(t, agentID); got != 0.0 {
			t.Fatalf("expected balance 0.0, got %v", got)
		}
	})
}
