package reservation

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

func TestListOpenReservationsByBillingEntity(t *testing.T) {
	t.Run("validation rejects missing both ids", func(t *testing.T) {
		t.Parallel()
		req := &ListOpenReservationsByBillingEntityParams{}
		api_errors.AssertApiError(t, ErrInvalidBillingEntity, req.Validate())
	})

	t.Run("validation rejects both ids provided", func(t *testing.T) {
		t.Parallel()
		req := &ListOpenReservationsByBillingEntityParams{OfficeID: 1, OrgID: 2}
		api_errors.AssertApiError(t, ErrInvalidBillingEntity, req.Validate())
	})

	t.Run("returns error when GetAgentsByOfficeID fails", func(t *testing.T) {
		t.Parallel()
		s := &Service{query: testQuerier()}
		et.MockEndpoint(accounts.GetAgentsByOfficeID, func(_ context.Context, _ accounts.GetAgentsByOfficeIDRequest) (*accounts.GetAgentsResponse, error) {
			return nil, api_errors.ErrNotFound
		})

		_, err := s.ListOpenReservationsByBillingEntity(context.Background(), &ListOpenReservationsByBillingEntityParams{OfficeID: 99})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when GetAgentsByOrganizationID fails", func(t *testing.T) {
		t.Parallel()
		s := &Service{query: testQuerier()}
		et.MockEndpoint(accounts.GetAgentsByOrganizationID, func(_ context.Context, _ accounts.GetAgentsByOrganizationIDRequest) (*accounts.GetAgentsResponse, error) {
			return nil, api_errors.ErrNotFound
		})

		_, err := s.ListOpenReservationsByBillingEntity(context.Background(), &ListOpenReservationsByBillingEntityParams{OrgID: 99})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns ErrOfficeInOrganicOrg when office belongs to organic org", func(t *testing.T) {
		t.Parallel()
		s := &Service{query: testQuerier()}
		et.MockEndpoint(accounts.GetAgentsByOfficeID, func(_ context.Context, _ accounts.GetAgentsByOfficeIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{1, 2}, IsOrganic: true}, nil
		})

		_, err := s.ListOpenReservationsByBillingEntity(context.Background(), &ListOpenReservationsByBillingEntityParams{OfficeID: 10})
		api_errors.AssertApiError(t, ErrOfficeInOrganicOrg, err)
	})

	t.Run("returns ErrOrgIsInorganic when org is inorganic", func(t *testing.T) {
		t.Parallel()
		s := &Service{query: testQuerier()}
		et.MockEndpoint(accounts.GetAgentsByOrganizationID, func(_ context.Context, _ accounts.GetAgentsByOrganizationIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{1, 2}, IsOrganic: false}, nil
		})

		_, err := s.ListOpenReservationsByBillingEntity(context.Background(), &ListOpenReservationsByBillingEntityParams{OrgID: 5})
		api_errors.AssertApiError(t, ErrOrgIsInorganic, err)
	})

	t.Run("returns reservations for multiple agents in an office (inorganic)", func(t *testing.T) {
		t.Parallel()
		const agent1, agent2 int64 = 10001, 10002
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		vn1 := "VCH-OFFICE-1"
		id1 := seedReservation(t, ctx, s, agent1, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-OFFICE-1"
			p.CurrencyCode = "USD"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id1, UserID: agent1, VoucherNumber: &vn1}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		vn2 := "VCH-OFFICE-2"
		id2 := seedReservation(t, ctx, s, agent2, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-OFFICE-2"
			p.CurrencyCode = "EUR"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id2, UserID: agent2, VoucherNumber: &vn2}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		et.MockEndpoint(accounts.GetAgentsByOfficeID, func(_ context.Context, _ accounts.GetAgentsByOfficeIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{agent1, agent2}, IsOrganic: false}, nil
		})

		resp, err := s.ListOpenReservationsByBillingEntity(ctx, &ListOpenReservationsByBillingEntityParams{OfficeID: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		assertCurrencyGroups(t, resp.CurrencyGroups, map[string]int{"USD": 1, "EUR": 1})
	})

	t.Run("returns reservations for multiple agents in an organic org", func(t *testing.T) {
		t.Parallel()
		const agent1, agent2 int64 = 20001, 20002
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		// Two reservations in USD, one in EUR — exercises multi-currency grouping.
		vn1 := "VCH-ORG-1"
		id1 := seedReservation(t, ctx, s, agent1, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-ORG-1"
			p.CurrencyCode = "USD"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id1, UserID: agent1, VoucherNumber: &vn1}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		vn2 := "VCH-ORG-2"
		id2 := seedReservation(t, ctx, s, agent2, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-ORG-2"
			p.CurrencyCode = "USD"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id2, UserID: agent2, VoucherNumber: &vn2}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		vn3 := "VCH-ORG-3"
		id3 := seedReservation(t, ctx, s, agent2, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-ORG-3"
			p.CurrencyCode = "EUR"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id3, UserID: agent2, VoucherNumber: &vn3}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		et.MockEndpoint(accounts.GetAgentsByOrganizationID, func(_ context.Context, _ accounts.GetAgentsByOrganizationIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{agent1, agent2}, IsOrganic: true}, nil
		})

		resp, err := s.ListOpenReservationsByBillingEntity(ctx, &ListOpenReservationsByBillingEntityParams{OrgID: 5})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		assertCurrencyGroups(t, resp.CurrencyGroups, map[string]int{"USD": 2, "EUR": 1})
	})

	t.Run("returns internal error when db query fails", func(t *testing.T) {
		t.Parallel()
		q, s := mockService(t)
		et.MockEndpoint(accounts.GetAgentsByOfficeID, func(_ context.Context, _ accounts.GetAgentsByOfficeIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{1}, IsOrganic: false}, nil
		})
		q.EXPECT().GetPaymentPendingReservationsByAgentsIDs(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListOpenReservationsByBillingEntity(context.Background(), &ListOpenReservationsByBillingEntityParams{OfficeID: 10})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns empty list when no open reservations exist for the agents", func(t *testing.T) {
		t.Parallel()
		const agent1, agent2 int64 = 30001, 30002
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		// Seed only booked reservations for these agents — they must not appear in open reservations.
		seedReservation(t, ctx, s, agent1, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-EMPTY-1"
		})
		seedReservation(t, ctx, s, agent2, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BILLING-EMPTY-2"
		})

		et.MockEndpoint(accounts.GetAgentsByOfficeID, func(_ context.Context, _ accounts.GetAgentsByOfficeIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{agent1, agent2}, IsOrganic: false}, nil
		})

		resp, err := s.ListOpenReservationsByBillingEntity(ctx, &ListOpenReservationsByBillingEntityParams{OfficeID: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.CurrencyGroups) != 0 {
			t.Fatalf("expected 0 currency groups, got %d", len(resp.CurrencyGroups))
		}
	})

	t.Run("maps price details and row data correctly", func(t *testing.T) {
		t.Parallel()
		// PurchasePrice=100, BrokerErpPrice=15 → carPurchasePrice=115
		// MarkupPercentage=45% → carSellingPrice=115*1.45=166.75, profit=51.75
		// BtErpPrice=20 → erpSellingPrice=20
		const agentID int64 = 40001
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		vn := "VCH-PRICE"
		id := seedReservation(t, ctx, s, agentID, func(p *CreateReservationParams) {
			p.BrokerReservationID = "BRK-PRICE"
			p.PurchasePrice = 100.00
			p.MarkupPercentage = 45.00
			p.BrokerErpPrice = 15.00
			p.BtErpPrice = 20
			p.CurrencyCode = "USD"
			p.PickupDate = "2026-07-01"
			p.DropoffDate = "2026-07-05"
		})
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id, UserID: agentID, VoucherNumber: &vn}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}

		et.MockEndpoint(accounts.GetAgentsByOrganizationID, func(_ context.Context, _ accounts.GetAgentsByOrganizationIDRequest) (*accounts.GetAgentsResponse, error) {
			return &accounts.GetAgentsResponse{IDs: []int64{agentID}, IsOrganic: true}, nil
		})

		resp, err := s.ListOpenReservationsByBillingEntity(ctx, &ListOpenReservationsByBillingEntityParams{OrgID: 5})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.CurrencyGroups) != 1 {
			t.Fatalf("expected 1 currency group, got %d", len(resp.CurrencyGroups))
		}
		if resp.CurrencyGroups[0].CurrencyCode != "USD" {
			t.Fatalf("expected currency group USD, got %s", resp.CurrencyGroups[0].CurrencyCode)
		}
		if len(resp.CurrencyGroups[0].Reservations) != 1 {
			t.Fatalf("expected 1 reservation in group, got %d", len(resp.CurrencyGroups[0].Reservations))
		}

		r := resp.CurrencyGroups[0].Reservations[0]
		if r.BrokerReservationID != "BRK-PRICE" {
			t.Errorf("expected BrokerReservationID=BRK-PRICE, got %s", r.BrokerReservationID)
		}
		if r.PaymentStatus != string(db.PaymentStatusUnpaid) {
			t.Errorf("expected payment_status=unpaid, got %s", r.PaymentStatus)
		}
		if r.ReservationStatus != string(db.ReservationStatusVouchered) {
			t.Errorf("expected reservation_status=vouchered, got %s", r.ReservationStatus)
		}
		if r.CurrencyCode != "USD" {
			t.Errorf("expected currency_code=USD, got %s", r.CurrencyCode)
		}
		if r.PickupDate != "2026-07-01" {
			t.Errorf("expected pickup_date=2026-07-01, got %s", r.PickupDate)
		}

		if r.CarPurchasePrice != roundPrice(115.0) { // 100 + 15
			t.Errorf("CarPurchasePrice: want 115.00, got %.2f", r.CarPurchasePrice)
		}
		if r.CarSellingPrice != roundPrice(166.75) { // 115 * 1.45
			t.Errorf("CarSellingPrice: want 166.75, got %.2f", r.CarSellingPrice)
		}
		if r.ProfitOnCar != roundPrice(51.75) { // 166.75 - 115
			t.Errorf("ProfitOnCar: want 51.75, got %.2f", r.ProfitOnCar)
		}
		if r.ERPSellingPrice != 20.0 {
			t.Errorf("ERPSellingPrice: want 20.00, got %.2f", r.ERPSellingPrice)
		}
	})
}

// assertCurrencyGroups verifies that the response contains exactly the expected
// currency groups, each with the expected number of reservations.
func assertCurrencyGroups(t *testing.T, groups []CurrencyGroup, want map[string]int) {
	t.Helper()
	if len(groups) != len(want) {
		t.Fatalf("expected %d currency groups, got %d (%v)", len(want), len(groups), groups)
	}
	got := make(map[string]int, len(groups))
	for _, g := range groups {
		if _, dup := got[g.CurrencyCode]; dup {
			t.Fatalf("duplicate currency group for %s", g.CurrencyCode)
		}
		got[g.CurrencyCode] = len(g.Reservations)
	}
	for code, n := range want {
		if got[code] != n {
			t.Errorf("currency %s: expected %d reservations, got %d", code, n, got[code])
		}
	}
}
