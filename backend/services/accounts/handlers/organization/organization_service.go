package organization

import (
	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
)

var (
	ErrNameAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Name already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeNameAlreadyExists},
	)

	ErrOrganizationOrganicRequiresIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Organic organization must have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeOrganicOrgRequiresIcountClientID, Field: "icountClientId"},
	)

	ErrOrganizationNonOrganicForbidsIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Non-organic organization must not have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeNonOrganicOrgForbidsIcountClientID, Field: "icountClientId"},
	)
)

type OrganizationService struct {
	query db.Querier
}

func NewOrganizationService(query db.Querier) *OrganizationService {
	return &OrganizationService{query: query}
}

type OrganizationResponse struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
	Obligo         *int32  `json:"obligo" encore:"optional"`
}

func toOrganizationResponse(o db.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:             o.ID,
		Name:           o.Name,
		IsOrganic:      o.IsOrganic,
		IcountClientID: o.IcountClientID,
		Phone:          o.Phone,
		Address:        o.Address,
		Obligo:         o.Obligo,
	}
}

// validateIcountClientIDConstraint enforces billing rules:
// organic orgs must have an icount_client_id; non-organic orgs must not.
func validateIcountClientIDConstraint(isOrganic bool, icountClientID *int32) error {
	if isOrganic && icountClientID == nil {
		return ErrOrganizationOrganicRequiresIcountClientID
	}
	if !isOrganic && icountClientID != nil {
		return ErrOrganizationNonOrganicForbidsIcountClientID
	}
	return nil
}
