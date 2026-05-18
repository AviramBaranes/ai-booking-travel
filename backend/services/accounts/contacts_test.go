package accounts

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.dev/beta/errs"
)

// --- Helpers ---

func validCreateContactParams(officeID *int64, orgID *int64) contact.CreateContactParams {
	return contact.CreateContactParams{
		FirstName:      "John",
		LastName:       "Doe",
		Role:           "manager",
		Cellphone:      "0521234567",
		Email:          "john.doe@test.com",
		OfficeID:       officeID,
		OrganizationID: orgID,
	}
}

func validUpdateContactParams() contact.UpdateContactParams {
	firstName := "Jane"
	lastName := "Smith"
	role := "director"
	cellphone := "0529876543"
	email := "jane.smith@test.com"
	return contact.UpdateContactParams{
		FirstName: &firstName,
		LastName:  &lastName,
		Role:      &role,
		Cellphone: &cellphone,
		Email:     &email,
	}
}

// seedOrgAndOffice creates an org and office for use in contact tests.
func seedOrgAndOffice(t *testing.T) (orgID int64, officeID int64) {
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

		page1, err := s.ListContacts(ctx, contact.ListContactsParams{Search: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page1.Contacts) != 15 {
			t.Fatalf("expected 15 contacts on page 1, got %d", len(page1.Contacts))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := s.ListContacts(ctx, contact.ListContactsParams{Search: prefix, Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page2.Contacts) != 3 {
			t.Fatalf("expected 3 contacts on page 2, got %d", len(page2.Contacts))
		}
		if page2.Total != 18 {
			t.Fatalf("expected total 18, got %d", page2.Total)
		}

		page1IDs := make(map[int64]bool)
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

		resp, err := s.ListContacts(ctx, contact.ListContactsParams{Search: prefix, Page: 3})
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

		resp, err := s.ListContacts(ctx, contact.ListContactsParams{Search: unique, Page: 1})
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

	t.Run("lists all contact fields correctly", func(t *testing.T) {
		t.Parallel()
		orgID, _ := seedOrgAndOffice(t)
		unique := fmt.Sprintf("AllFieldsList%d", time.Now().UnixNano())

		p := validCreateContactParams(nil, &orgID)
		p.FirstName = unique
		p.LastName = "FieldsTest"
		p.Role = "tester"
		p.Cellphone = "0521111111"
		p.Email = fmt.Sprintf("%s@test.com", unique)
		p.IsPaymentResponsible = true
		_, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		resp, err := s.ListContacts(ctx, contact.ListContactsParams{Search: unique, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Contacts) != 1 {
			t.Fatalf("expected 1 result, got %d", len(resp.Contacts))
		}
		c := resp.Contacts[0]
		if c.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if c.FirstName != unique {
			t.Fatalf("expected firstName %q, got %q", unique, c.FirstName)
		}
		if c.LastName != "FieldsTest" {
			t.Fatalf("expected lastName %q, got %q", "FieldsTest", c.LastName)
		}
		if c.Role != "tester" {
			t.Fatalf("expected role %q, got %q", "tester", c.Role)
		}
		if c.Cellphone != "0521111111" {
			t.Fatalf("expected cellphone %q, got %q", "0521111111", c.Cellphone)
		}
		if c.Email != p.Email {
			t.Fatalf("expected email %q, got %q", p.Email, c.Email)
		}
		if c.OrganizationID == nil || *c.OrganizationID != orgID {
			t.Fatalf("expected organizationId %d, got %v", orgID, c.OrganizationID)
		}
		if c.OfficeID != nil {
			t.Fatalf("expected nil officeId, got %v", c.OfficeID)
		}
		if !c.IsPaymentResponsible {
			t.Fatal("expected isPaymentResponsible to be true")
		}
		if c.OrganizationName == nil || *c.OrganizationName == "" {
			t.Fatalf("expected non-empty organizationName, got %v", c.OrganizationName)
		}
		if c.OfficeName != nil {
			t.Fatalf("expected nil officeName, got %v", c.OfficeName)
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

		resp, err := s.ListContacts(ctx, contact.ListContactsParams{Search: prefix, OfficeID: officeA, Page: 1})
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

		resp, err := s.ListContacts(ctx, contact.ListContactsParams{Search: prefix, OrgID: orgA, Page: 1})
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
		p := contact.ListContactsParams{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
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
		if resp.IsPaymentResponsible {
			t.Fatal("expected isPaymentResponsible to default to false for office contact")
		}
	})

	t.Run("creates contact with organizationId", func(t *testing.T) {
		t.Parallel()
		orgID, _ := seedOrgAndOffice(t)
		p := validCreateContactParams(nil, &orgID)
		p.Email = fmt.Sprintf("create_org_%d@test.com", time.Now().UnixNano())
		p.IsPaymentResponsible = true

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
		if resp.OrganizationID == nil || *resp.OrganizationID != orgID {
			t.Fatalf("expected organizationId %d, got %v", orgID, resp.OrganizationID)
		}
		if resp.OfficeID != nil {
			t.Fatalf("expected nil officeId, got %v", resp.OfficeID)
		}
		if !resp.IsPaymentResponsible {
			t.Fatal("expected isPaymentResponsible to be true")
		}
	})

	t.Run("validation rejects both officeId and organizationId", func(t *testing.T) {
		t.Parallel()
		officeID := int64(1)
		orgID := int64(1)
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
		officeID := int64(1)
		p := validCreateContactParams(&officeID, nil)
		p.FirstName = ""
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects blank lastName", func(t *testing.T) {
		t.Parallel()
		officeID := int64(1)
		p := validCreateContactParams(&officeID, nil)
		p.LastName = ""
		api_errors.AssertApiError(t, invalidValueErr("lastName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		officeID := int64(1)
		p := validCreateContactParams(&officeID, nil)
		p.Email = "not-an-email"
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
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
		resp, err := s.UpdateContact(ctx, created.ID, contact.UpdateContactParams{FirstName: &newFirst})
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
		orgID, _ := seedOrgAndOffice(t)
		p := validCreateContactParams(nil, &orgID)
		p.Email = fmt.Sprintf("update_full_%d@test.com", time.Now().UnixNano())
		created, err := s.CreateContact(ctx, p)
		if err != nil {
			t.Fatalf("failed to create contact: %v", err)
		}

		isTrue := true
		params := validUpdateContactParams()
		params.Email = ptrStr(fmt.Sprintf("updated_full_%d@test.com", time.Now().UnixNano()))
		params.IsPaymentResponsible = &isTrue
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
		if !resp.IsPaymentResponsible {
			t.Fatal("expected isPaymentResponsible to be true after full update")
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
		p := contact.UpdateContactParams{FirstName: &blank}
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		bad := "not-an-email"
		p := contact.UpdateContactParams{Email: &bad}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
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

}

func ptrStr(s string) *string {
	return &s
}

// --- Helpers ---

func seedOrg(t *testing.T, isOrganic bool) db.Organization {
	t.Helper()
	var icountClientID *int32
	if isOrganic {
		id := int32(999)
		icountClientID = &id
	}
	org, err := query.CreateOrganization(context.Background(), db.CreateOrganizationParams{
		Name: randomName(), IsOrganic: isOrganic, IcountClientID: icountClientID,
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	return org
}

func seedOffice(t *testing.T, orgID int64) db.CreateOfficeRow {
	t.Helper()
	o, err := query.CreateOffice(context.Background(), db.CreateOfficeParams{
		Name: randomName(), OrganizationID: orgID,
	})
	if err != nil {
		t.Fatalf("create office: %v", err)
	}
	return o
}

// seedContact creates a contact attached to either an office or an organization.
// Exactly one of officeID or orgID must be non-nil.
func seedContact(t *testing.T, officeID, orgID *int64, isPaymentResponsible bool) db.Contact {
	t.Helper()
	c, err := query.CreateContact(context.Background(), db.CreateContactParams{
		FirstName:            randomName(),
		LastName:             randomName(),
		Role:                 "billing",
		Cellphone:            "0521234567",
		Email:                randomName() + "@test.com",
		OfficeID:             officeID,
		OrganizationID:       orgID,
		IsPaymentResponsible: isPaymentResponsible,
	})
	if err != nil {
		t.Fatalf("create contact: %v", err)
	}
	return c
}

// sortBillingContact returns a copy of c with Offices and their Agents sorted by ID,
// allowing order-independent comparison of the EP's response.
func sortBillingContact(c contact.BillingContact) contact.BillingContact {
	offices := append([]contact.Office(nil), c.Offices...)
	sort.Slice(offices, func(i, j int) bool { return offices[i].ID < offices[j].ID })
	for i := range offices {
		agents := append([]contact.Agent(nil), offices[i].Agents...)
		sort.Slice(agents, func(a, b int) bool { return agents[a].ID < agents[b].ID })
		offices[i].Agents = agents
	}
	c.Offices = offices
	return c
}

// --- Tests ---

func TestGetBillingContacts(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("returns empty when nil ids or agent has no payment-responsible contact", func(t *testing.T) {
		t.Parallel()
		// Nil ID slice short-circuits to an empty response.
		resp, err := s.GetBillingContacts(ctx, contact.GetBillingContactsParams{AgentsIDs: nil})
		if err != nil || len(resp.Contacts) != 0 {
			t.Fatalf("nil ids: err=%v, contacts=%d", err, len(resp.Contacts))
		}

		// Agent whose only associated contact is not payment-responsible yields no results.
		org := seedOrg(t, false)
		office := seedOffice(t, org.ID)
		agent := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), office.ID)
		seedContact(t, &office.ID, nil, false)

		resp, err = s.GetBillingContacts(ctx, contact.GetBillingContactsParams{AgentsIDs: []int64{agent.ID}})
		if err != nil || len(resp.Contacts) != 0 {
			t.Fatalf("no payresp: err=%v, contacts=%d", err, len(resp.Contacts))
		}
	})

	// Exercises all join rules and grouping shapes against a single request containing
	// many agent IDs, mirroring how the EP is invoked in production.
	t.Run("matches join rules and groups across many agents in one call", func(t *testing.T) {
		t.Parallel()

		// Organization A: organic, two offices, one org-level payment contact (expected match).
		// Also seeded: an office-level contact (excluded by the organic + office mismatch)
		// and a non-payment-responsible org-level contact (excluded by the flag).
		orgA := seedOrg(t, true)
		officeA1, officeA2 := seedOffice(t, orgA.ID), seedOffice(t, orgA.ID)
		agentA1 := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeA1.ID)
		agentA2 := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeA2.ID)
		contactA := seedContact(t, nil, &orgA.ID, true)
		seedContact(t, &officeA1.ID, nil, true)
		seedContact(t, nil, &orgA.ID, false)

		// Organization B: non-organic, two offices.
		//   Office B1: two agents and one office-level payment contact, expected to produce
		//              a single BillingContact grouping both agents under one office.
		//   Office B2: one agent and one office-level payment contact, producing a
		//              separate BillingContact.
		// Also seeded: an org-level payment contact (excluded by the non-organic + org mismatch).
		orgB := seedOrg(t, false)
		officeB1, officeB2 := seedOffice(t, orgB.ID), seedOffice(t, orgB.ID)
		agentB1a := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeB1.ID)
		agentB1b := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeB1.ID)
		agentB2 := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeB2.ID)
		contactB1 := seedContact(t, &officeB1.ID, nil, true)
		contactB2 := seedContact(t, &officeB2.ID, nil, true)
		seedContact(t, nil, &orgB.ID, true)

		// Organization C: non-organic, single office with two payment-responsible contacts,
		// expected to produce two BillingContacts that share the same office.
		// Also seeded: an agent omitted from the request, which must not appear in the response.
		orgC := seedOrg(t, false)
		officeC := seedOffice(t, orgC.ID)
		agentC := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeC.ID)
		agentCExcluded := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), officeC.ID)
		contactC1 := seedContact(t, &officeC.ID, nil, true)
		contactC2 := seedContact(t, &officeC.ID, nil, true)

		// Non-agent users whose IDs are passed in the request must be excluded by the role filter.
		admin, err := query.CreateStaffUser(ctx, db.CreateStaffUserParams{
			Role:      db.UserRoleAdmin,
			FirstName: "A", LastName: "A", Email: randomName() + "@test.com", PasswordHash: "h",
		})
		if err != nil {
			t.Fatalf("create admin: %v", err)
		}
		t.Cleanup(func() { _ = query.DeleteUser(ctx, admin.ID) })
		customer, err := query.CreateCustomer(ctx, db.CreateCustomerParams{
			FirstName: "C", LastName: "C", Email: randomName() + "@test.com", PasswordHash: "h",
		})
		if err != nil {
			t.Fatalf("create customer: %v", err)
		}
		t.Cleanup(func() { _ = query.DeleteUser(ctx, customer.ID) })

		// Single request carrying every agent, admin and customer ID under test.
		// agentCExcluded is intentionally omitted to verify scoping by request IDs.
		ids := []int64{
			agentA1.ID, agentA2.ID,
			agentB1a.ID, agentB1b.ID, agentB2.ID,
			agentC.ID,
			admin.ID, customer.ID,
		}
		resp, err := s.GetBillingContacts(ctx, contact.GetBillingContactsParams{AgentsIDs: ids})
		if err != nil {
			t.Fatalf("GetBillingContacts: %v", err)
		}

		// seedAgent constructs every agent with the same first and last name,
		// so all expected Agent.Name values resolve to this constant.
		const agentName = "Test Agent"
		want := map[string]contact.BillingContact{
			contactA.Email: {
				ContactName: contactA.FirstName + " " + contactA.LastName, ContactEmail: contactA.Email,
				OrganizationID: orgA.ID, OrganizationName: orgA.Name, IsOrganic: true,
				Offices: []contact.Office{
					{ID: officeA1.ID, Name: officeA1.Name, Agents: []contact.Agent{{ID: agentA1.ID, Name: agentName}}},
					{ID: officeA2.ID, Name: officeA2.Name, Agents: []contact.Agent{{ID: agentA2.ID, Name: agentName}}},
				},
			},
			contactB1.Email: {
				ContactName: contactB1.FirstName + " " + contactB1.LastName, ContactEmail: contactB1.Email,
				OrganizationID: orgB.ID, OrganizationName: orgB.Name, IsOrganic: false,
				Offices: []contact.Office{{ID: officeB1.ID, Name: officeB1.Name, Agents: []contact.Agent{
					{ID: agentB1a.ID, Name: agentName}, {ID: agentB1b.ID, Name: agentName},
				}}},
			},
			contactB2.Email: {
				ContactName: contactB2.FirstName + " " + contactB2.LastName, ContactEmail: contactB2.Email,
				OrganizationID: orgB.ID, OrganizationName: orgB.Name, IsOrganic: false,
				Offices: []contact.Office{{ID: officeB2.ID, Name: officeB2.Name, Agents: []contact.Agent{{ID: agentB2.ID, Name: agentName}}}},
			},
			contactC1.Email: {
				ContactName: contactC1.FirstName + " " + contactC1.LastName, ContactEmail: contactC1.Email,
				OrganizationID: orgC.ID, OrganizationName: orgC.Name, IsOrganic: false,
				Offices: []contact.Office{{ID: officeC.ID, Name: officeC.Name, Agents: []contact.Agent{{ID: agentC.ID, Name: agentName}}}},
			},
			contactC2.Email: {
				ContactName: contactC2.FirstName + " " + contactC2.LastName, ContactEmail: contactC2.Email,
				OrganizationID: orgC.ID, OrganizationName: orgC.Name, IsOrganic: false,
				Offices: []contact.Office{{ID: officeC.ID, Name: officeC.Name, Agents: []contact.Agent{{ID: agentC.ID, Name: agentName}}}},
			},
		}

		if len(resp.Contacts) != len(want) {
			t.Fatalf("got %d contacts, want %d", len(resp.Contacts), len(want))
		}

		got := make(map[string]contact.BillingContact, len(resp.Contacts))
		for _, c := range resp.Contacts {
			got[c.ContactEmail] = c
			// Verify that excluded user IDs never appear under any contact.
			for _, o := range c.Offices {
				for _, a := range o.Agents {
					if a.ID == agentCExcluded.ID || a.ID == admin.ID || a.ID == customer.ID {
						t.Errorf("excluded user id %d leaked into contact %q", a.ID, c.ContactEmail)
					}
				}
			}
		}

		for email, w := range want {
			g, ok := got[email]
			if !ok {
				t.Errorf("missing expected contact %q", email)
				continue
			}
			if diff := compareBillingContact(sortBillingContact(w), sortBillingContact(g)); diff != "" {
				t.Errorf("contact %q mismatch: %s", email, diff)
			}
		}
	})
}

// compareBillingContact reports the first field-level difference between want and got,
// or returns an empty string if the two contacts are equal.
func compareBillingContact(want, got contact.BillingContact) string {
	if want.ContactName != got.ContactName {
		return "ContactName: want " + want.ContactName + ", got " + got.ContactName
	}
	if want.OrganizationID != got.OrganizationID || want.OrganizationName != got.OrganizationName || want.IsOrganic != got.IsOrganic {
		return "org fields differ"
	}
	if len(want.Offices) != len(got.Offices) {
		return "offices length differs"
	}
	for i := range want.Offices {
		w, g := want.Offices[i], got.Offices[i]
		if w.ID != g.ID || w.Name != g.Name || len(w.Agents) != len(g.Agents) {
			return "office differs"
		}
		for j := range w.Agents {
			if w.Agents[j] != g.Agents[j] {
				return "agent differs"
			}
		}
	}
	return ""
}
