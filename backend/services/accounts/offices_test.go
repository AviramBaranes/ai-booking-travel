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

func validCreateOfficeParams() CreateOfficeRequest {
	phone := "0521234567"
	address := "123 Office St"
	return CreateOfficeRequest{
		Name:           "Test Office",
		OrganizationID: 1, // will be overridden in tests with a real org ID
		Phone:          &phone,
		Address:        &address,
	}
}

func validUpdateOfficeParams() UpdateOfficeRequest {
	name := "Updated Office"
	orgID := int32(1)
	phone := "0529876543"
	address := "456 Updated St"
	return UpdateOfficeRequest{
		Name:           &name,
		OrganizationID: &orgID,
		Phone:          &phone,
		Address:        &address,
	}
}

func officeMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func createTestOffice(t *testing.T, s *Service, orgID int32, name string) *OfficeResponse {
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

		page1, err := s.ListOffices(ctx, &ListOfficesRequest{
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

		page2, err := s.ListOffices(ctx, &ListOfficesRequest{
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
		page1IDs := make(map[int32]bool)
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

		resp, err := s.ListOffices(ctx, &ListOfficesRequest{
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

		resp, err := s.ListOffices(ctx, &ListOfficesRequest{
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

		resp, err := s.ListOffices(ctx, &ListOfficesRequest{
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
		resp, err = s.ListOffices(ctx, &ListOfficesRequest{
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

		resp, err := s.ListOffices(ctx, &ListOfficesRequest{
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
		p := ListOfficesRequest{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("returns error when list db fails", func(t *testing.T) {
		t.Parallel()
		q, s := officeMockService(t)
		q.EXPECT().ListOffices(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListOffices(ctx, &ListOfficesRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count db fails", func(t *testing.T) {
		t.Parallel()
		q, s := officeMockService(t)
		q.EXPECT().ListOffices(gomock.Any(), gomock.Any()).Return([]db.ListOfficesRow{}, nil)
		q.EXPECT().CountOffices(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.ListOffices(ctx, &ListOfficesRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
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

		resp, err := s.CreateOffice(ctx, CreateOfficeRequest{
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
		api_errors.AssertApiError(t, ErrNameAlreadyExists, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := officeMockService(t)
		q.EXPECT().CreateOffice(gomock.Any(), gomock.Any()).Return(db.Office{}, errors.New("db error"))

		_, err := s.CreateOffice(ctx, validCreateOfficeParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateOffice(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("updates all fields", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "UpdateOfficeFullOrg")
		created := createTestOffice(t, s, org.ID, "Update Full Office")

		params := validUpdateOfficeParams()
		params.OrganizationID = &org.ID // keep same org
		resp, err := s.UpdateOffice(ctx, created.ID, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != *params.Name {
			t.Fatalf("expected name %q, got %q", *params.Name, resp.Name)
		}
		if resp.Phone == nil || *resp.Phone != *params.Phone {
			t.Fatalf("expected phone %q, got %v", *params.Phone, resp.Phone)
		}
		if resp.Address == nil || *resp.Address != *params.Address {
			t.Fatalf("expected address %q, got %v", *params.Address, resp.Address)
		}
	})

	t.Run("partial update only changes provided fields", func(t *testing.T) {
		t.Parallel()
		org := createTestOrg(t, s, "PartialUpdateOfficeOrg")
		created := createTestOffice(t, s, org.ID, "Partial Update Office")

		newName := "Partial Updated Name"
		resp, err := s.UpdateOffice(ctx, created.ID, UpdateOfficeRequest{Name: &newName})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Partial Updated Name" {
			t.Fatalf("expected name 'Partial Updated Name', got %q", resp.Name)
		}
		if resp.OrganizationID != created.OrganizationID {
			t.Fatalf("expected organizationId unchanged %d, got %d", created.OrganizationID, resp.OrganizationID)
		}
		if (resp.Phone == nil) != (created.Phone == nil) || (resp.Phone != nil && *resp.Phone != *created.Phone) {
			t.Fatalf("expected phone unchanged %v, got %v", created.Phone, resp.Phone)
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
		_, err := s.UpdateOffice(ctx, officeB.ID, UpdateOfficeRequest{Name: &dupName})
		api_errors.AssertApiError(t, ErrNameAlreadyExists, err)
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		t.Parallel()
		p := validUpdateOfficeParams()
		blank := "   "
		p.Name = &blank
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := officeMockService(t)
		q.EXPECT().UpdateOffice(gomock.Any(), gomock.Any()).Return(db.Office{}, errors.New("db error"))

		_, err := s.UpdateOffice(ctx, 1, validUpdateOfficeParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
