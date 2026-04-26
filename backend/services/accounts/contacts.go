package accounts

import (
	"context"
	"errors"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// --- Request / Response types ---

type ContactResponse struct {
	ID                   int32   `json:"id"`
	FirstName            string  `json:"firstName"`
	LastName             string  `json:"lastName"`
	Role                 string  `json:"role"`
	Cellphone            string  `json:"cellphone"`
	Email                string  `json:"email"`
	OfficeID             *int32  `json:"officeId"`
	OrganizationID       *int32  `json:"organizationId"`
	IsPaymentResponsible bool    `json:"isPaymentResponsible"`
	OfficeName           *string `json:"officeName"`
	OrganizationName     *string `json:"organizationName"`
}

type ListContactsRequest struct {
	Search   string `query:"search"`
	OfficeID int32  `query:"officeId" validate:"omitempty,gte=1"`
	OrgID    int32  `query:"orgId" validate:"omitempty,gte=1"`
	Page     int32  `query:"page" validate:"required,gte=1"`
}

func (p ListContactsRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type ListContactsResponse struct {
	Contacts []ContactResponse `json:"contacts"`
	Total    int64             `json:"total"`
}

type CreateContactRequest struct {
	FirstName            string `json:"firstName" validate:"required,notblank"`
	LastName             string `json:"lastName" validate:"required,notblank"`
	Role                 string `json:"role" validate:"required,notblank"`
	Cellphone            string `json:"cellphone" validate:"required,notblank"`
	Email                string `json:"email" validate:"required,email"`
	OfficeID             *int32 `json:"officeId" encore:"optional"`
	OrganizationID       *int32 `json:"organizationId" encore:"optional"`
	IsPaymentResponsible bool   `json:"isPaymentResponsible" encore:"optional"`
}

func (p CreateContactRequest) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}

	hasOffice := p.OfficeID != nil
	hasOrg := p.OrganizationID != nil
	if hasOffice == hasOrg {
		return api_errors.NewErrorWithDetail(
			errs.InvalidArgument,
			"Exactly one of officeId or organizationId must be provided",
			api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue},
		)
	}

	return nil
}

type UpdateContactRequest struct {
	FirstName            *string `json:"firstName" validate:"omitempty,notblank" encore:"optional"`
	LastName             *string `json:"lastName" validate:"omitempty,notblank" encore:"optional"`
	Role                 *string `json:"role" validate:"omitempty,notblank" encore:"optional"`
	Cellphone            *string `json:"cellphone" validate:"omitempty,notblank" encore:"optional"`
	Email                *string `json:"email" validate:"omitempty,email" encore:"optional"`
	OfficeID             *int32  `json:"officeId" encore:"optional"`
	OrganizationID       *int32  `json:"organizationId" encore:"optional"`
	IsPaymentResponsible *bool   `json:"isPaymentResponsible" encore:"optional"`
}

func (p UpdateContactRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

const contactsPageSize int64 = 15

func toContactResponse(c db.Contact) ContactResponse {
	return ContactResponse{
		ID:                   c.ID,
		FirstName:            c.FirstName,
		LastName:             c.LastName,
		Role:                 c.Role,
		Cellphone:            c.Cellphone,
		Email:                c.Email,
		OfficeID:             c.OfficeID,
		OrganizationID:       c.OrganizationID,
		IsPaymentResponsible: c.IsPaymentResponsible,
	}
}

func toContactResponseFromRow(r db.ListContactsRow) ContactResponse {
	return ContactResponse{
		ID:                   r.ID,
		FirstName:            r.FirstName,
		LastName:             r.LastName,
		Role:                 r.Role,
		Cellphone:            r.Cellphone,
		Email:                r.Email,
		OfficeID:             r.OfficeID,
		OrganizationID:       r.OrganizationID,
		IsPaymentResponsible: r.IsPaymentResponsible,
		OfficeName:           r.OfficeName,
		OrganizationName:     r.OrganizationName,
	}
}

// --- Endpoints ---

// ListContacts lists contacts with optional filtering and pagination.
//
//encore:api auth method=GET path=/contacts tag:admin
func (s *Service) ListContacts(ctx context.Context, params *ListContactsRequest) (*ListContactsResponse, error) {
	offset := int64(params.Page-1) * contactsPageSize

	var searchPtr *string
	if s := strings.TrimSpace(params.Search); s != "" {
		searchPtr = &s
	}

	var officeIDPtr *int32
	if params.OfficeID > 0 {
		officeIDPtr = &params.OfficeID
	}

	var orgIDPtr *int32
	if params.OrgID > 0 {
		orgIDPtr = &params.OrgID
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

// CreateContact creates a new contact.
//
//encore:api auth method=POST path=/contacts tag:admin
func (s *Service) CreateContact(ctx context.Context, params CreateContactRequest) (*ContactResponse, error) {
	row, err := s.query.CreateContact(ctx, db.CreateContactParams{
		FirstName:            params.FirstName,
		LastName:             params.LastName,
		Role:                 params.Role,
		Cellphone:            params.Cellphone,
		Email:                params.Email,
		OfficeID:             params.OfficeID,
		OrganizationID:       params.OrganizationID,
		IsPaymentResponsible: params.IsPaymentResponsible,
	})
	if err != nil {
		rlog.Error("failed to create contact", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toContactResponse(row)
	return &resp, nil
}

// UpdateContact updates an existing contact.
//
//encore:api auth method=PUT path=/contacts/:id tag:admin
func (s *Service) UpdateContact(ctx context.Context, id int32, params UpdateContactRequest) (*ContactResponse, error) {
	row, err := s.query.UpdateContact(ctx, db.UpdateContactParams{
		ID:                   id,
		FirstName:            params.FirstName,
		LastName:             params.LastName,
		Role:                 params.Role,
		Cellphone:            params.Cellphone,
		Email:                params.Email,
		OfficeID:             params.OfficeID,
		OrganizationID:       params.OrganizationID,
		IsPaymentResponsible: params.IsPaymentResponsible,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update contact", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toContactResponse(row)
	return &resp, nil
}

// DeleteContact deletes a contact by its ID.
//
//encore:api auth method=DELETE path=/contacts/:id tag:admin
func (s *Service) DeleteContact(ctx context.Context, id int32) error {
	err := s.query.DeleteContact(ctx, id)
	if err != nil {
		rlog.Error("failed to delete contact", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
