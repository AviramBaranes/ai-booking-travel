package reservation

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	queries "encore.app/services/reservation/handlers/queries"
)

// --- Helpers ---

// seedOfficeReservation creates a reservation billed to an office of an inorganic organization.
func seedOfficeReservation(t *testing.T, ctx context.Context, s *Service, agentID, officeID int64, brokerReservationID string) int64 {
	t.Helper()
	orgID := int64(20)
	isOrganic := false
	return seedReservation(t, ctx, s, agentID, func(p *CreateReservationParams) {
		p.BrokerReservationID = brokerReservationID
		p.CurrencyCode = "USD"
		p.OfficeID = &officeID
		p.OrganizationID = &orgID
		p.IsOrganizationOrganic = &isOrganic
	})
}

func voucher(t *testing.T, ctx context.Context, s *Service, id, agentID int64) {
	t.Helper()
	vn := "VCH-CANCELED"
	if err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: id, UserID: agentID, VoucherNumber: &vn, CurrencyRate: dbadapters.NumericFromFloat64(1)}); err != nil {
		t.Fatalf("failed to apply voucher: %v", err)
	}
}

func cancel(t *testing.T, ctx context.Context, s *Service, id int64) {
	t.Helper()
	if err := s.query.CancelReservation(ctx, id); err != nil {
		t.Fatalf("failed to cancel reservation: %v", err)
	}
}

func resolve(t *testing.T, ctx context.Context, s *Service, ids ...int64) {
	t.Helper()
	if err := s.query.ResolveReservationsPayment(ctx, ids); err != nil {
		t.Fatalf("failed to resolve reservations: %v", err)
	}
}

func assertStatuses(t *testing.T, ctx context.Context, s *Service, id int64, wantReservation db.ReservationStatus, wantPayment db.PaymentStatus) {
	t.Helper()
	r, err := s.query.GetReservationByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to get reservation: %v", err)
	}
	if r.ReservationStatus != wantReservation || r.PaymentStatus != wantPayment {
		t.Fatalf("expected %s/%s, got %s/%s", wantReservation, wantPayment, r.ReservationStatus, r.PaymentStatus)
	}
}

// --- Tests ---

func TestFindCanceledUnpaidReservation(t *testing.T) {
	t.Run("validation rejects missing supplier reservation id", func(t *testing.T) {
		t.Parallel()
		req := &FindCanceledUnpaidReservationParams{OfficeID: 1}
		api_errors.AssertApiError(t, invalidValueErr("broker_reservation_id"), req.Validate())
	})

	t.Run("validation rejects missing both ids", func(t *testing.T) {
		t.Parallel()
		req := &FindCanceledUnpaidReservationParams{BrokerReservationID: "X"}
		api_errors.AssertApiError(t, ErrInvalidBillingEntity, req.Validate())
	})

	t.Run("validation rejects both ids provided", func(t *testing.T) {
		t.Parallel()
		req := &FindCanceledUnpaidReservationParams{OfficeID: 1, OrgID: 2, BrokerReservationID: "X"}
		api_errors.AssertApiError(t, ErrInvalidBillingEntity, req.Validate())
	})

	t.Run("finds a canceled unpaid reservation of the billing entity", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID int64 = 20001, 2001
		ctx := context.Background()
		s := &Service{query: testQuerier()}
		qs := queries.NewQueryService(s.query)

		id := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-FOUND-1")
		voucher(t, ctx, s, id, agentID)
		cancel(t, ctx, s, id)

		resp, err := qs.FindCanceledUnpaidReservation(ctx, &FindCanceledUnpaidReservationParams{
			OfficeID:            officeID,
			BrokerReservationID: "  CANCELED-FOUND-1 ",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		r := resp.Reservation
		if r.ID != id || r.PaymentStatus != PaymentStatusUnpaid || r.ReservationStatus != ReservationStatusCanceled {
			t.Fatalf("unexpected reservation: %+v", r)
		}
		if r.CurrencyCode != "USD" || r.TotalPrice <= 0 {
			t.Fatalf("expected a priced USD reservation, got %+v", r)
		}
	})

	t.Run("finds a reservation canceled before it was vouchered", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID int64 = 20002, 2002
		ctx := context.Background()
		s := &Service{query: testQuerier()}
		qs := queries.NewQueryService(s.query)

		id := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-BOOKED-1")
		cancel(t, ctx, s, id)

		resp, err := qs.FindCanceledUnpaidReservation(ctx, &FindCanceledUnpaidReservationParams{
			OfficeID:            officeID,
			BrokerReservationID: "CANCELED-BOOKED-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Reservation.ID != id {
			t.Fatalf("expected reservation %d, got %d", id, resp.Reservation.ID)
		}
	})

	t.Run("does not find a reservation of another billing entity", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID, otherOfficeID int64 = 20003, 2003, 2004
		ctx := context.Background()
		s := &Service{query: testQuerier()}
		qs := queries.NewQueryService(s.query)

		id := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-OTHER-1")
		cancel(t, ctx, s, id)

		_, err := qs.FindCanceledUnpaidReservation(ctx, &FindCanceledUnpaidReservationParams{
			OfficeID:            otherOfficeID,
			BrokerReservationID: "CANCELED-OTHER-1",
		})
		api_errors.AssertApiError(t, ErrCanceledReservationNotFound, err)
	})

	t.Run("does not find an open or a refund-pending reservation", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID int64 = 20004, 2005
		ctx := context.Background()
		s := &Service{query: testQuerier()}
		qs := queries.NewQueryService(s.query)

		openID := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-OPEN-1")
		voucher(t, ctx, s, openID, agentID)

		refundID := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-REFUND-1")
		voucher(t, ctx, s, refundID, agentID)
		resolve(t, ctx, s, refundID)
		cancel(t, ctx, s, refundID)

		for _, brokerID := range []string{"CANCELED-OPEN-1", "CANCELED-REFUND-1"} {
			_, err := qs.FindCanceledUnpaidReservation(ctx, &FindCanceledUnpaidReservationParams{
				OfficeID:            officeID,
				BrokerReservationID: brokerID,
			})
			api_errors.AssertApiError(t, ErrCanceledReservationNotFound, err)
		}
	})
}

func TestListCanceledUnpaidReservationsByIDs(t *testing.T) {
	t.Run("returns only the canceled unpaid reservations of the billing entity", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID, otherOfficeID int64 = 20005, 2006, 2007
		ctx := context.Background()
		s := &Service{query: testQuerier()}
		qs := queries.NewQueryService(s.query)

		canceledID := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-IDS-1")
		cancel(t, ctx, s, canceledID)

		openID := seedOfficeReservation(t, ctx, s, agentID, officeID, "CANCELED-IDS-2")
		voucher(t, ctx, s, openID, agentID)

		otherID := seedOfficeReservation(t, ctx, s, agentID, otherOfficeID, "CANCELED-IDS-3")
		cancel(t, ctx, s, otherID)

		resp, err := qs.ListCanceledUnpaidReservationsByIDs(ctx, &ListCanceledUnpaidReservationsByIDsParams{
			OfficeID: officeID,
			IDs:      []int64{canceledID, openID, otherID},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Reservations) != 1 || resp.Reservations[0].ID != canceledID {
			t.Fatalf("expected only reservation %d, got %+v", canceledID, resp.Reservations)
		}
	})
}

func TestResolveReservationsPayment(t *testing.T) {
	t.Run("a billed canceled unpaid reservation stays canceled and is owed back", func(t *testing.T) {
		t.Parallel()
		const agentID, officeID int64 = 20006, 2008
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		canceledID := seedOfficeReservation(t, ctx, s, agentID, officeID, "RESOLVE-CANCELED-1")
		voucher(t, ctx, s, canceledID, agentID)
		cancel(t, ctx, s, canceledID)

		openID := seedOfficeReservation(t, ctx, s, agentID, officeID, "RESOLVE-OPEN-1")
		voucher(t, ctx, s, openID, agentID)

		resolve(t, ctx, s, canceledID, openID)

		assertStatuses(t, ctx, s, canceledID, db.ReservationStatusCanceled, db.PaymentStatusRefundPending)
		assertStatuses(t, ctx, s, openID, db.ReservationStatusVouchered, db.PaymentStatusPaid)

		// The refund is settled on a later bill.
		resolve(t, ctx, s, canceledID)
		assertStatuses(t, ctx, s, canceledID, db.ReservationStatusCanceled, db.PaymentStatusRefunded)
	})
}
