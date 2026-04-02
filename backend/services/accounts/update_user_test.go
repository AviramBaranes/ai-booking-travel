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

		agentCode := "agent123"
		params := UpdateUserParams{
			ID:        999,
			AgentCode: &agentCode,
		}

		// Mock DB to return ErrNoRows
		q.EXPECT().
			UpdateUser(gomock.Any(), db.UpdateUserParams{
				ID:          params.ID,
				AgentCode:   db.TextParam(params.AgentCode),
				PhoneNumber: db.TextParam(params.PhoneNumber),
				OfficeCode:  db.TextParam(params.OfficeCode),
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

		agentCode := "agent123"
		params := UpdateUserParams{
			ID:        123,
			AgentCode: &agentCode,
		}

		q.EXPECT().
			UpdateUser(gomock.Any(), gomock.Any()).
			Return(db.UpdateUserRow{}, errors.New("db error"))

		s := &Service{query: q}
		_, err := s.UpdateUser(ctx, params)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("Successful update", func(t *testing.T) {
		// Create test user
		adminUsername := "admin_update_test"
		admin, delAdmin, err := registerAdmin(ctx, adminUsername, testPassword)
		if err != nil {
			t.Fatalf("Failed to create test admin: %v", err)
		}
		defer delAdmin()

		// Update user with new values
		newPhone := "123-456-7890"
		newOffice := "Office A"
		newAgentCode := "Agent 007"

		params := UpdateUserParams{
			ID:          admin.ID,
			PhoneNumber: &newPhone,
			OfficeCode:  &newOffice,
			AgentCode:   &newAgentCode,
		}

		resp, err := UpdateUser(ctx, params)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify response fields
		if resp.ID != admin.ID {
			t.Errorf("Expected ID %d, got %d", admin.ID, resp.ID)
		}
		if resp.PhoneNumber != newPhone {
			t.Errorf("Expected PhoneNumber %s, got %s", newPhone, resp.PhoneNumber)
		}
		if resp.OfficeCode != newOffice {
			t.Errorf("Expected OfficeCode %s, got %s", newOffice, resp.OfficeCode)
		}
		if resp.AgentCode != newAgentCode {
			t.Errorf("Expected AgentCode %s, got %s", newAgentCode, resp.AgentCode)
		}

		// Verify changes persisted in DB
		updatedUser, err := query.GetUserById(ctx, admin.ID)
		if err != nil {
			t.Fatalf("Failed to get user from DB: %v", err)
		}

		if db.StringFromTextParam(updatedUser.PhoneNumber) != newPhone {
			t.Errorf("DB: Expected PhoneNumber %s, got %s", newPhone, db.StringFromTextParam(updatedUser.PhoneNumber))
		}
		if db.StringFromTextParam(updatedUser.OfficeCode) != newOffice {
			t.Errorf("DB: Expected OfficeCode %s, got %s", newOffice, db.StringFromTextParam(updatedUser.OfficeCode))
		}
		if db.StringFromTextParam(updatedUser.AgentCode) != newAgentCode {
			t.Errorf("DB: Expected AgentCode %s, got %s", newAgentCode, db.StringFromTextParam(updatedUser.AgentCode))
		}
	})
}
