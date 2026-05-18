package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	organization "encore.app/services/accounts/handlers/organization"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateOrgParams() organization.CreateOrganizationParams {
	phone := "0521234567"
	address := "123 Test St"
	obligo := int32(1000)
	icountClientID := int32(100)
	return organization.CreateOrganizationParams{
		Name:           "Test Organization",
		IsOrganic:      true,
		IcountClientID: &icountClientID,
		Phone:          &phone,
		Address:        &address,
		Obligo:         &obligo,
	}
}

func validUpdateOrgParams() organization.UpdateOrganizationParams {
	name := "Updated Organization"
	isOrganic := true
	icountClientID := int32(200)
	phone := "0529876543"
	address := "456 Updated St"
	obligo := int32(2000)
	return organization.UpdateOrganizationParams{
		Name:           &name,
		IsOrganic:      &isOrganic,
		IcountClientID: &icountClientID,
		Phone:          &phone,
		Address:        &address,
		Obligo:         &obligo,
	}
}

func createTestOrg(t *testing.T, s *Service, name string) *organization.OrganizationResponse {
	t.Helper()
	p := validCreateOrgParams()
	p.Name = name
	resp, err := s.CreateOrganization(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to seed org %s: %v", name, err)
	}
	return resp
}

// --- Tests grouped by endpoint ---

func TestListOrganizations(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("pagination returns max 15 per page", func(t *testing.T) {
		// Create 18 orgs with a unique prefix for filtering
		for i := 1; i <= 18; i++ {
			createTestOrg(t, s, fmt.Sprintf("PagOrg%02d Corp", i))
		}

		page1, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "PagOrg",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page1.Organizations) != 15 {
			t.Fatalf("expected 15 orgs on page 1, got %d", len(page1.Organizations))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "PagOrg",
			Page:   2,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page2.Organizations) != 3 {
			t.Fatalf("expected 3 orgs on page 2, got %d", len(page2.Organizations))
		}
		if page2.Total != 18 {
			t.Fatalf("expected total 18, got %d", page2.Total)
		}

		// No overlap between pages
		page1IDs := make(map[int64]bool)
		for _, o := range page1.Organizations {
			page1IDs[o.ID] = true
		}
		for _, o := range page2.Organizations {
			if page1IDs[o.ID] {
				t.Fatalf("org %d (%s) appeared on both pages", o.ID, o.Name)
			}
		}
	})

	t.Run("empty page returns no results", func(t *testing.T) {
		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "PagOrg",
			Page:   3,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 0 {
			t.Fatalf("expected 0 orgs on page 3, got %d", len(resp.Organizations))
		}
		if resp.Total != 18 {
			t.Fatalf("expected total 18 (unchanged), got %d", resp.Total)
		}
	})

	t.Run("filters by search", func(t *testing.T) {
		createTestOrg(t, s, "Searchable UniqueXYZ Corp")

		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "UniqueXYZ",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 1 {
			t.Fatalf("expected 1 result, got %d", len(resp.Organizations))
		}
		if resp.Organizations[0].Name != "Searchable UniqueXYZ Corp" {
			t.Fatalf("unexpected org: %s", resp.Organizations[0].Name)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("lists icount_client_id when set", func(t *testing.T) {
		icountID := int32(777)
		_, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "IcountListOrg UniqueXXX", IsOrganic: true, IcountClientID: &icountID,
		})
		if err != nil {
			t.Fatalf("failed to create org: %v", err)
		}
		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{Search: "IcountListOrg UniqueXXX", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 1 {
			t.Fatalf("expected 1 org, got %d", len(resp.Organizations))
		}
		o := resp.Organizations[0]
		if o.IcountClientID == nil || *o.IcountClientID != 777 {
			t.Fatalf("expected icountClientId 777, got %v", o.IcountClientID)
		}
	})

	t.Run("search returns no results for non-matching query", func(t *testing.T) {
		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "ZZZNoMatchHere999",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 0 {
			t.Fatalf("expected 0 results, got %d", len(resp.Organizations))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})

	t.Run("filters by isOrganic", func(t *testing.T) {
		// Create one organic and one non-organic with unique prefix
		orgP := validCreateOrgParams()
		orgP.Name = "OrganicFilterTest Org1"
		orgP.IsOrganic = true
		_, err := s.CreateOrganization(ctx, orgP)
		if err != nil {
			t.Fatalf("failed to create organic org: %v", err)
		}

		nonOrgP := validCreateOrgParams()
		nonOrgP.Name = "OrganicFilterTest Org2"
		nonOrgP.IsOrganic = false
		nonOrgP.IcountClientID = nil
		_, err = s.CreateOrganization(ctx, nonOrgP)
		if err != nil {
			t.Fatalf("failed to create non-organic org: %v", err)
		}

		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search:    "OrganicFilterTest",
			IsOrganic: "true",
			Page:      1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 1 {
			t.Fatalf("expected 1 organic org, got %d", len(resp.Organizations))
		}
		if resp.Organizations[0].Name != "OrganicFilterTest Org1" {
			t.Fatalf("expected organic org, got %s", resp.Organizations[0].Name)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("returns correct office contact and agent counts", func(t *testing.T) {
		// Create org
		org := createTestOrg(t, s, "Counts Test Org")

		// Create 2 offices under this org
		office1, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name: "Counts Office A", OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("failed to create office1: %v", err)
		}
		office2, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name: "Counts Office B", OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("failed to create office2: %v", err)
		}

		// Create contacts: 1 on org directly, 1 on office1, 1 on office2
		_, err = query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "OrgContact", LastName: "One", Role: "manager",
			Cellphone: "0501111111", Email: "orgc1@test.com",
			OrganizationID: &org.ID,
		})
		if err != nil {
			t.Fatalf("failed to create org contact: %v", err)
		}
		_, err = query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "OfficeContact", LastName: "Two", Role: "sales",
			Cellphone: "0502222222", Email: "offc2@test.com",
			OfficeID: &office1.ID,
		})
		if err != nil {
			t.Fatalf("failed to create office1 contact: %v", err)
		}
		_, err = query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "OfficeContact", LastName: "Three", Role: "sales",
			Cellphone: "0503333333", Email: "offc3@test.com",
			OfficeID: &office2.ID,
		})
		if err != nil {
			t.Fatalf("failed to create office2 contact: %v", err)
		}

		// Create 2 agents: one in office1, one in office2
		_, err = query.CreateAgent(ctx, db.CreateAgentParams{
			Email: "agent1@counts.com", PasswordHash: "hash",
			OfficeID: &office1.ID,
		})
		if err != nil {
			t.Fatalf("failed to create agent1: %v", err)
		}
		_, err = query.CreateAgent(ctx, db.CreateAgentParams{
			Email: "agent2@counts.com", PasswordHash: "hash",
			OfficeID: &office2.ID,
		})
		if err != nil {
			t.Fatalf("failed to create agent2: %v", err)
		}

		// Also create an unrelated org with its own office/contact/agent
		unrelatedOrg := createTestOrg(t, s, "Counts Unrelated Org")
		unrelatedOffice, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name: "Unrelated Office", OrganizationID: unrelatedOrg.ID,
		})
		if err != nil {
			t.Fatalf("failed to create unrelated office: %v", err)
		}
		_, err = query.CreateContact(ctx, db.CreateContactParams{
			FirstName: "Unrelated", LastName: "Contact", Role: "other",
			Cellphone: "0504444444", Email: "unrelated@test.com",
			OfficeID: &unrelatedOffice.ID,
		})
		if err != nil {
			t.Fatalf("failed to create unrelated contact: %v", err)
		}
		_, err = query.CreateAgent(ctx, db.CreateAgentParams{
			Email: "unrelated-agent@counts.com", PasswordHash: "hash",
			OfficeID: &unrelatedOffice.ID,
		})
		if err != nil {
			t.Fatalf("failed to create unrelated agent: %v", err)
		}

		// Query and find our target org
		resp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "Counts Test Org",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Organizations) != 1 {
			t.Fatalf("expected 1 org, got %d", len(resp.Organizations))
		}

		o := resp.Organizations[0]
		if o.OfficeCount != 2 {
			t.Fatalf("expected 2 offices, got %d", o.OfficeCount)
		}
		if o.ContactCount != 3 {
			t.Fatalf("expected 3 contacts, got %d", o.ContactCount)
		}
		if o.AgentCount != 2 {
			t.Fatalf("expected 2 agents, got %d", o.AgentCount)
		}

		// Verify unrelated org has its own counts
		unResp, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{
			Search: "Counts Unrelated Org",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(unResp.Organizations) != 1 {
			t.Fatalf("expected 1 org, got %d", len(unResp.Organizations))
		}
		u := unResp.Organizations[0]
		if u.OfficeCount != 1 {
			t.Fatalf("expected 1 office for unrelated org, got %d", u.OfficeCount)
		}
		if u.ContactCount != 1 {
			t.Fatalf("expected 1 contact for unrelated org, got %d", u.ContactCount)
		}
		if u.AgentCount != 1 {
			t.Fatalf("expected 1 agent for unrelated org, got %d", u.AgentCount)
		}
	})

	t.Run("validation rejects page 0", func(t *testing.T) {
		p := organization.ListOrganizationsParams{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("validation rejects negative page", func(t *testing.T) {
		p := organization.ListOrganizationsParams{Page: -1}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("returns error when list db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListOrganizations(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListOrganizations(gomock.Any(), gomock.Any()).Return([]db.ListOrganizationsRow{}, nil)
		q.EXPECT().CountOrganizations(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.ListOrganizations(ctx, organization.ListOrganizationsParams{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateOrganization(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("validation rejects organic without icount_client_id", func(t *testing.T) {
		p := validCreateOrgParams()
		p.IcountClientID = nil
		api_errors.AssertApiError(t, organization.ErrOrganizationOrganicRequiresIcountClientID, p.Validate())
	})

	t.Run("validation rejects non-organic with icount_client_id", func(t *testing.T) {
		icountID := int32(99)
		p := organization.CreateOrganizationParams{Name: "Foo", IsOrganic: false, IcountClientID: &icountID}
		api_errors.AssertApiError(t, organization.ErrOrganizationNonOrganicForbidsIcountClientID, p.Validate())
	})

	t.Run("validation rejects missing name", func(t *testing.T) {
		p := validCreateOrgParams()
		p.Name = ""
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		p := validCreateOrgParams()
		p.Name = "   "
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects negative obligo", func(t *testing.T) {
		p := validCreateOrgParams()
		neg := int32(-1)
		p.Obligo = &neg
		api_errors.AssertApiError(t, invalidValueErr("obligo"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateOrgParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation accepts nil optional fields", func(t *testing.T) {
		p := organization.CreateOrganizationParams{Name: "Minimal Org", IsOrganic: false}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("creates organization successfully", func(t *testing.T) {
		p := validCreateOrgParams()
		p.Name = "Create Test Org"
		resp, err := s.CreateOrganization(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.Name != "Create Test Org" {
			t.Fatalf("expected name 'Create Test Org', got %q", resp.Name)
		}
		if !resp.IsOrganic {
			t.Fatal("expected isOrganic to be true")
		}
		if resp.Phone == nil || *resp.Phone != *p.Phone {
			t.Fatalf("expected phone %q, got %v", *p.Phone, resp.Phone)
		}
		if resp.Address == nil || *resp.Address != *p.Address {
			t.Fatalf("expected address %q, got %v", *p.Address, resp.Address)
		}
		if resp.Obligo == nil || *resp.Obligo != *p.Obligo {
			t.Fatalf("expected obligo %v, got %v", *p.Obligo, resp.Obligo)
		}
		if resp.IcountClientID == nil || *resp.IcountClientID != *p.IcountClientID {
			t.Fatalf("expected icountClientId %v, got %v", p.IcountClientID, resp.IcountClientID)
		}
	})

	t.Run("creates organization with nil optional fields", func(t *testing.T) {
		resp, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name:      "Minimal Create Org",
			IsOrganic: false,
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
		if resp.Obligo != nil {
			t.Fatalf("expected nil obligo, got %v", resp.Obligo)
		}
	})

	t.Run("returns error on duplicate name", func(t *testing.T) {
		createTestOrg(t, s, "Duplicate Name Org")

		p := validCreateOrgParams()
		p.Name = "Duplicate Name Org"
		_, err := s.CreateOrganization(ctx, p)
		api_errors.AssertApiError(t, organization.ErrNameAlreadyExists, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(db.Organization{}, errors.New("db error"))

		_, err := s.CreateOrganization(ctx, validCreateOrgParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateOrganization(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("validation rejects non-organic with icount_client_id", func(t *testing.T) {
		isOrganic := false
		icountID := int32(50)
		p := organization.UpdateOrganizationParams{IsOrganic: &isOrganic, IcountClientID: &icountID}
		api_errors.AssertApiError(t, organization.ErrOrganizationNonOrganicForbidsIcountClientID, p.Validate())
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		p := validUpdateOrgParams()
		blank := "   "
		p.Name = &blank
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects negative obligo", func(t *testing.T) {
		p := validUpdateOrgParams()
		neg := int32(-1)
		p.Obligo = &neg
		api_errors.AssertApiError(t, invalidValueErr("obligo"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validUpdateOrgParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation accepts partial update with only name", func(t *testing.T) {
		name := "Partial"
		p := organization.UpdateOrganizationParams{Name: &name}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates organization successfully", func(t *testing.T) {
		created := createTestOrg(t, s, "Update Full Org")

		params := validUpdateOrgParams()
		resp, err := s.UpdateOrganization(ctx, created.ID, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != *params.Name {
			t.Fatalf("expected name %q, got %q", *params.Name, resp.Name)
		}
		if resp.IsOrganic != *params.IsOrganic {
			t.Fatalf("expected isOrganic %v, got %v", *params.IsOrganic, resp.IsOrganic)
		}
		if resp.Phone == nil || *resp.Phone != *params.Phone {
			t.Fatalf("expected phone %q, got %v", *params.Phone, resp.Phone)
		}
		if resp.Address == nil || *resp.Address != *params.Address {
			t.Fatalf("expected address %q, got %v", *params.Address, resp.Address)
		}
		if resp.Obligo == nil || *resp.Obligo != *params.Obligo {
			t.Fatalf("expected obligo %v, got %v", *params.Obligo, resp.Obligo)
		}
		if resp.IcountClientID == nil || *resp.IcountClientID != *params.IcountClientID {
			t.Fatalf("expected icountClientId %v, got %v", params.IcountClientID, resp.IcountClientID)
		}
	})

	t.Run("returns error when adding icount to non-organic org", func(t *testing.T) {
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "NonOrganic IcountAdd UniqueNOIA", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create non-organic org: %v", err)
		}
		icountID := int32(50)
		_, err = s.UpdateOrganization(ctx, org.ID, organization.UpdateOrganizationParams{IcountClientID: &icountID})
		api_errors.AssertApiError(t, organization.ErrOrganizationNonOrganicForbidsIcountClientID, err)
	})

	t.Run("returns error when switching to organic without icount_client_id", func(t *testing.T) {
		org, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name: "NonOrganic ToOrganic UniqueNOTO", IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create non-organic org: %v", err)
		}
		isOrganic := true
		_, err = s.UpdateOrganization(ctx, org.ID, organization.UpdateOrganizationParams{IsOrganic: &isOrganic})
		api_errors.AssertApiError(t, organization.ErrOrganizationOrganicRequiresIcountClientID, err)
	})

	t.Run("partial update only changes provided fields", func(t *testing.T) {
		created := createTestOrg(t, s, "Partial Update Org")

		newName := "Partial Updated Name"
		resp, err := s.UpdateOrganization(ctx, created.ID, organization.UpdateOrganizationParams{Name: &newName})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Partial Updated Name" {
			t.Fatalf("expected name 'Partial Updated Name', got %q", resp.Name)
		}
		if resp.IsOrganic != created.IsOrganic {
			t.Fatalf("expected isOrganic unchanged %v, got %v", created.IsOrganic, resp.IsOrganic)
		}
		if (resp.Phone == nil) != (created.Phone == nil) || (resp.Phone != nil && *resp.Phone != *created.Phone) {
			t.Fatalf("expected phone unchanged %v, got %v", created.Phone, resp.Phone)
		}
	})

	t.Run("returns not found when org does not exist", func(t *testing.T) {
		_, err := s.UpdateOrganization(ctx, 999999, validUpdateOrgParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error on duplicate name", func(t *testing.T) {
		createTestOrg(t, s, "Dup Update Org A")
		orgB := createTestOrg(t, s, "Dup Update Org B")

		dupName := "Dup Update Org A"
		_, err := s.UpdateOrganization(ctx, orgB.ID, organization.UpdateOrganizationParams{Name: &dupName})
		api_errors.AssertApiError(t, organization.ErrNameAlreadyExists, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateOrganization(gomock.Any(), gomock.Any()).Return(db.Organization{}, errors.New("db error"))

		_, err := s.UpdateOrganization(ctx, 1, validUpdateOrgParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestListOrganicOrganizations(t *testing.T) {
	ctx := context.Background()
	t.Run("it returns all organic organizations", func(t *testing.T) {
		newDb, err := et.NewTestDatabase(ctx, "accounts")
		if err != nil {
			t.Fatalf("failed to create test database: %v", err)
		}
		s := newService(newDb)

		icountID1, icountID2 := int32(10), int32(11)
		org1, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name:           "organic_org_1",
			IsOrganic:      true,
			IcountClientID: &icountID1,
		})

		if err != nil {
			t.Fatalf("failed to create org1: %v", err)
		}

		org2, err := s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name:           "organic_org_2",
			IsOrganic:      true,
			IcountClientID: &icountID2,
		})

		if err != nil {
			t.Fatalf("failed to create org2: %v", err)
		}

		_, err = s.CreateOrganization(ctx, organization.CreateOrganizationParams{
			Name:      "inorganic_org",
			IsOrganic: false,
		})

		if err != nil {
			t.Fatalf("failed to create inorganic org: %v", err)
		}

		resp, err := s.ListOrganicOrganizations(ctx)

		if err != nil {
			t.Fatalf("failed to list organic organizations: %v", err)
		}

		if len(resp.Organizations) != 2 {
			t.Fatalf("expected 2 organic organizations, got %d", len(resp.Organizations))
		}

		expectedResults := []organization.OrganicOrganization{
			{ID: org1.ID, Name: org1.Name},
			{ID: org2.ID, Name: org2.Name},
		}

		for i, org := range resp.Organizations {
			if org.ID != expectedResults[i].ID || org.Name != expectedResults[i].Name {
				t.Fatalf("expected organization %v, got %v", expectedResults[i], org)
			}
		}
	})
}
