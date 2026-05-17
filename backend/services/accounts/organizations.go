package accounts

import (
	"context"

	"encore.app/services/accounts/handlers/organization_handlers"
)

// --- Endpoints ---

// ListOrganizations lists organizations with optional search and pagination.
//
//encore:api auth method=GET path=/organizations tag:admin
func (s *Service) ListOrganizations(ctx context.Context, p organization_handlers.ListOrganizationsParams) (*organization_handlers.ListOrganizationsResponse, error) {
	handler := organization_handlers.NewOrganizationService(s.query)
	return handler.ListOrganizations(ctx, p)
}

// CreateOrganization creates a new organization.
//
//encore:api auth method=POST path=/organizations tag:admin
func (s *Service) CreateOrganization(ctx context.Context, p organization_handlers.CreateOrganizationParams) (*organization_handlers.OrganizationResponse, error) {
	handler := organization_handlers.NewOrganizationService(s.query)
	return handler.CreateOrganization(ctx, p)
}

// UpdateOrganization updates an existing organization.
//
//encore:api auth method=PUT path=/organizations/:id tag:admin
func (s *Service) UpdateOrganization(ctx context.Context, id int64, p organization_handlers.UpdateOrganizationParams) (*organization_handlers.OrganizationResponse, error) {
	handler := organization_handlers.NewOrganizationService(s.query)
	return handler.UpdateOrganization(ctx, id, p)
}

// ListOrganicOrganizations lists all organic organizations for accountant use.
//
//encore:api auth method=GET path=/organic-organizations tag:accountant
func (s *Service) ListOrganicOrganizations(ctx context.Context) (*organization_handlers.ListOrganicOrganizationResponse, error) {
	handler := organization_handlers.NewOrganizationService(s.query)
	return handler.ListOrganicOrganizations(ctx)
}
