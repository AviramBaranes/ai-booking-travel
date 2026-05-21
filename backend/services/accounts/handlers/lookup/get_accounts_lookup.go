package lookup

import (
	"context"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetAccountsLookup(ctx context.Context, p GetAccountsLookupParams) (*GetAccountsLookupResponse, error) {
	var organizations []AccountName
	var offices []AccountName
	var users []AccountName

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		rows, err := s.query.GetOrganizationNamesByIDs(groupCtx, p.OrganizationIDs)
		if err != nil {
			return fmt.Errorf("get organization names: %w", err)
		}
		organizations = make([]AccountName, 0, len(rows))
		for _, row := range rows {
			organizations = append(organizations, AccountName{ID: row.ID, Name: row.Name})
		}
		return nil
	})

	group.Go(func() error {
		rows, err := s.query.GetOfficeNamesByIDs(groupCtx, p.OfficeIDs)
		if err != nil {
			return fmt.Errorf("get office names: %w", err)
		}
		offices = make([]AccountName, 0, len(rows))
		for _, row := range rows {
			offices = append(offices, AccountName{ID: row.ID, Name: row.Name})
		}
		return nil
	})

	group.Go(func() error {
		rows, err := s.query.GetUserNamesByIDs(groupCtx, p.UserIDs)
		if err != nil {
			return fmt.Errorf("get user names: %w", err)
		}
		users = make([]AccountName, 0, len(rows))
		for _, row := range rows {
			users = append(users, AccountName{ID: row.ID, Name: row.FirstName + " " + row.LastName})
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		rlog.Error("failed to get accounts lookup", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetAccountsLookupResponse{
		Organizations: organizations,
		Offices:       offices,
		Users:         users,
	}, nil
}
