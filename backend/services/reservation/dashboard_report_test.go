package reservation

import (
	"fmt"
	"testing"
	"time"

	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/reports"
)

// today returns the current calendar date in the business timezone, which is how the
// dashboard interprets its from/to filters.
func today(t *testing.T) string {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatalf("failed to load business timezone: %v", err)
	}
	return time.Now().In(loc).Format("2006-01-02")
}

func findDashboardRow(t *testing.T, rows []reports.DashboardReservation, reservationID int64) reports.DashboardReservation {
	t.Helper()
	for _, row := range rows {
		if row.ReservationID == reservationID {
			return row
		}
	}
	t.Fatalf("expected dashboard row for reservation %d", reservationID)
	return reports.DashboardReservation{}
}

func TestGetDashboardReportValidation(t *testing.T) {
	t.Run("rejects missing from", func(t *testing.T) {
		p := &reports.DashboardParams{To: "2026-08-12"}
		if err := p.Validate(); err == nil {
			t.Fatal("expected an error for a missing from date")
		}
	})

	t.Run("rejects missing to", func(t *testing.T) {
		p := &reports.DashboardParams{From: "2026-08-12"}
		if err := p.Validate(); err == nil {
			t.Fatal("expected an error for a missing to date")
		}
	})

	t.Run("rejects a non-date value", func(t *testing.T) {
		p := &reports.DashboardParams{From: "12/08/2026", To: "2026-08-12"}
		if err := p.Validate(); err == nil {
			t.Fatal("expected an error for a malformed from date")
		}
	})

	t.Run("rejects an inverted range", func(t *testing.T) {
		p := &reports.DashboardParams{From: "2026-08-12", To: "2026-08-01"}
		if err := p.Validate(); err == nil {
			t.Fatal("expected an error when to precedes from")
		}
	})

	t.Run("accepts a single-day range", func(t *testing.T) {
		p := &reports.DashboardParams{From: "2026-08-12", To: "2026-08-12"}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestGetDashboardReport(t *testing.T) {
	ctx := adminAuthContext(900101)
	query := testQuerier()
	s := &Service{query: query}
	unique := time.Now().UnixNano()
	day := today(t)

	org := createBusinessReportOrg(t, ctx, fmt.Sprintf("Dashboard Org %d", unique), int32(200))
	office := createBusinessReportOffice(t, ctx, org.ID, fmt.Sprintf("Dashboard Office %d", unique))
	agent := createBusinessReportAgent(t, ctx, office.ID, "DashAgent", fmt.Sprintf("dash_agent_%d@test.com", unique), unique%100000000)
	isOrganic := org.IsOrganic

	// An agent booking: automatic gearbox, ERP included, no coupon.
	agentReservationID := seedReservation(t, ctx, s, agent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("DASH-AGENT-%d", unique)
		p.OfficeID = &office.ID
		p.OrganizationID = &org.ID
		p.IsOrganizationOrganic = &isOrganic
		p.Broker = "flex"
		p.CountryCode = "IL"
		p.CurrencyCode = "USD"
		p.CurrencyRate = 4
		p.PurchasePrice = 100
		p.MarkupPercentage = 50
		p.BrokerErpPrice = 20
		p.BtErpPrice = 30
		p.DiscountPercentage = 0
		p.DriverAge = 41
		p.RentalDays = 4
	})

	// A customer booking: no office/organization, manual gearbox, no ERP, coupon applied.
	customerReservationID := seedReservation(t, ctx, s, agent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("DASH-CUSTOMER-%d", unique)
		p.OfficeID = nil
		p.OrganizationID = nil
		p.IsOrganizationOrganic = nil
		p.Broker = "hertz"
		p.CountryCode = "GR"
		p.CurrencyRate = 1
		p.PurchasePrice = 200
		p.MarkupPercentage = 25
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 10
		p.CouponName = "SUMMER10"
		p.DriverAge = 27
		p.CarDetails = &broker.CarDetails{
			Model:        "VW Golf",
			CarGroup:     "Compact",
			SupplierName: "Avis",
			CarType:      "Hatchback",
			IsAutoGear:   false,
			IsElectric:   false,
		}
	})

	resp, err := GetDashboardReport(ctx, &reports.DashboardParams{From: day, To: day})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("classifies agent bookings as business", func(t *testing.T) {
		row := findDashboardRow(t, resp.Reservations, agentReservationID)

		if !row.IsBusiness {
			t.Fatal("expected the agent reservation to be classified as business")
		}
		if row.OrganizationID == nil || *row.OrganizationID != org.ID {
			t.Fatalf("expected organization %d, got %v", org.ID, row.OrganizationID)
		}
		if row.OfficeID == nil || *row.OfficeID != office.ID {
			t.Fatalf("expected office %d, got %v", office.ID, row.OfficeID)
		}
	})

	t.Run("classifies bookings without an office as private", func(t *testing.T) {
		row := findDashboardRow(t, resp.Reservations, customerReservationID)

		if row.IsBusiness {
			t.Fatal("expected the reservation without an office to be classified as private")
		}
	})

	t.Run("derives gear type from car details", func(t *testing.T) {
		if got := findDashboardRow(t, resp.Reservations, agentReservationID).GearType; got != reports.GearTypeAuto {
			t.Fatalf("expected gear type %q, got %q", reports.GearTypeAuto, got)
		}
		if got := findDashboardRow(t, resp.Reservations, customerReservationID).GearType; got != reports.GearTypeManual {
			t.Fatalf("expected gear type %q, got %q", reports.GearTypeManual, got)
		}
	})

	t.Run("reports dimensions taken from car details", func(t *testing.T) {
		row := findDashboardRow(t, resp.Reservations, customerReservationID)

		if row.SupplierName != "Avis" {
			t.Fatalf("expected supplier name %q, got %q", "Avis", row.SupplierName)
		}
		if row.CarType != "Hatchback" {
			t.Fatalf("expected car type %q, got %q", "Hatchback", row.CarType)
		}
		if row.CarGroup != "Compact" {
			t.Fatalf("expected car group %q, got %q", "Compact", row.CarGroup)
		}
		if row.CountryCode != "GR" {
			t.Fatalf("expected country code %q, got %q", "GR", row.CountryCode)
		}
		if row.CouponName != "SUMMER10" {
			t.Fatalf("expected coupon name %q, got %q", "SUMMER10", row.CouponName)
		}
	})

	t.Run("converts money to ILS and derives profit net of the discount", func(t *testing.T) {
		row := findDashboardRow(t, resp.Reservations, agentReservationID)

		// purchase 100 + broker ERP 20 = 120 cost; marked up 50% = 180; plus BT ERP 30 = 210 revenue.
		assertFloatEqual(t, row.CostILS, 120*4)
		assertFloatEqual(t, row.RevenueILS, 210*4)
		assertFloatEqual(t, row.ProfitILS, (210-120)*4)
		assertFloatEqual(t, row.ErpRevenueILS, 30*4)
		assertFloatEqual(t, row.ErpCostILS, 20*4)
		assertFloatEqual(t, row.DiscountILS, 0)

		if !row.HasERP {
			t.Fatal("expected the reservation to be marked as including ERP")
		}
	})

	t.Run("subtracts the coupon discount from revenue and profit", func(t *testing.T) {
		row := findDashboardRow(t, resp.Reservations, customerReservationID)

		// purchase 200 cost; marked up 25% = 250; 10% coupon = 225 revenue.
		assertFloatEqual(t, row.CostILS, 200)
		assertFloatEqual(t, row.RevenueILS, 225)
		assertFloatEqual(t, row.ProfitILS, 25)
		assertFloatEqual(t, row.DiscountILS, 25)

		if row.HasERP {
			t.Fatal("expected the reservation without ERP prices to be marked as excluding ERP")
		}
	})

	t.Run("keeps profit equal to revenue minus cost on every row", func(t *testing.T) {
		for _, row := range resp.Reservations {
			assertFloatEqual(t, row.ProfitILS, row.RevenueILS-row.CostILS)
		}
	})

	t.Run("returns the entity names the rows reference", func(t *testing.T) {
		if !containsAccountName(resp.Organizations, org.ID, org.Name) {
			t.Fatalf("expected organization %d (%s) in the lookup", org.ID, org.Name)
		}
		if !containsAccountName(resp.Offices, office.ID, office.Name) {
			t.Fatalf("expected office %d (%s) in the lookup", office.ID, office.Name)
		}
		if !containsAccountName(resp.Users, agent.ID, "Report DashAgent") {
			t.Fatalf("expected user %d in the lookup", agent.ID)
		}
	})

	t.Run("excludes reservations created outside the window", func(t *testing.T) {
		past, err := GetDashboardReport(ctx, &reports.DashboardParams{From: "2000-01-01", To: "2000-01-31"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, row := range past.Reservations {
			if row.ReservationID == agentReservationID || row.ReservationID == customerReservationID {
				t.Fatalf("reservation %d leaked into a window it was not created in", row.ReservationID)
			}
		}
	})

	t.Run("reports supplier payment state", func(t *testing.T) {
		if findDashboardRow(t, resp.Reservations, agentReservationID).SupplierPaid {
			t.Fatal("expected a freshly created reservation to be unpaid to the supplier")
		}

		if err := s.query.MarkReservationsPaidToSupplier(ctx, db.MarkReservationsPaidToSupplierParams{
			SupplierPaidAt: dbadapters.DBTime(time.Now()),
			Ids:            []int64{agentReservationID},
		}); err != nil {
			t.Fatalf("failed to mark the reservation paid to supplier: %v", err)
		}

		paid, err := GetDashboardReport(ctx, &reports.DashboardParams{From: day, To: day})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !findDashboardRow(t, paid.Reservations, agentReservationID).SupplierPaid {
			t.Fatal("expected the reservation to be reported as paid to the supplier")
		}
	})

	t.Run("folds a penalty into the reservation row", func(t *testing.T) {
		if err := s.query.CancelReservation(ctx, customerReservationID); err != nil {
			t.Fatalf("failed to cancel the reservation: %v", err)
		}
		if _, err := s.query.InsertReservationPenalty(ctx, db.InsertReservationPenaltyParams{
			ReservationID: customerReservationID,
			PenaltyType:   db.PenaltyTypeCancellation,
			CurrencyCode:  "USD",
			CurrencyRate:  dbadapters.NumericFromFloat64(4),
			Amount:        dbadapters.NumericFromFloat64(50),
		}); err != nil {
			t.Fatalf("failed to insert a penalty: %v", err)
		}

		withPenalty, err := GetDashboardReport(ctx, &reports.DashboardParams{From: day, To: day})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		row := findDashboardRow(t, withPenalty.Reservations, customerReservationID)
		if row.PenaltyType != string(db.PenaltyTypeCancellation) {
			t.Fatalf("expected penalty type %q, got %q", db.PenaltyTypeCancellation, row.PenaltyType)
		}
		assertFloatEqual(t, row.PenaltyAmountILS, 200)
		if row.PenaltyPaid {
			t.Fatal("expected a new penalty to be unpaid")
		}
		if row.PenaltySupplierPaid {
			t.Fatal("expected a new penalty to be unpaid to the supplier")
		}
		if row.Status != string(db.ReservationStatusCanceled) {
			t.Fatalf("expected status %q, got %q", db.ReservationStatusCanceled, row.Status)
		}
	})
}

func containsAccountName(names []accounts.AccountName, id int64, expected string) bool {
	for _, name := range names {
		if name.ID == id && name.Name == expected {
			return true
		}
	}
	return false
}
