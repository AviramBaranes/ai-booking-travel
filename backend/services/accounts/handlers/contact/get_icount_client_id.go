package contact

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

type GetIcountClientIDParams struct {
	OfficeID       *int64 `encore:"optional"`
	OrganizationID *int64 `encore:"optional"`
}

type GetIcountClientIDResponse struct {
	ClientID int32
}

func (s *ContactService) GetIcountClientID(ctx context.Context, p GetIcountClientIDParams) (*GetIcountClientIDResponse, error) {
	if p.OfficeID != nil {
		icountClientID, err := s.query.GetOfficeIcountClientID(ctx, *p.OfficeID)
		if err != nil {
			rlog.Error("failed to get iCount client ID for office", "error", err, "office_id", p.OfficeID)
			return nil, api_errors.ErrInternalError
		}
		if icountClientID == nil {
			return nil, api_errors.ErrNotFound
		}
		return &GetIcountClientIDResponse{ClientID: *icountClientID}, nil
	}

	if p.OrganizationID != nil {
		icountClientID, err := s.query.GetOrganizationIcountClientID(ctx, *p.OrganizationID)
		if err != nil {
			rlog.Error("failed to get iCount client ID for organization", "error", err, "organization_id", p.OrganizationID)
			return nil, api_errors.ErrInternalError
		}
		if icountClientID == nil {
			return nil, api_errors.ErrNotFound
		}
		return &GetIcountClientIDResponse{ClientID: *icountClientID}, nil
	}

	return nil, api_errors.ErrInvalidValue
}
