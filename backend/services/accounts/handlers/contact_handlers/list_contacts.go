package contact_handlers

import (
	"context"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type ListContactsParams struct {
	Search   string `query:"search"`
	OfficeID int64  `query:"officeId" validate:"omitempty,gte=1"`
	OrgID    int64  `query:"orgId" validate:"omitempty,gte=1"`
	Page     int32  `query:"page" validate:"required,gte=1"`
}

func (p ListContactsParams) Validate() error {
	return validation.ValidateStruct(p)
}

type ListContactsResponse struct {
	Contacts []ContactResponse `json:"contacts"`
	Total    int64             `json:"total"`
}

func (s *ContactService) ListContacts(ctx context.Context, p ListContactsParams) (*ListContactsResponse, error) {
	offset := int64(p.Page-1) * contactsPageSize

	var searchPtr *string
	if trimmed := strings.TrimSpace(p.Search); trimmed != "" {
		searchPtr = &trimmed
	}

	var officeIDPtr *int64
	if p.OfficeID > 0 {
		officeIDPtr = &p.OfficeID
	}

	var orgIDPtr *int64
	if p.OrgID > 0 {
		orgIDPtr = &p.OrgID
	}

	rows, err := s.query.ListContacts(ctx, db.ListContactsParams{
		Name:           searchPtr,
		OfficeID:       officeIDPtr,
		OrganizationID: orgIDPtr,
		PageOffset:     offset,
		PageSize:       contactsPageSize,
	})
	if err != nil {
		rlog.Error("failed to list contacts", "error", err)
		return nil, api_errors.ErrInternalError
	}

	total, err := s.query.CountContacts(ctx, db.CountContactsParams{
		Name:           searchPtr,
		OfficeID:       officeIDPtr,
		OrganizationID: orgIDPtr,
	})
	if err != nil {
		rlog.Error("failed to count contacts", "error", err)
		return nil, api_errors.ErrInternalError
	}

	contacts := make([]ContactResponse, 0, len(rows))
	for _, r := range rows {
		contacts = append(contacts, toContactResponseFromRow(r))
	}

	return &ListContactsResponse{Contacts: contacts, Total: total}, nil
}
