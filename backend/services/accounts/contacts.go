package accounts

import (
	"context"

	contact_handlers "encore.app/services/accounts/handlers/contact_handlers"
)

// ListContacts lists contacts with optional filtering and pagination.
//
//encore:api auth method=GET path=/contacts tag:admin
func (s *Service) ListContacts(ctx context.Context, p contact_handlers.ListContactsParams) (*contact_handlers.ListContactsResponse, error) {
	h := contact_handlers.NewContactService(s.query)
	return h.ListContacts(ctx, p)
}

// CreateContact creates a new contact.
//
//encore:api auth method=POST path=/contacts tag:admin
func (s *Service) CreateContact(ctx context.Context, p contact_handlers.CreateContactParams) (*contact_handlers.ContactResponse, error) {
	h := contact_handlers.NewContactService(s.query)
	return h.CreateContact(ctx, p)
}

// UpdateContact updates an existing contact.
//
//encore:api auth method=PUT path=/contacts/:id tag:admin
func (s *Service) UpdateContact(ctx context.Context, id int64, p contact_handlers.UpdateContactParams) (*contact_handlers.ContactResponse, error) {
	h := contact_handlers.NewContactService(s.query)
	return h.UpdateContact(ctx, id, p)
}

// DeleteContact deletes a contact by its ID.
//
//encore:api auth method=DELETE path=/contacts/:id tag:admin
func (s *Service) DeleteContact(ctx context.Context, id int64) error {
	h := contact_handlers.NewContactService(s.query)
	return h.DeleteContact(ctx, id)
}

// encore:api private
func (s *Service) GetBillingContacts(ctx context.Context, p contact_handlers.GetBillingContactsParams) (*contact_handlers.GetBillingContactsResponse, error) {
	h := contact_handlers.NewContactService(s.query)
	return h.GetBillingContacts(ctx, p)
}

// GetIcountClientID returns the iCount client ID for a given office or organization.
//
// encore:api private
func (s *Service) GetIcountClientID(ctx context.Context, p contact_handlers.GetIcountClientIDParams) (*contact_handlers.GetIcountClientIDResponse, error) {
	h := contact_handlers.NewContactService(s.query)
	return h.GetIcountClientID(ctx, p)
}
