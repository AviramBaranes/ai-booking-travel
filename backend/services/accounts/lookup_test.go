package accounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"encore.app/services/accounts/db"
	"encore.app/services/accounts/handlers/lookup"
)

func TestGetAccountsLookup(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	seedLookupRows := func(t *testing.T) (db.Organization, db.CreateOfficeRow, db.CreateAgentRow, db.CreateStaffUserRow) {
		t.Helper()
		unique := time.Now().UnixNano()

		org, err := query.CreateOrganization(ctx, db.CreateOrganizationParams{
			Name:      fmt.Sprintf("lookup_org_%d", unique),
			IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("failed to create organization: %v", err)
		}

		office, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name:           fmt.Sprintf("lookup_office_%d", unique),
			OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("failed to create office: %v", err)
		}

		phone := randomIsraeliPhoneNumber()
		agent, err := query.CreateAgent(ctx, db.CreateAgentParams{
			FirstName:    "Lookup",
			LastName:     "Agent",
			Email:        fmt.Sprintf("lookup_agent_%d@test.com", unique),
			PhoneNumber:  &phone,
			PasswordHash: "hash",
			OfficeID:     &office.ID,
		})
		if err != nil {
			t.Fatalf("failed to create agent: %v", err)
		}
		t.Cleanup(func() { _ = query.DeleteUser(ctx, agent.ID) })

		admin, err := query.CreateStaffUser(ctx, db.CreateStaffUserParams{
			Role:         db.UserRoleAdmin,
			FirstName:    "Lookup",
			LastName:     "Admin",
			Email:        fmt.Sprintf("lookup_admin_%d@test.com", unique),
			PasswordHash: "hash",
		})
		if err != nil {
			t.Fatalf("failed to create admin: %v", err)
		}
		t.Cleanup(func() { _ = query.DeleteUser(ctx, admin.ID) })

		return org, office, agent, admin
	}

	t.Run("returns seeded account names", func(t *testing.T) {
		org, office, agent, admin := seedLookupRows(t)

		resp, err := s.GetAccountsLookup(ctx, lookup.GetAccountsLookupParams{
			OrganizationIDs: []int64{org.ID},
			OfficeIDs:       []int64{office.ID},
			UserIDs:         []int64{agent.ID, admin.ID},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		assertLookupName(t, resp.Organizations, org.ID, org.Name)
		assertLookupName(t, resp.Offices, office.ID, office.Name)
		assertLookupName(t, resp.Users, agent.ID, "Lookup Agent")
		assertLookupName(t, resp.Users, admin.ID, "Lookup Admin")
	})

	t.Run("ignores missing ids", func(t *testing.T) {
		org, office, agent, _ := seedLookupRows(t)
		missingID := int64(999999999)

		resp, err := s.GetAccountsLookup(ctx, lookup.GetAccountsLookupParams{
			OrganizationIDs: []int64{org.ID, missingID},
			OfficeIDs:       []int64{office.ID, missingID},
			UserIDs:         []int64{agent.ID, missingID},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(resp.Organizations) != 1 {
			t.Fatalf("expected 1 organization, got %d", len(resp.Organizations))
		}
		if len(resp.Offices) != 1 {
			t.Fatalf("expected 1 office, got %d", len(resp.Offices))
		}
		if len(resp.Users) != 1 {
			t.Fatalf("expected 1 user, got %d", len(resp.Users))
		}
	})
}

func assertLookupName(t *testing.T, rows []lookup.AccountName, id int64, want string) {
	t.Helper()
	for _, row := range rows {
		if row.ID == id {
			if row.Name != want {
				t.Fatalf("expected name %q for id %d, got %q", want, id, row.Name)
			}
			return
		}
	}
	t.Fatalf("expected id %d in lookup response", id)
}
