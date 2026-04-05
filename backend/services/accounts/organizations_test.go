package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateOrgParams() CreateOrganizationRequest {
	phone := "0521234567"
	address := "123 Test St"
	obligo := 1000.0
	return CreateOrganizationRequest{
		Name:      "Test Organization",
		IsOrganic: true,
		Phone:     &phone,
		Address:   &address,
		Obligo:    &obligo,
	}
}

func validUpdateOrgParams() UpdateOrganizationRequest {
	name := "Updated Organization"
	isOrganic := false
	phone := "0529876543"
	address := "456 Updated St"
	obligo := 2000.0
	return UpdateOrganizationRequest{
		Name:      &name,
		IsOrganic: &isOrganic,
		Phone:     &phone,
		Address:   &address,
		Obligo:    &obligo,
	}
}

func orgMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func createTestOrg(t *testing.T, s *Service, name string) *OrganizationResponse {
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

		page1, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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

		page2, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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
		page1IDs := make(map[int32]bool)
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
		resp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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

		resp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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

	t.Run("search returns no results for non-matching query", func(t *testing.T) {
		resp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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
		_, err = s.CreateOrganization(ctx, nonOrgP)
		if err != nil {
			t.Fatalf("failed to create non-organic org: %v", err)
		}

		resp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
			Search:    "OrganicFilterTest",
			IsOrganic: true,
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
		resp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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
		unResp, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{
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
		p := ListOrganizationsRequest{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("validation rejects negative page", func(t *testing.T) {
		p := ListOrganizationsRequest{Page: -1}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("returns error when list db fails", func(t *testing.T) {
		q, s := orgMockService(t)
		q.EXPECT().ListOrganizations(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count db fails", func(t *testing.T) {
		q, s := orgMockService(t)
		q.EXPECT().ListOrganizations(gomock.Any(), gomock.Any()).Return([]db.ListOrganizationsRow{}, nil)
		q.EXPECT().CountOrganizations(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.ListOrganizations(ctx, &ListOrganizationsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateOrganization(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

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
		neg := -1.0
		p.Obligo = &neg
		api_errors.AssertApiError(t, invalidValueErr("obligo"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateOrgParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation accepts nil optional fields", func(t *testing.T) {
		p := CreateOrganizationRequest{Name: "Minimal Org", IsOrganic: false}
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
	})

	t.Run("creates organization with nil optional fields", func(t *testing.T) {
		resp, err := s.CreateOrganization(ctx, CreateOrganizationRequest{
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
		api_errors.AssertApiError(t, ErrNameAlreadyExists, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := orgMockService(t)
		q.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(db.Organization{}, errors.New("db error"))

		_, err := s.CreateOrganization(ctx, validCreateOrgParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateOrganization(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("validation rejects blank name", func(t *testing.T) {
		p := validUpdateOrgParams()
		blank := "   "
		p.Name = &blank
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects negative obligo", func(t *testing.T) {
		p := validUpdateOrgParams()
		neg := -1.0
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
		p := UpdateOrganizationRequest{Name: &name}
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
	})

	t.Run("partial update only changes provided fields", func(t *testing.T) {
		created := createTestOrg(t, s, "Partial Update Org")

		newName := "Partial Updated Name"
		resp, err := s.UpdateOrganization(ctx, created.ID, UpdateOrganizationRequest{Name: &newName})
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
		_, err := s.UpdateOrganization(ctx, orgB.ID, UpdateOrganizationRequest{Name: &dupName})
		api_errors.AssertApiError(t, ErrNameAlreadyExists, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := orgMockService(t)
		q.EXPECT().UpdateOrganization(gomock.Any(), gomock.Any()).Return(db.Organization{}, errors.New("db error"))

		_, err := s.UpdateOrganization(ctx, 1, validUpdateOrgParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
