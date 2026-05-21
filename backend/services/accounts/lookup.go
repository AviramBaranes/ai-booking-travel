package accounts

import (
	"context"

	"encore.app/services/accounts/handlers/lookup"
)

// encore:api private
func (s *Service) GetAccountsLookup(ctx context.Context, p lookup.GetAccountsLookupParams) (*lookup.GetAccountsLookupResponse, error) {
	h := lookup.NewService(s.query)
	return h.GetAccountsLookup(ctx, p)
}
