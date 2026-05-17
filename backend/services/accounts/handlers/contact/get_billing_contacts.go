package contact

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type GetBillingContactsParams struct {
	AgentsIDs []int64
}

type GetBillingContactsResponse struct {
	Contacts []BillingContact
}

// Agent represents an agent associated with a billing contact.
type Agent struct {
	ID   int64
	Name string
}

// Office represents an office associated with a billing contact.
type Office struct {
	ID     int64
	Name   string
	Agents []Agent
}

// BillingContact represents the billing contact information for an organization billing responsible contact.
type BillingContact struct {
	ContactName      string
	ContactEmail     string
	OrganizationID   int64
	OrganizationName string
	IsOrganic        bool
	Offices          []Office
}

func (s *ContactService) GetBillingContacts(ctx context.Context, p GetBillingContactsParams) (*GetBillingContactsResponse, error) {
	rows, err := s.query.GetAgentsBillingContacts(ctx, p.AgentsIDs)
	if err != nil {
		rlog.Error("failed to get billing contacts for agents", "error", err)
		return nil, api_errors.ErrInternalError
	}

	contactsMap := createContactsMap(rows)

	contacts := make([]BillingContact, 0, len(contactsMap))
	for _, contact := range contactsMap {
		contacts = append(contacts, contact)
	}

	return &GetBillingContactsResponse{Contacts: contacts}, nil
}

func createContactsMap(rows []db.GetAgentsBillingContactsRow) map[int64]BillingContact {
	contactsMap := make(map[int64]BillingContact)
	for _, r := range rows {
		contact, exists := contactsMap[r.ContactID]
		if !exists {
			contact = BillingContact{
				ContactName:      r.ContactFirstName + " " + r.ContactLastName,
				ContactEmail:     r.Email,
				OrganizationID:   r.OrganizationID,
				OrganizationName: r.OrganizationName,
				IsOrganic:        r.IsOrganic,
			}
		}

		agent := Agent{
			ID:   r.AgentID,
			Name: r.AgentFirstName + " " + r.AgentLastName,
		}

		officeIndex := -1
		for i, o := range contact.Offices {
			if o.ID == r.OfficeID {
				officeIndex = i
				break
			}
		}

		if officeIndex == -1 {
			newOffice := Office{
				ID:     r.OfficeID,
				Name:   r.OfficeName,
				Agents: []Agent{agent},
			}
			contact.Offices = append(contact.Offices, newOffice)
		} else {
			contact.Offices[officeIndex].Agents = append(contact.Offices[officeIndex].Agents, agent)
		}

		contactsMap[r.ContactID] = contact
	}

	return contactsMap
}
