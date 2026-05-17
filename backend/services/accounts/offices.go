package accounts

import (
	"context"

	"encore.app/services/accounts/handlers/office_handlers"
)

// --- Endpoints ---

// ListOffices lists offices with optional filtering and pagination.
//
//encore:api auth method=GET path=/offices tag:admin
func (s *Service) ListOffices(ctx context.Context, p office_handlers.ListOfficesParams) (*office_handlers.ListOfficesResponse, error) {
	handler := office_handlers.NewOfficeService(s.query)
	return handler.ListOffices(ctx, p)
}

// CreateOffice creates a new office.
//
//encore:api auth method=POST path=/offices tag:admin
func (s *Service) CreateOffice(ctx context.Context, p office_handlers.CreateOfficeParams) (*office_handlers.OfficeResponse, error) {
	handler := office_handlers.NewOfficeService(s.query)
	return handler.CreateOffice(ctx, p)
}

// UpdateOffice updates an existing office.
//
//encore:api auth method=PUT path=/offices/:id tag:admin
func (s *Service) UpdateOffice(ctx context.Context, id int64, p office_handlers.UpdateOfficeParams) (*office_handlers.OfficeResponse, error) {
	handler := office_handlers.NewOfficeService(s.query)
	return handler.UpdateOffice(ctx, id, p)
}

// ListInorganicOffices lists all inorganic offices for accountant use.
//
//encore:api auth method=GET path=/inorganic-offices tag:accountant
func (s *Service) ListInorganicOffices(ctx context.Context) (*office_handlers.ListInorganicOfficeResponse, error) {
	handler := office_handlers.NewOfficeService(s.query)
	return handler.ListInorganicOffices(ctx)
}
