package accounts

import (
	"context"

	contact "encore.app/services/accounts/handlers/contact"
)

// ListContacts lists contacts with optional filtering and pagination.
//
//encore:api auth method=GET path=/contacts tag:admin
func (s *Service) ListContacts(ctx context.Context, p contact.ListContactsParams) (*contact.ListContactsResponse, error) {
	h := contact.NewContactService(s.query)
	return h.ListContacts(ctx, p)
}

// CreateContact creates a new contact.
//
//encore:api auth method=POST path=/contacts tag:admin
func (s *Service) CreateContact(ctx context.Context, p contact.CreateContactParams) (*contact.ContactResponse, error) {
	h := contact.NewContactService(s.query)
	return h.CreateContact(ctx, p)
}

// UpdateContact updates an existing contact.
//
//encore:api auth method=PUT path=/contacts/:id tag:admin
func (s *Service) UpdateContact(ctx context.Context, id int64, p contact.UpdateContactParams) (*contact.ContactResponse, error) {
	h := contact.NewContactService(s.query)
	return h.UpdateContact(ctx, id, p)
}

// DeleteContact deletes a contact by its ID.
//
//encore:api auth method=DELETE path=/contacts/:id tag:admin
func (s *Service) DeleteContact(ctx context.Context, id int64) error {
	h := contact.NewContactService(s.query)
	return h.DeleteContact(ctx, id)
}

// encore:api private
func (s *Service) GetBillingContacts(ctx context.Context, p contact.GetBillingContactsParams) (*contact.GetBillingContactsResponse, error) {
	h := contact.NewContactService(s.query)
	return h.GetBillingContacts(ctx, p)
}

// GetIcountClientID returns the iCount client ID for a given office or organization.
//
// encore:api private
func (s *Service) GetIcountClientID(ctx context.Context, p contact.GetIcountClientIDParams) (*contact.GetIcountClientIDResponse, error) {
	h := contact.NewContactService(s.query)
	return h.GetIcountClientID(ctx, p)
}
