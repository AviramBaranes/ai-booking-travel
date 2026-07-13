package accounts

import (
	"context"

	"encore.app/services/accounts/handlers/customer"
)

// encore:api auth path=/customers method=GET tag:admin
func (s *Service) ListCustomers(ctx context.Context, p customer.ListCustomersParams) (*customer.ListCustomersResponse, error) {
	cs := customer.NewCustomerService(s.query)
	return cs.ListCustomers(ctx, p)
}

// encore:api auth path=/update/me method=PUT tag:customer
func (s *Service) UpdateCustomer(ctx context.Context, p customer.UpdateCustomerParams) error {
	authData := GetAuthData()
	cs := customer.NewCustomerService(s.query)
	return cs.UpdateCustomer(ctx, p, authData.UserID)
}
