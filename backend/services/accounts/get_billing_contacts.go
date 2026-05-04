package accounts

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// GetBillingContacts is the response type for GetBillingContacts.
type GetBillingContactsResponse struct {
	Contacts []BillingContact
}

// Agent represents an agent associated with a billing contact.
type Agent struct {
	ID   int32
	Name string
}

// Office represents an office associated with a billing contact.
type Office struct {
	ID     int32
	Name   string
	Agents []Agent
}

// BillingContact represents the billing contact information for an organization billing responsible contact.
type BillingContact struct {
	ContactName      string
	ContactEmail     string
	OrganizationID   int32
	OrganizationName string
	IsOrganic        bool
	Offices          []Office
}

// GetBillingContactsRequest is the request type for GetBillingContacts.
type GetBillingContactsRequest struct {
	AgentsIDs []int32
}

// encore:api private
func (s *Service) GetBillingContacts(ctx context.Context, p *GetBillingContactsRequest) (*GetBillingContactsResponse, error) {
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

	return &GetBillingContactsResponse{
		Contacts: contacts,
	}, nil
}

// createContactsMap converts a db.ListContactsRow to a ContactResponse.
func createContactsMap(rows []db.GetAgentsBillingContactsRow) map[int32]BillingContact {
	contactsMap := make(map[int32]BillingContact)
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

// GetIcountClientIDRequest is the request type for GetIcountClientID.
type GetIcountClientIDParams struct {
	OfficeID       *int32
	OrganizationID *int32
}

// GetIcountClientIDResponse is the response type for GetIcountClientID.
type GetIcountClientIDResponse struct {
	ClientID int32
}

// GetIcountClientID returns the iCount client ID for a given office or organization.
//
// encore:api private
func (s *Service) GetIcountClientID(ctx context.Context, p *GetIcountClientIDParams) (*GetIcountClientIDResponse, error) {
	if p.OfficeID != nil {
		icountClientID, err := s.query.GetOfficeIcountClientID(ctx, *p.OfficeID)
		if err != nil {
			rlog.Error("failed to get iCount client ID for office", "error", err, "office_id", p.OfficeID)
			return nil, api_errors.ErrInternalError
		}
		if icountClientID == nil {
			return nil, api_errors.ErrNotFound
		}
		return &GetIcountClientIDResponse{ClientID: *icountClientID}, nil
	}

	if p.OrganizationID != nil {
		icountClientID, err := s.query.GetOrganizationIcountClientID(ctx, *p.OrganizationID)
		if err != nil {
			rlog.Error("failed to get iCount client ID for organization", "error", err, "organization_id", p.OrganizationID)
			return nil, api_errors.ErrInternalError
		}
		if icountClientID == nil {
			return nil, api_errors.ErrNotFound
		}
		return &GetIcountClientIDResponse{ClientID: *icountClientID}, nil
	}

	return nil, api_errors.ErrInvalidValue
}
