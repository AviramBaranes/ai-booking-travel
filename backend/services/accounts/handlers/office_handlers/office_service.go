package office_handlers

import (
	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.app/services/accounts/db"
)

var (
	ErrNameAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists, "Name already exists",
		api_errors.ErrorDetails{Code: api_errors.CodeNameAlreadyExists},
	)

	ErrOfficeOrganicForbidsIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Office under an organic organization must not have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeOfficeOrganicForbidsIcountClientID, Field: "icountClientId"},
	)

	ErrOfficeNonOrganicRequiresIcountClientID = api_errors.NewErrorWithDetail(
		errs.InvalidArgument, "Office under a non-organic organization must have an icount_client_id",
		api_errors.ErrorDetails{Code: api_errors.CodeOfficeNonOrganicRequiresIcountClientID, Field: "icountClientId"},
	)
)

type OfficeService struct {
	query db.Querier
}

func NewOfficeService(query db.Querier) *OfficeService {
	return &OfficeService{query: query}
}

type OfficeResponse struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	OrganizationID   int64   `json:"organizationId"`
	OrganizationName string  `json:"organizationName"`
	IcountClientID   *int32  `json:"icountClientId" encore:"optional"`
	Phone            *string `json:"phone" encore:"optional"`
	Address          *string `json:"address" encore:"optional"`
	ContactCount     int64   `json:"contactCount"`
	AgentCount       int64   `json:"agentCount"`
}

func toOfficeResponse(o db.ListOfficesRow) OfficeResponse {
	return OfficeResponse{
		ID:               o.ID,
		Name:             o.Name,
		OrganizationID:   o.OrganizationID,
		OrganizationName: o.OrganizationName,
		IcountClientID:   o.IcountClientID,
		Phone:            o.Phone,
		Address:          o.Address,
		ContactCount:     o.ContactCount,
		AgentCount:       o.AgentCount,
	}
}

// validateOfficeIcountClientIDConstraint enforces billing rules for offices.
// Offices under organic orgs must NOT have an icount_client_id;
// offices under non-organic orgs MUST have one.
func validateOfficeIcountClientIDConstraint(isOrganic bool, icountClientID *int32) error {
	if isOrganic && icountClientID != nil {
		return ErrOfficeOrganicForbidsIcountClientID
	}
	if !isOrganic && icountClientID == nil {
		return ErrOfficeNonOrganicRequiresIcountClientID
	}
	return nil
}
