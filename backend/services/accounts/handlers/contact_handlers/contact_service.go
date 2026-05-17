package contact_handlers

import "encore.app/services/accounts/db"

type ContactService struct {
	query db.Querier
}

func NewContactService(query db.Querier) *ContactService {
	return &ContactService{query: query}
}

// ContactResponse is the shared response type for contact operations.
type ContactResponse struct {
	ID                   int64   `json:"id"`
	FirstName            string  `json:"firstName"`
	LastName             string  `json:"lastName"`
	Role                 string  `json:"role"`
	Cellphone            string  `json:"cellphone"`
	Email                string  `json:"email"`
	OfficeID             *int64  `json:"officeId" encore:"optional"`
	OrganizationID       *int64  `json:"organizationId" encore:"optional"`
	IsPaymentResponsible bool    `json:"isPaymentResponsible"`
	OfficeName           *string `json:"officeName" encore:"optional"`
	OrganizationName     *string `json:"organizationName" encore:"optional"`
}

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
