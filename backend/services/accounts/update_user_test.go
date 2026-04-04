package accounts

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid params", func(t *testing.T) {
		p := UpdateUserParams{}
		err := p.Validate()
		expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code:  api_errors.CodeInvalidValue,
			Field: "id",
		})
		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("User not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)

		params := UpdateUserParams{
			ID: 999,
		}

		// Mock DB to return ErrNoRows
		q.EXPECT().
			UpdateUser(gomock.Any(), db.UpdateUserParams{
				ID:          params.ID,
				PhoneNumber: params.PhoneNumber,
			}).
			Return(db.UpdateUserRow{}, db.ErrNoRows)

		s := &Service{query: q}
		_, err := s.UpdateUser(ctx, params)
		api_errors.AssertApiError(t, ErrUserNotFound, err)
	})

	t.Run("Database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)

		params := UpdateUserParams{
			ID: 123,
		}

		q.EXPECT().
			UpdateUser(gomock.Any(), gomock.Any()).
			Return(db.UpdateUserRow{}, errors.New("db error"))

		s := &Service{query: q}
		_, err := s.UpdateUser(ctx, params)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful update", func(t *testing.T) {
		agentEmail := "agent@example.com"
		phoneNumber := "0505050502"
		agent, delAgent, err := registerAgent(ctx, RegisterAgentParams{
			Email:       agentEmail,
			PhoneNumber: phoneNumber,
			Password:    testPassword,
		})
		if err != nil {
			t.Fatalf("Failed to create test agent: %v", err)
		}
		defer delAgent()

		// Update user with new values
		newPhone := "0555555555"
		org, err := query.CreateOrganization(ctx, db.CreateOrganizationParams{
			Name:      randomName(),
			IsOrganic: false,
		})
		if err != nil {
			t.Fatalf("Failed to create organization: %v", err)
		}

		office, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name:           randomName(),
			OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("Failed to create office: %v", err)
		}
		newOffice := office.ID

		params := UpdateUserParams{
			ID:          agent.ID,
			PhoneNumber: &newPhone,
			OfficeID:    &newOffice,
		}

		resp, err := UpdateUser(ctx, params)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify response fields
		if resp.ID != agent.ID {
			t.Errorf("Expected ID %d, got %d", agent.ID, resp.ID)
		}
		if resp.PhoneNumber != newPhone {
			t.Errorf("Expected PhoneNumber %s, got %s", newPhone, resp.PhoneNumber)
		}
		if *resp.OfficeID != newOffice {
			t.Errorf("Expected OfficeID %d, got %d", newOffice, *resp.OfficeID)
		}

		// Verify changes persisted in DB
		updatedUser, err := query.GetUserById(ctx, agent.ID)
		if err != nil {
			t.Fatalf("Failed to get user from DB: %v", err)
		}

		if *updatedUser.PhoneNumber != newPhone {
			t.Errorf("DB: Expected PhoneNumber %s, got %s", newPhone, *updatedUser.PhoneNumber)
		}
	})
}
