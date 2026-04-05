package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateContactParams(officeID *int32, orgID *int32) CreateContactRequest {
	return CreateContactRequest{
		FirstName:      "John",
		LastName:       "Doe",
		Role:           "manager",
		Cellphone:      "0521234567",
		Email:          "john.doe@test.com",
		OfficeID:       officeID,
		OrganizationID: orgID,
	}
}

func validUpdateContactParams() UpdateContactRequest {
	firstName := "Jane"
	lastName := "Smith"
	role := "director"
	cellphone := "0529876543"
	email := "jane.smith@test.com"
	return UpdateContactRequest{
		FirstName: &firstName,
		LastName:  &lastName,
		Role:      &role,
		Cellphone: &cellphone,
		Email:     &email,
	}
}

func contactMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

// seedOrgAndOffice creates an org and office for use in contact tests.
func seedOrgAndOffice(t *testing.T) (orgID int32, officeID int32) {
	t.Helper()
	ctx := context.Background()
	org, err := query.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:      randomName(),
		IsOrganic: false,
	})
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	office, err := query.CreateOffice(ctx, db.CreateOfficeParams{
		Name:           randomName(),
		OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("failed to create office: %v", err)
	}
	return org.ID, office.ID
}

// --- Tests ---

func TestListContacts(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("pagination returns max 15 per page", func(t *testing.T) {
		t.Parallel()
		orgID, officeID := seedOrgAndOffice(t)
		prefix := fmt.Sprintf("PagCon%d", time.Now().UnixNano())
		for i := 1; i <= 18; i++ {
			p := validCreateContactParams(&officeID, nil)
			p.FirstName = fmt.Sprintf("%s_%02d", prefix, i)
			p.Email = fmt.Sprintf("%s_%02d@test.com", prefix, i)
			_, err := s.CreateContact(ctx, p)
			if err != nil {
				t.Fatalf("failed to create contact %d: %v", i, err)
			}
			_ = orgID
		}

		page1, err := s.ListContacts(ctx, &ListContactsRequest{Search: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page1.Contacts) != 15 {
			t.Fatalf("expected 15 contacts on page 1, got %d", len(page1.Contacts))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := s.ListContacts(ctx, &ListContactsRequest{Search: prefix, Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page2.Contacts) != 3 {
			t.Fatalf("expected 3 contacts on page 2, got %d", len(page2.Contacts))
		}
		if page2.Total != 18 {
			t.Fatalf("expected total 18, got %d", page2.Total)
		}

		page1IDs := make(map[int32]bool)
		for _, c := range page1.Contacts {
			page1IDs[c.ID] = true
		}
		for _, c := range page2.Contacts {
			if page1IDs[c.ID] {
				t.Fatalf("contact %d appeared on both pages", c.ID)
			}
		}
	})

	t.Run("empty page returns no results with same total", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		prefix := fmt.Sprintf("EmptyPgCon%d", time.Now().UnixNano())
		for i := 1; i <= 18; i++ {
			p := validCreateContactParams(&officeID, nil)
			p.FirstName = fmt.Sprintf("%s_%02d", prefix, i)
			p.Email = fmt.Sprintf("%s_%02d@test.com", prefix, i)
			_, err := s.CreateContact(ctx, p)
			if err != nil {
				t.Fatalf("failed to create contact %d: %v", i, err)
			}
		}

		resp, err := s.ListContacts(ctx, &ListContactsRequest{Search: prefix, Page: 3})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Contacts) != 0 {
			t.Fatalf("expected 0 contacts on page 3, got %d", len(resp.Contacts))
		}
		if resp.Total != 18 {
			t.Fatalf("expected total 18, got %d", resp.Total)
		}
	})

	t.Run("filters by search name substring", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		unique := fmt.Sprintf("UniqueSearchCon%d", time.Now().UnixNano())
		p := validCreateContactParams(&officeID, nil)
		p.FirstName = unique
		p.Email = fmt.Sprintf("%s@test.com", unique)
		_, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		resp, err := s.ListContacts(ctx, &ListContactsRequest{Search: unique, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Contacts) != 1 {
			t.Fatalf("expected 1 result, got %d", len(resp.Contacts))
		}
		if resp.Contacts[0].FirstName != unique {
			t.Fatalf("expected firstName %q, got %q", unique, resp.Contacts[0].FirstName)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("filters by officeId", func(t *testing.T) {
		t.Parallel()
		orgID, officeA := seedOrgAndOffice(t)
		officeB, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name:           fmt.Sprintf("FilterOfficeB_%d", time.Now().UnixNano()),
			OrganizationID: orgID,
		})
		if err != nil {
			t.Fatalf("failed to create officeB: %v", err)
		}

		prefix := fmt.Sprintf("OffFilter%d", time.Now().UnixNano())
		pA := validCreateContactParams(&officeA, nil)
		pA.FirstName = prefix + "_A"
		pA.Email = fmt.Sprintf("%s_a@test.com", prefix)
		_, err = s.CreateContact(ctx, pA)
		if err != nil {
			t.Fatalf("failed to create contact A: %v", err)
		}

		pB := validCreateContactParams(&officeB.ID, nil)
		pB.FirstName = prefix + "_B"
		pB.Email = fmt.Sprintf("%s_b@test.com", prefix)
		_, err = s.CreateContact(ctx, pB)
		if err != nil {
			t.Fatalf("failed to create contact B: %v", err)
		}

		resp, err := s.ListContacts(ctx, &ListContactsRequest{Search: prefix, OfficeID: officeA, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Contacts) != 1 {
			t.Fatalf("expected 1 contact for officeA, got %d", len(resp.Contacts))
		}
		if resp.Contacts[0].FirstName != prefix+"_A" {
			t.Fatalf("expected contact A, got %q", resp.Contacts[0].FirstName)
		}
	})

	t.Run("filters by orgId", func(t *testing.T) {
		t.Parallel()
		orgA, officeA := seedOrgAndOffice(t)
		orgB, _ := seedOrgAndOffice(t)

		prefix := fmt.Sprintf("OrgFilter%d", time.Now().UnixNano())

		// org-level contact for orgA
		pOrgA := validCreateContactParams(nil, &orgA)
		pOrgA.FirstName = prefix + "_OrgA"
		pOrgA.Email = fmt.Sprintf("%s_orga@test.com", prefix)
		_, err := s.CreateContact(ctx, pOrgA)
		if err != nil {
			t.Fatalf("failed to create org contact A: %v", err)
		}

		// office-level contact for orgA (should NOT appear when filtering by orgId)
		pOffA := validCreateContactParams(&officeA, nil)
		pOffA.FirstName = prefix + "_OffA"
		pOffA.Email = fmt.Sprintf("%s_offa@test.com", prefix)
		_, err = s.CreateContact(ctx, pOffA)
		if err != nil {
			t.Fatalf("failed to create office contact A: %v", err)
		}

		// org-level contact for orgB
		pOrgB := validCreateContactParams(nil, &orgB)
		pOrgB.FirstName = prefix + "_OrgB"
		pOrgB.Email = fmt.Sprintf("%s_orgb@test.com", prefix)
		_, err = s.CreateContact(ctx, pOrgB)
		if err != nil {
			t.Fatalf("failed to create org contact B: %v", err)
		}

		resp, err := s.ListContacts(ctx, &ListContactsRequest{Search: prefix, OrgID: orgA, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Contacts) != 1 {
			t.Fatalf("expected 1 org-level contact for orgA, got %d", len(resp.Contacts))
		}
		if resp.Contacts[0].FirstName != prefix+"_OrgA" {
			t.Fatalf("expected OrgA contact, got %q", resp.Contacts[0].FirstName)
		}
	})

	t.Run("validation rejects page 0", func(t *testing.T) {
		t.Parallel()
		p := ListContactsRequest{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("returns error when list db fails", func(t *testing.T) {
		t.Parallel()
		q, s := contactMockService(t)
		q.EXPECT().ListContacts(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListContacts(ctx, &ListContactsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count db fails", func(t *testing.T) {
		t.Parallel()
		q, s := contactMockService(t)
		q.EXPECT().ListContacts(gomock.Any(), gomock.Any()).Return([]db.Contact{}, nil)
		q.EXPECT().CountContacts(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.ListContacts(ctx, &ListContactsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateContact(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates contact with officeId", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		p := validCreateContactParams(&officeID, nil)
		p.Email = fmt.Sprintf("create_office_%d@test.com", time.Now().UnixNano())

		resp, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.FirstName != p.FirstName {
			t.Fatalf("expected firstName %q, got %q", p.FirstName, resp.FirstName)
		}
		if resp.LastName != p.LastName {
			t.Fatalf("expected lastName %q, got %q", p.LastName, resp.LastName)
		}
		if resp.Role != p.Role {
			t.Fatalf("expected role %q, got %q", p.Role, resp.Role)
		}
		if resp.Cellphone != p.Cellphone {
			t.Fatalf("expected cellphone %q, got %q", p.Cellphone, resp.Cellphone)
		}
		if resp.Email != p.Email {
			t.Fatalf("expected email %q, got %q", p.Email, resp.Email)
		}
		if resp.OfficeID == nil || *resp.OfficeID != officeID {
			t.Fatalf("expected officeId %d, got %v", officeID, resp.OfficeID)
		}
		if resp.OrganizationID != nil {
			t.Fatalf("expected nil organizationId, got %v", resp.OrganizationID)
		}
	})

	t.Run("creates contact with organizationId", func(t *testing.T) {
		t.Parallel()
		orgID, _ := seedOrgAndOffice(t)
		p := validCreateContactParams(nil, &orgID)
		p.Email = fmt.Sprintf("create_org_%d@test.com", time.Now().UnixNano())

		resp, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.OrganizationID == nil || *resp.OrganizationID != orgID {
			t.Fatalf("expected organizationId %d, got %v", orgID, resp.OrganizationID)
		}
		if resp.OfficeID != nil {
			t.Fatalf("expected nil officeId, got %v", resp.OfficeID)
		}
	})

	t.Run("validation rejects both officeId and organizationId", func(t *testing.T) {
		t.Parallel()
		officeID := int32(1)
		orgID := int32(1)
		p := validCreateContactParams(&officeID, &orgID)

		err := p.Validate()
		wantErr := api_errors.NewErrorWithDetail(
			errs.InvalidArgument,
			"Exactly one of officeId or organizationId must be provided",
			api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue},
		)
		api_errors.AssertApiError(t, wantErr, err)
	})

	t.Run("validation rejects neither officeId nor organizationId", func(t *testing.T) {
		t.Parallel()
		p := validCreateContactParams(nil, nil)

		err := p.Validate()
		wantErr := api_errors.NewErrorWithDetail(
			errs.InvalidArgument,
			"Exactly one of officeId or organizationId must be provided",
			api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue},
		)
		api_errors.AssertApiError(t, wantErr, err)
	})

	t.Run("validation rejects blank firstName", func(t *testing.T) {
		t.Parallel()
		officeID := int32(1)
		p := validCreateContactParams(&officeID, nil)
		p.FirstName = ""
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects blank lastName", func(t *testing.T) {
		t.Parallel()
		officeID := int32(1)
		p := validCreateContactParams(&officeID, nil)
		p.LastName = ""
		api_errors.AssertApiError(t, invalidValueErr("lastName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		officeID := int32(1)
		p := validCreateContactParams(&officeID, nil)
		p.Email = "not-an-email"
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := contactMockService(t)
		q.EXPECT().CreateContact(gomock.Any(), gomock.Any()).Return(db.Contact{}, errors.New("db error"))

		officeID := int32(1)
		_, err := s.CreateContact(ctx, validCreateContactParams(&officeID, nil))
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateContact(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("updates only provided fields", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		p := validCreateContactParams(&officeID, nil)
		p.Email = fmt.Sprintf("update_partial_%d@test.com", time.Now().UnixNano())
		created, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		newFirst := "UpdatedFirst"
		resp, err := s.UpdateContact(ctx, created.ID, UpdateContactRequest{FirstName: &newFirst})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.FirstName != "UpdatedFirst" {
			t.Fatalf("expected firstName %q, got %q", "UpdatedFirst", resp.FirstName)
		}
		// Unchanged fields
		if resp.LastName != created.LastName {
			t.Fatalf("expected lastName unchanged %q, got %q", created.LastName, resp.LastName)
		}
		if resp.Role != created.Role {
			t.Fatalf("expected role unchanged %q, got %q", created.Role, resp.Role)
		}
		if resp.Cellphone != created.Cellphone {
			t.Fatalf("expected cellphone unchanged %q, got %q", created.Cellphone, resp.Cellphone)
		}
		if resp.Email != created.Email {
			t.Fatalf("expected email unchanged %q, got %q", created.Email, resp.Email)
		}
	})

	t.Run("full update changes all fields", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		p := validCreateContactParams(&officeID, nil)
		p.Email = fmt.Sprintf("update_full_%d@test.com", time.Now().UnixNano())
		created, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		params := validUpdateContactParams()
		params.Email = ptrStr(fmt.Sprintf("updated_full_%d@test.com", time.Now().UnixNano()))
		resp, err := s.UpdateContact(ctx, created.ID, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.FirstName != *params.FirstName {
			t.Fatalf("expected firstName %q, got %q", *params.FirstName, resp.FirstName)
		}
		if resp.LastName != *params.LastName {
			t.Fatalf("expected lastName %q, got %q", *params.LastName, resp.LastName)
		}
		if resp.Role != *params.Role {
			t.Fatalf("expected role %q, got %q", *params.Role, resp.Role)
		}
		if resp.Cellphone != *params.Cellphone {
			t.Fatalf("expected cellphone %q, got %q", *params.Cellphone, resp.Cellphone)
		}
		if resp.Email != *params.Email {
			t.Fatalf("expected email %q, got %q", *params.Email, resp.Email)
		}
	})

	t.Run("returns not found for non-existent id", func(t *testing.T) {
		t.Parallel()
		_, err := s.UpdateContact(ctx, 999999, validUpdateContactParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("validation rejects blank firstName", func(t *testing.T) {
		t.Parallel()
		blank := "   "
		p := UpdateContactRequest{FirstName: &blank}
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		bad := "not-an-email"
		p := UpdateContactRequest{Email: &bad}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := contactMockService(t)
		q.EXPECT().UpdateContact(gomock.Any(), gomock.Any()).Return(db.Contact{}, errors.New("db error"))

		_, err := s.UpdateContact(ctx, 1, validUpdateContactParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestDeleteContact(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("deletes contact successfully", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		p := validCreateContactParams(&officeID, nil)
		p.Email = fmt.Sprintf("delete_%d@test.com", time.Now().UnixNano())
		created, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		err = s.DeleteContact(ctx, created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify deleted: update should return not found
		_, err = s.UpdateContact(ctx, created.ID, validUpdateContactParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := contactMockService(t)
		q.EXPECT().DeleteContact(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		err := s.DeleteContact(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func ptrStr(s string) *string {
	return &s
}
