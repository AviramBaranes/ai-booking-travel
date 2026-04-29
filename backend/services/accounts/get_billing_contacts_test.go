package accounts

import (
	"context"
	"errors"
	"sort"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func billingMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func seedOrg(t *testing.T, isOrganic bool) db.Organization {
	t.Helper()
	org, err := query.CreateOrganization(context.Background(), db.CreateOrganizationParams{
		Name: randomName(), IsOrganic: isOrganic,
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	return org
}

func seedOffice(t *testing.T, orgID int32) db.Office {
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
func seedContact(t *testing.T, officeID, orgID *int32, isPaymentResponsible bool) db.Contact {
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
func sortBillingContact(c BillingContact) BillingContact {
	offices := append([]Office(nil), c.Offices...)
	sort.Slice(offices, func(i, j int) bool { return offices[i].ID < offices[j].ID })
	for i := range offices {
		agents := append([]Agent(nil), offices[i].Agents...)
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

	t.Run("returns internal error when db fails", func(t *testing.T) {
		t.Parallel()
		q, ms := billingMockService(t)
		q.EXPECT().GetAgentsBillingContacts(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db error"))

		_, err := ms.GetBillingContacts(ctx, &GetBillingContactsRequest{AgentsIDs: []int32{1}})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns empty when nil ids or agent has no payment-responsible contact", func(t *testing.T) {
		t.Parallel()
		// Nil ID slice short-circuits to an empty response.
		resp, err := s.GetBillingContacts(ctx, &GetBillingContactsRequest{AgentsIDs: nil})
		if err != nil || len(resp.Contacts) != 0 {
			t.Fatalf("nil ids: err=%v, contacts=%d", err, len(resp.Contacts))
		}

		// Agent whose only associated contact is not payment-responsible yields no results.
		org := seedOrg(t, false)
		office := seedOffice(t, org.ID)
		agent := seedAgent(t, s, randomName()+"@test.com", randomIsraeliPhoneNumber(), office.ID)
		seedContact(t, &office.ID, nil, false)

		resp, err = s.GetBillingContacts(ctx, &GetBillingContactsRequest{AgentsIDs: []int32{agent.ID}})
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
		admin, err := query.CreateAdmin(ctx, db.CreateAdminParams{
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
		ids := []int32{
			agentA1.ID, agentA2.ID,
			agentB1a.ID, agentB1b.ID, agentB2.ID,
			agentC.ID,
			admin.ID, customer.ID,
		}
		resp, err := s.GetBillingContacts(ctx, &GetBillingContactsRequest{AgentsIDs: ids})
		if err != nil {
			t.Fatalf("GetBillingContacts: %v", err)
		}

		// seedAgent constructs every agent with the same first and last name,
		// so all expected Agent.Name values resolve to this constant.
		const agentName = "Test Agent"
		want := map[string]BillingContact{
			contactA.Email: {
				ContactName: contactA.FirstName + " " + contactA.LastName, ContactEmail: contactA.Email,
				OrganizationID: orgA.ID, OrganizationName: orgA.Name, IsOrganic: true,
				Offices: []Office{
					{ID: officeA1.ID, Name: officeA1.Name, Agents: []Agent{{ID: agentA1.ID, Name: agentName}}},
					{ID: officeA2.ID, Name: officeA2.Name, Agents: []Agent{{ID: agentA2.ID, Name: agentName}}},
				},
			},
			contactB1.Email: {
				ContactName: contactB1.FirstName + " " + contactB1.LastName, ContactEmail: contactB1.Email,
				OrganizationID: orgB.ID, OrganizationName: orgB.Name, IsOrganic: false,
				Offices: []Office{{ID: officeB1.ID, Name: officeB1.Name, Agents: []Agent{
					{ID: agentB1a.ID, Name: agentName}, {ID: agentB1b.ID, Name: agentName},
				}}},
			},
			contactB2.Email: {
				ContactName: contactB2.FirstName + " " + contactB2.LastName, ContactEmail: contactB2.Email,
				OrganizationID: orgB.ID, OrganizationName: orgB.Name, IsOrganic: false,
				Offices: []Office{{ID: officeB2.ID, Name: officeB2.Name, Agents: []Agent{{ID: agentB2.ID, Name: agentName}}}},
			},
			contactC1.Email: {
				ContactName: contactC1.FirstName + " " + contactC1.LastName, ContactEmail: contactC1.Email,
				OrganizationID: orgC.ID, OrganizationName: orgC.Name, IsOrganic: false,
				Offices: []Office{{ID: officeC.ID, Name: officeC.Name, Agents: []Agent{{ID: agentC.ID, Name: agentName}}}},
			},
			contactC2.Email: {
				ContactName: contactC2.FirstName + " " + contactC2.LastName, ContactEmail: contactC2.Email,
				OrganizationID: orgC.ID, OrganizationName: orgC.Name, IsOrganic: false,
				Offices: []Office{{ID: officeC.ID, Name: officeC.Name, Agents: []Agent{{ID: agentC.ID, Name: agentName}}}},
			},
		}

		if len(resp.Contacts) != len(want) {
			t.Fatalf("got %d contacts, want %d", len(resp.Contacts), len(want))
		}

		got := make(map[string]BillingContact, len(resp.Contacts))
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
func compareBillingContact(want, got BillingContact) string {
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
