package reservation

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

func TestGetOpenReservations(t *testing.T) {
	t.Run("returns internal error when db query fails", func(t *testing.T) {
		q, ms := mockService(t)
		q.EXPECT().GetPaymentPendingReservations(gomock.Any()).
			Return(nil, errors.New("db error"))
		et.MockService[Interface]("reservation", ms)

		_, err := GetOpenReservations(context.Background())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns empty list when no open reservations exist", func(t *testing.T) {
		const userID int32 = 8888
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		// Seed only non-open reservations (default status is "booked").
		seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
			p.BrokerReservationID = "EMPTY-BOOKED-1"
		})
		seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
			p.BrokerReservationID = "EMPTY-BOOKED-2"
		})

		resp, err := GetOpenReservations(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got := indexByBookingID(resp.Reservations)
		if _, ok := got["EMPTY-BOOKED-1"]; ok {
			t.Fatal("booked reservation must not appear in open reservations")
		}
		if _, ok := got["EMPTY-BOOKED-2"]; ok {
			t.Fatal("booked reservation must not appear in open reservations")
		}
	})

	t.Run("returns only open reservations with correct field values", func(t *testing.T) {
		const userID int32 = 7777
		ctx := context.Background()
		s := &Service{query: testQuerier()}

		// 1) "booked" reservation — should NOT appear in open reservations.
		seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
			p.BrokerReservationID = "OPEN-BOOKED"
			p.DriverTitle = "Mr"
			p.DriverFirstName = "Booked"
			p.DriverLastName = "User"
		})

		// 2) "vouchered" + "unpaid" — SHOULD appear.
		vouchID := seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
			p.BrokerReservationID = "OPEN-VOUCH"
			p.DriverTitle = "Mr"
			p.DriverFirstName = "Voucher"
			p.DriverLastName = "Holder"
			p.PickupDate = "2026-07-01"
			p.ReturnDate = "2026-07-05"
			p.RentalDays = 4
			p.CountryCode = "US"
			p.CurrencyCode = "USD"
			p.PurchasePrice = 100.00
			p.MarkupPercentage = 45.00
			p.BrokerErpPrice = 15.00
			p.BtErpPrice = 20
		})

		// 3) "canceled" + "refund_pending" — SHOULD appear.
		cancID := seedReservation(t, ctx, s, userID, func(p *CreateReservationRequest) {
			p.BrokerReservationID = "OPEN-CANC"
			p.DriverTitle = "Ms"
			p.DriverFirstName = "Cancel"
			p.DriverLastName = "Person"
			p.PickupDate = "2026-08-10"
			p.ReturnDate = "2026-08-13"
			p.RentalDays = 3
			p.CountryCode = "FR"
			p.CurrencyCode = "EUR"
			p.PurchasePrice = 200.00
			p.MarkupPercentage = 50.00
			p.BrokerErpPrice = 30.00
			p.BtErpPrice = 10
		})

		// Transition reservation states via the existing sqlc-generated queries.
		voucherNumber := "VCH-12345"
		if _, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{
			ID:            vouchID,
			UserID:        userID,
			VoucherNumber: &voucherNumber,
		}); err != nil {
			t.Fatalf("failed to apply voucher: %v", err)
		}
		if err := s.query.CancelReservation(ctx, cancID); err != nil {
			t.Fatalf("failed to cancel reservation: %v", err)
		}

		// Re-read to capture authoritative timestamp values.
		vouchRow, err := s.query.GetReservationByID(ctx, vouchID)
		if err != nil {
			t.Fatalf("failed to fetch vouchered reservation: %v", err)
		}
		cancRow, err := s.query.GetReservationByID(ctx, cancID)
		if err != nil {
			t.Fatalf("failed to fetch canceled reservation: %v", err)
		}

		resp, err := GetOpenReservations(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got := indexByBookingID(resp.Reservations)

		if _, ok := got["OPEN-BOOKED"]; ok {
			t.Fatal("booked reservation must not appear in open reservations")
		}

		wantVouch := OpenReservation{
			ID:                  vouchID,
			PaymentStatus:       "unpaid",
			BrokerReservationID: "OPEN-VOUCH",
			AgentID:             userID,
			CreatedAt:           db.TimestamptzToString(vouchRow.CreatedAt),
			VoucheredAt:         db.TimestamptzToString(vouchRow.VoucheredAt),
			VoucherNumber:       voucherNumber,
			DriverName:          "Mr Voucher Holder",
			PickupDate:          "2026-07-01",
			DropoffDate:         "2026-07-05",
			RentalDays:          4,
			CountryCode:         "US",
			CurrencyCode:        "USD",
			// 100 * 1.45 = 145
			CarPrice: 145.0,
			// 15 * 1.45 + 20 = 21.75 + 20 = 41.75
			ERPPrice: 41.75,
			// 145 + 41.75 = 186.75
			TotalPrice: 186.75,
		}
		assertOpenReservationEqual(t, "vouchered row", wantVouch, got["OPEN-VOUCH"])

		wantCanc := OpenReservation{
			ID:                  cancID,
			PaymentStatus:       "refund_pending",
			BrokerReservationID: "OPEN-CANC",
			AgentID:             userID,
			CreatedAt:           db.TimestamptzToString(cancRow.CreatedAt),
			VoucheredAt:         "",
			VoucherNumber:       "",
			DriverName:          "Ms Cancel Person",
			PickupDate:          "2026-08-10",
			DropoffDate:         "2026-08-13",
			RentalDays:          3,
			CountryCode:         "FR",
			CurrencyCode:        "EUR",
			// 200 * 1.5 = 300
			CarPrice: 300.0,
			// 30 * 1.5 + 10 = 45 + 10 = 55
			ERPPrice: 55.0,
			// 300 + 55 = 355
			TotalPrice: 355.0,
		}
		assertOpenReservationEqual(t, "canceled row", wantCanc, got["OPEN-CANC"])
	})
}

// indexByBookingID returns reservations keyed by their broker reservation id.
func indexByBookingID(rows []OpenReservation) map[string]OpenReservation {
	m := make(map[string]OpenReservation, len(rows))
	for _, r := range rows {
		m[r.BrokerReservationID] = r
	}
	return m
}

func assertOpenReservationEqual(t *testing.T, label string, want, got OpenReservation) {
	t.Helper()
	if got == (OpenReservation{}) {
		t.Fatalf("%s: not found in response", label)
	}
	if got.ID != want.ID {
		t.Errorf("%s: ID = %d, want %d", label, got.ID, want.ID)
	}
	if got.PaymentStatus != want.PaymentStatus {
		t.Errorf("%s: PaymentStatus = %q, want %q", label, got.PaymentStatus, want.PaymentStatus)
	}
	if got.BrokerReservationID != want.BrokerReservationID {
		t.Errorf("%s: BrokerReservationID = %q, want %q", label, got.BrokerReservationID, want.BrokerReservationID)
	}
	if got.AgentID != want.AgentID {
		t.Errorf("%s: AgentID = %d, want %d", label, got.AgentID, want.AgentID)
	}
	if got.CreatedAt != want.CreatedAt {
		t.Errorf("%s: CreatedAt = %q, want %q", label, got.CreatedAt, want.CreatedAt)
	}
	if got.VoucheredAt != want.VoucheredAt {
		t.Errorf("%s: VoucheredAt = %q, want %q", label, got.VoucheredAt, want.VoucheredAt)
	}
	if got.VoucherNumber != want.VoucherNumber {
		t.Errorf("%s: VoucherNumber = %q, want %q", label, got.VoucherNumber, want.VoucherNumber)
	}
	if got.DriverName != want.DriverName {
		t.Errorf("%s: DriverName = %q, want %q", label, got.DriverName, want.DriverName)
	}
	if got.PickupDate != want.PickupDate {
		t.Errorf("%s: PickupDate = %q, want %q", label, got.PickupDate, want.PickupDate)
	}
	if got.DropoffDate != want.DropoffDate {
		t.Errorf("%s: DropoffDate = %q, want %q", label, got.DropoffDate, want.DropoffDate)
	}
	if got.RentalDays != want.RentalDays {
		t.Errorf("%s: RentalDays = %d, want %d", label, got.RentalDays, want.RentalDays)
	}
	if got.CountryCode != want.CountryCode {
		t.Errorf("%s: CountryCode = %q, want %q", label, got.CountryCode, want.CountryCode)
	}
	if got.CurrencyCode != want.CurrencyCode {
		t.Errorf("%s: CurrencyCode = %q, want %q", label, got.CurrencyCode, want.CurrencyCode)
	}
	if got.CarPrice != want.CarPrice {
		t.Errorf("%s: CarPrice = %v, want %v", label, got.CarPrice, want.CarPrice)
	}
	if got.ERPPrice != want.ERPPrice {
		t.Errorf("%s: ERPPrice = %v, want %v", label, got.ERPPrice, want.ERPPrice)
	}
	if got.TotalPrice != want.TotalPrice {
		t.Errorf("%s: TotalPrice = %v, want %v", label, got.TotalPrice, want.TotalPrice)
	}
}
