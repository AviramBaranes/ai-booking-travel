package accounts

import (
	"context"

	contact_handlers "encore.app/services/accounts/handlers/contact_handlers"
)

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
