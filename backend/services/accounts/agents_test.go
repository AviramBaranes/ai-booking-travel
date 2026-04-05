package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func agentMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

// seedAgent creates an org, office, and agent for use in agent tests.
// Returns the agent response, orgID, officeID.
func seedAgent(t *testing.T, s *Service, email, phone string, officeID int32) *CreateAgentResponse {
	t.Helper()
	resp, err := s.CreateAgent(context.Background(), CreateAgentRequest{
		Email:       email,
		Password:    "ValidPass123!",
		PhoneNumber: phone,
		OfficeID:    officeID,
	})
	if err != nil {
		t.Fatalf("failed to create agent %s: %v", email, err)
	}
	return resp
}

// --- Tests ---

func TestListAgents(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("pagination returns max 15 per page", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		prefix := fmt.Sprintf("PagAgent%d", time.Now().UnixNano())

		for i := 1; i <= 18; i++ {
			seedAgent(t, s,
				fmt.Sprintf("%s_%02d@test.com", prefix, i),
				fmt.Sprintf("05%08d", i),
				officeID,
			)
		}

		page1, err := s.ListAgents(ctx, &ListAgentsRequest{Search: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page1.Agents) != 15 {
			t.Fatalf("expected 15 agents on page 1, got %d", len(page1.Agents))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := s.ListAgents(ctx, &ListAgentsRequest{Search: prefix, Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(page2.Agents) != 3 {
			t.Fatalf("expected 3 agents on page 2, got %d", len(page2.Agents))
		}
		if page2.Total != 18 {
			t.Fatalf("expected total 18, got %d", page2.Total)
		}

		page1IDs := make(map[int32]bool)
		for _, a := range page1.Agents {
			page1IDs[a.ID] = true
		}
		for _, a := range page2.Agents {
			if page1IDs[a.ID] {
				t.Fatalf("agent %d appeared on both pages", a.ID)
			}
		}
	})

	t.Run("empty page returns no results with same total", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		prefix := fmt.Sprintf("EmptyPgAgent%d", time.Now().UnixNano())

		for i := 1; i <= 18; i++ {
			seedAgent(t, s,
				fmt.Sprintf("%s_%02d@test.com", prefix, i),
				fmt.Sprintf("06%08d", i),
				officeID,
			)
		}

		resp, err := s.ListAgents(ctx, &ListAgentsRequest{Search: prefix, Page: 3})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Agents) != 0 {
			t.Fatalf("expected 0 agents on page 3, got %d", len(resp.Agents))
		}
		if resp.Total != 18 {
			t.Fatalf("expected total 18, got %d", resp.Total)
		}
	})

	t.Run("filters by search email substring", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		unique := fmt.Sprintf("UniqueAgent%d", time.Now().UnixNano())

		seedAgent(t, s,
			fmt.Sprintf("%s@test.com", unique),
			randomIsraeliPhoneNumber(),
			officeID,
		)

		resp, err := s.ListAgents(ctx, &ListAgentsRequest{Search: unique, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Agents) != 1 {
			t.Fatalf("expected 1 result, got %d", len(resp.Agents))
		}
		if resp.Agents[0].Email != fmt.Sprintf("%s@test.com", unique) {
			t.Fatalf("expected email %s@test.com, got %s", unique, resp.Agents[0].Email)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("filters by officeId", func(t *testing.T) {
		t.Parallel()
		orgID, officeA := seedOrgAndOffice(t)
		officeB, err := query.CreateOffice(ctx, db.CreateOfficeParams{
			Name:           fmt.Sprintf("AgentFilterOfficeB_%d", time.Now().UnixNano()),
			OrganizationID: orgID,
		})
		if err != nil {
			t.Fatalf("failed to create officeB: %v", err)
		}

		prefix := fmt.Sprintf("OffFilterAgent%d", time.Now().UnixNano())
		seedAgent(t, s, fmt.Sprintf("%s_a@test.com", prefix), randomIsraeliPhoneNumber(), officeA)
		seedAgent(t, s, fmt.Sprintf("%s_b@test.com", prefix), randomIsraeliPhoneNumber(), officeB.ID)

		resp, err := s.ListAgents(ctx, &ListAgentsRequest{Search: prefix, OfficeID: officeA, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Agents) != 1 {
			t.Fatalf("expected 1 agent for officeA, got %d", len(resp.Agents))
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	t.Run("filters by orgId", func(t *testing.T) {
		t.Parallel()
		orgA, officeA := seedOrgAndOffice(t)
		_, officeB := seedOrgAndOffice(t) // different org

		prefix := fmt.Sprintf("OrgFilterAgent%d", time.Now().UnixNano())
		seedAgent(t, s, fmt.Sprintf("%s_a@test.com", prefix), randomIsraeliPhoneNumber(), officeA)
		seedAgent(t, s, fmt.Sprintf("%s_b@test.com", prefix), randomIsraeliPhoneNumber(), officeB)

		resp, err := s.ListAgents(ctx, &ListAgentsRequest{Search: prefix, OrgID: orgA, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Agents) != 1 {
			t.Fatalf("expected 1 agent for orgA, got %d", len(resp.Agents))
		}
		if resp.Agents[0].Email != fmt.Sprintf("%s_a@test.com", prefix) {
			t.Fatalf("expected agent A, got %s", resp.Agents[0].Email)
		}
	})

	t.Run("validation rejects page 0", func(t *testing.T) {
		t.Parallel()
		p := ListAgentsRequest{Page: 0}
		api_errors.AssertApiError(t, invalidValueErr("page"), p.Validate())
	})

	t.Run("returns error when list db fails", func(t *testing.T) {
		t.Parallel()
		q, s := agentMockService(t)
		q.EXPECT().ListAgents(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListAgents(ctx, &ListAgentsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count db fails", func(t *testing.T) {
		t.Parallel()
		q, s := agentMockService(t)
		q.EXPECT().ListAgents(gomock.Any(), gomock.Any()).Return([]db.ListAgentsRow{}, nil)
		q.EXPECT().CountAgents(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.ListAgents(ctx, &ListAgentsRequest{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateAgent(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates agent successfully", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		resp, err := s.CreateAgent(ctx, CreateAgentRequest{
			Email:       fmt.Sprintf("create_agent_ok_%d@test.com", time.Now().UnixNano()),
			Password:    "ValidPass123!",
			PhoneNumber: randomIsraeliPhoneNumber(),
			OfficeID:    officeID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
	})

	t.Run("returns error on duplicate email", func(t *testing.T) {
		t.Parallel()
		_, officeID := seedOrgAndOffice(t)
		email := fmt.Sprintf("dup_agent_%d@test.com", time.Now().UnixNano())
		seedAgent(t, s, email, randomIsraeliPhoneNumber(), officeID)

		_, err := s.CreateAgent(ctx, CreateAgentRequest{
			Email:       email,
			Password:    "ValidPass123!",
			PhoneNumber: randomIsraeliPhoneNumber(),
			OfficeID:    officeID,
		})
		api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		p := CreateAgentRequest{Email: "not-an-email", Password: "ValidPass123!", PhoneNumber: "0521234567", OfficeID: 1}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects empty email", func(t *testing.T) {
		t.Parallel()
		p := CreateAgentRequest{Email: "", Password: "ValidPass123!", PhoneNumber: "0521234567", OfficeID: 1}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects empty phoneNumber", func(t *testing.T) {
		t.Parallel()
		p := CreateAgentRequest{Email: "agent@test.com", Password: "ValidPass123!", PhoneNumber: "", OfficeID: 1}
		api_errors.AssertApiError(t, invalidValueErr("phoneNumber"), p.Validate())
	})

	t.Run("validation rejects officeId 0", func(t *testing.T) {
		t.Parallel()
		p := CreateAgentRequest{Email: "agent@test.com", Password: "ValidPass123!", PhoneNumber: "0521234567", OfficeID: 0}
		api_errors.AssertApiError(t, invalidValueErr("officeId"), p.Validate())
	})

	t.Run("validation rejects weak password", func(t *testing.T) {
		t.Parallel()
		p := CreateAgentRequest{Email: "agent@test.com", Password: "short", PhoneNumber: "0521234567", OfficeID: 1}
		api_errors.AssertApiError(t, ErrPasswordTooShort, p.Validate())
	})

	t.Run("returns error when check exists db fails", func(t *testing.T) {
		t.Parallel()
		q, s := agentMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("db error"))

		_, err := s.CreateAgent(ctx, CreateAgentRequest{
			Email:       "db_fail@test.com",
			Password:    "ValidPass123!",
			PhoneNumber: "0521234567",
			OfficeID:    1,
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when create db fails", func(t *testing.T) {
		t.Parallel()
		q, s := agentMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), db.ErrNoRows)
		q.EXPECT().CreateAgent(gomock.Any(), gomock.Any()).Return(db.CreateAgentRow{}, errors.New("db error"))

		_, err := s.CreateAgent(ctx, CreateAgentRequest{
			Email:       "db_create_fail@test.com",
			Password:    "ValidPass123!",
			PhoneNumber: "0521234567",
			OfficeID:    1,
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
