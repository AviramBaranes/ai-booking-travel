package accounts

import (
	"context"

	organization "encore.app/services/accounts/handlers/organization"
)

// --- Endpoints ---

// ListOrganizations lists organizations with optional search and pagination.
//
//encore:api auth method=GET path=/organizations tag:admin
func (s *Service) ListOrganizations(ctx context.Context, p organization.ListOrganizationsParams) (*organization.ListOrganizationsResponse, error) {
	handler := organization.NewOrganizationService(s.query)
	return handler.ListOrganizations(ctx, p)
}

// CreateOrganization creates a new organization.
//
//encore:api auth method=POST path=/organizations tag:admin
func (s *Service) CreateOrganization(ctx context.Context, p organization.CreateOrganizationParams) (*organization.OrganizationResponse, error) {
	handler := organization.NewOrganizationService(s.query)
	return handler.CreateOrganization(ctx, p)
}

// UpdateOrganization updates an existing organization.
//
//encore:api auth method=PUT path=/organizations/:id tag:admin
func (s *Service) UpdateOrganization(ctx context.Context, id int64, p organization.UpdateOrganizationParams) (*organization.OrganizationResponse, error) {
	handler := organization.NewOrganizationService(s.query)
	return handler.UpdateOrganization(ctx, id, p)
}

// ListOrganicOrganizations lists all organic organizations for accountant use.
//
//encore:api auth method=GET path=/organic-organizations tag:accountant
func (s *Service) ListOrganicOrganizations(ctx context.Context) (*organization.ListOrganicOrganizationResponse, error) {
	handler := organization.NewOrganizationService(s.query)
	return handler.ListOrganicOrganizations(ctx)
}

// encore:api private
func (s *Service) UpdateOrganizationBalanceDue(ctx context.Context, p organization.UpdateOrganizationBalanceDueParams) error {
	handler := organization.NewOrganizationService(s.query)
	return handler.UpdateOrganizationBalanceDue(ctx, p)
}
