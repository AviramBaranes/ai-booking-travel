package reservation

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/reports"
	"encore.dev/beta/auth"
)

type businessReportSeed struct {
	orgA      *accounts.OrganizationResponse
	officeA   *accounts.OfficeResponse
	agentA    *accounts.CreateAgentResponse
	admin     *accounts.CreateAdminResponse
	bookingA  string
	supplierA string

	orgB      *accounts.OrganizationResponse
	officeB   *accounts.OfficeResponse
	agentB    *accounts.CreateAgentResponse
	bookingB  string
	supplierB string

	reservationAID int64
	reservationBID int64
}

type businessesBalancesReportSeed struct {
	organicOrg      *accounts.OrganizationResponse
	inorganicOffice *accounts.OfficeResponse
}

func adminAuthContext(userID int64) context.Context {
	uid := auth.UID(strconv.FormatInt(userID, 10))
	return auth.WithContext(context.Background(), uid, &accounts.AuthData{
		UserID: userID,
		Role:   accounts.UserRoleAdmin,
	})
}

func TestGetBusinessReport(t *testing.T) {
	t.Run("validation rejects zero page", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("page"), reports.ReportParams{Page: 0, PageSize: 25}.Validate())
	})

	t.Run("validation rejects zero page size", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("pageSize"), reports.ReportParams{PageSize: 0, Page: 1}.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := (reports.ReportParams{Page: 1, PageSize: 25, UserType: "both"}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	ctx := adminAuthContext(900001)
	query := testQuerier()
	s := &Service{query: query}
	seed := seedBusinessReportData(t, ctx, s)
	baseParams := reports.ReportParams{
		Page:           1,
		PageSize:       25,
		PickupDateFrom: "2099-01-01",
		PickupDateTo:   "2099-12-31",
		UserType:       "both",
	}

	t.Run("returns account names and price fields", func(t *testing.T) {
		params := baseParams
		params.Supplier = seed.supplierA

		resp, err := GetBusinessReport(ctx, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Count != 1 {
			t.Fatalf("expected total 1, got %d", resp.Count)
		}
		if len(resp.Reservations) != 1 {
			t.Fatalf("expected 1 reservation, got %d", len(resp.Reservations))
		}

		assertFloatEqual(t, resp.TotalSales, 648)

		row := resp.Reservations[0]
		if row.BrokerReservationID != seed.bookingA {
			t.Fatalf("expected booking %q, got %q", seed.bookingA, row.BrokerReservationID)
		}
		if row.OrganizationName != seed.orgA.Name {
			t.Fatalf("expected organization name %q, got %q", seed.orgA.Name, row.OrganizationName)
		}
		if row.OfficeName != seed.officeA.Name {
			t.Fatalf("expected office name %q, got %q", seed.officeA.Name, row.OfficeName)
		}
		if row.UserName != "Report AgentA" {
			t.Fatalf("expected agent name %q, got %q", "Report AgentA", row.UserName)
		}
		if row.AdminName == nil || *row.AdminName != "Report Admin" {
			t.Fatalf("expected admin name %q, got %v", "Report Admin", row.AdminName)
		}
		if row.DriverName != "Ms Alice Report" {
			t.Fatalf("expected driver name %q, got %q", "Ms Alice Report", row.DriverName)
		}
		if row.SupplierName != "Hertz" {
			t.Fatalf("expected supplier name %q, got %q", "Hertz", row.SupplierName)
		}
		assertFloatEqual(t, row.CurrencyRate, 4)                     // seeded USD-to-ILS rate
		assertFloatEqual(t, row.CarSellPriceWithBrokerERP, 132)      // (100 purchase + 20 broker ERP) with 10% markup
		assertFloatEqual(t, row.CarSellPriceWithBrokerERPInILS, 528) // 132 * currency rate 4
		assertFloatEqual(t, row.BtERPPrice, 30)                      // seeded BT ERP fee
		assertFloatEqual(t, row.BtERPPriceInILS, 120)                // 30 * currency rate 4
		assertFloatEqual(t, row.TotalPrice, 162)                     // 132 sell price + 30 BT ERP, no discount
		assertFloatEqual(t, row.TotalPriceInILS, 648)                // 162 * currency rate 4
	})

	t.Run("filters by organization", func(t *testing.T) {
		params := baseParams
		params.OrganizationID = seed.orgA.ID
		params.Supplier = seed.supplierA
		assertBusinessReportBookings(t, ctx, params, seed.bookingA)
	})

	t.Run("filters by office", func(t *testing.T) {
		params := baseParams
		params.OfficeID = seed.officeA.ID
		params.Supplier = seed.supplierA
		assertBusinessReportBookings(t, ctx, params, seed.bookingA)
	})

	t.Run("filters by agent", func(t *testing.T) {
		params := baseParams
		params.AgentID = seed.agentA.ID
		params.Supplier = seed.supplierA
		assertBusinessReportBookings(t, ctx, params, seed.bookingA)
	})

	t.Run("filters by broker", func(t *testing.T) {
		params := baseParams
		params.Broker = "hertz"
		params.Supplier = seed.supplierB
		assertBusinessReportBookings(t, ctx, params, seed.bookingB)
	})

	t.Run("filters by status", func(t *testing.T) {
		params := baseParams
		params.Status = "canceled"
		params.Supplier = seed.supplierB
		assertBusinessReportBookings(t, ctx, params, seed.bookingB)
	})

	t.Run("filters by pickup date range", func(t *testing.T) {
		params := reports.ReportParams{
			Page:           1,
			PickupDateFrom: "2099-02-01",
			PickupDateTo:   "2099-02-28",
			Supplier:       seed.supplierB,
			PageSize:       25,
			UserType:       "both",
		}
		assertBusinessReportBookings(t, ctx, params, seed.bookingB)
	})

	t.Run("non-matching supplier returns empty", func(t *testing.T) {
		params := baseParams
		params.Supplier = fmt.Sprintf("missing_supplier_%d", time.Now().UnixNano())

		resp, err := GetBusinessReport(ctx, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Count != 0 {
			t.Fatalf("expected total 0, got %d", resp.Count)
		}
		if len(resp.Reservations) != 0 {
			t.Fatalf("expected 0 reservations, got %d", len(resp.Reservations))
		}
	})

	t.Run("userType=agent shows only business reservations", func(t *testing.T) {
		params := baseParams
		params.UserType = "agent"
		params.Supplier = seed.supplierA
		assertBusinessReportBookings(t, ctx, params, seed.bookingA)
	})

	t.Run("userType=customer returns empty for business reservations", func(t *testing.T) {
		params := baseParams
		params.UserType = "customer"
		params.Supplier = seed.supplierA

		resp, err := GetBusinessReport(ctx, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Count != 0 {
			t.Fatalf("expected 0 for customer filter on business reservation, got %d", resp.Count)
		}
	})
}

func TestGetProfitReport(t *testing.T) {
	ctx := adminAuthContext(900002)
	query := testQuerier()
	s := &Service{query: query}
	seed := seedBusinessReportData(t, ctx, s)
	v := "TEST-VOUCHER"
	query.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            seed.reservationAID,
		UserID:        seed.agentA.ID,
		VoucherNumber: &v,
		CurrencyRate:  dbadapters.NumericFromFloat64(4.0),
	})

	t.Run("returns shared report fields and profit fields", func(t *testing.T) {
		resp, err := GetProfitReport(ctx, reports.ReportParams{
			Page:           1,
			PageSize:       25,
			PickupDateFrom: "2099-01-01",
			PickupDateTo:   "2099-12-31",
			Supplier:       seed.supplierA,
			UserType:       "both",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Count != 1 {
			t.Fatalf("expected total 1, got %d", resp.Count)
		}
		if len(resp.Reservations) != 1 {
			t.Fatalf("expected 1 reservation, got %d", len(resp.Reservations))
		}

		row := resp.Reservations[0]
		if row.BrokerReservationID != seed.bookingA {
			t.Fatalf("expected booking %q, got %q", seed.bookingA, row.BrokerReservationID)
		}
		if row.OrganizationName != seed.orgA.Name {
			t.Fatalf("expected organization name %q, got %q", seed.orgA.Name, row.OrganizationName)
		}
		if row.OfficeName != seed.officeA.Name {
			t.Fatalf("expected office name %q, got %q", seed.officeA.Name, row.OfficeName)
		}
		assertFloatEqual(t, row.CarSellPriceWithBrokerERP, 132)      // (100 purchase + 20 broker ERP) with 10% markup
		assertFloatEqual(t, row.CarSellPriceWithBrokerERPInILS, 528) // 132 * currency rate 4
		assertFloatEqual(t, row.PurchasePrice, 120)                  // 100 purchase + 20 broker ERP
		assertFloatEqual(t, row.PurchasePriceInILS, 480)             // 120 * currency rate 4
		assertFloatEqual(t, row.Profit, 42)                          // 132 sell price - 120 purchase price + 30 bt ERP
		assertFloatEqual(t, row.ProfitInILS, 168)                    // 42 * currency rate 4

		assertFloatEqual(t, resp.TotalSales, 648)
		assertFloatEqual(t, resp.TotalProfit, 168)
		assertFloatEqual(t, resp.ProfitPercentage, (168.0/648.0)*100.0)
	})
}

func TestGetBusinessesBalancesReport(t *testing.T) {
	ctx := adminAuthContext(900003)
	query := testQuerier()
	s := &Service{query: query}
	seed := seedBusinessesBalancesReportData(t, ctx, s)

	t.Run("returns billing entity names and currency-bucketed balances", func(t *testing.T) {
		resp, err := GetBusinessesBalancesReport(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total < 2 {
			t.Fatalf("expected at least 2 businesses, got %d", resp.Total)
		}

		// Organic org: 2 USD (vouchered 77 + canceled 120 = 197), 1 EUR (100), 1 ILS canceled (50 * 3.6 = 180)
		organicRow := requireBusinessesBalanceRow(t, resp.Businesses, reports.BillingEntityBusiness, seed.organicOrg.ID)
		if organicRow.BillingEntityName != seed.organicOrg.Name {
			t.Fatalf("expected organic billing entity name %q, got %q", seed.organicOrg.Name, organicRow.BillingEntityName)
		}
		assertFloatEqual(t, organicRow.TotalOpenBalanceInDollar, -43) // 77 (vouchered) + 120 (canceled)
		assertFloatEqual(t, organicRow.TotalOpenBalanceInEuro, 100)   // EUR vouchered reservation
		assertFloatEqual(t, organicRow.TotalInOtherCurrency, -180)    // 50 (ILS canceled) * rate 3.6

		// Inorganic office: 2 USD only (vouchered 220 + canceled 33 = 253)
		inorganicOfficeRow := requireBusinessesBalanceRow(t, resp.Businesses, reports.BillingEntityOffice, seed.inorganicOffice.ID)
		if inorganicOfficeRow.BillingEntityName != seed.inorganicOffice.Name {
			t.Fatalf("expected inorganic office billing entity name %q, got %q", seed.inorganicOffice.Name, inorganicOfficeRow.BillingEntityName)
		}
		assertFloatEqual(t, inorganicOfficeRow.TotalOpenBalanceInDollar, 187) // 220 (vouchered) + 33 (canceled)
		assertFloatEqual(t, inorganicOfficeRow.TotalOpenBalanceInEuro, 0)
		assertFloatEqual(t, inorganicOfficeRow.TotalInOtherCurrency, 0)
	})
}

func seedBusinessReportData(t *testing.T, ctx context.Context, s *Service) businessReportSeed {
	t.Helper()
	unique := time.Now().UnixNano()
	icountClientID := int32(100)

	orgA := createBusinessReportOrg(t, ctx, fmt.Sprintf("Report Org A %d", unique), icountClientID)
	officeA := createBusinessReportOffice(t, ctx, orgA.ID, fmt.Sprintf("Report Office A %d", unique))
	agentA := createBusinessReportAgent(t, ctx, officeA.ID, "AgentA", fmt.Sprintf("report_agent_a_%d@test.com", unique), unique%100000000)
	admin := createBusinessReportAdmin(t, ctx, fmt.Sprintf("report_admin_%d@test.com", unique))

	orgB := createBusinessReportOrg(t, ctx, fmt.Sprintf("Report Org B %d", unique), icountClientID+1)
	officeB := createBusinessReportOffice(t, ctx, orgB.ID, fmt.Sprintf("Report Office B %d", unique))
	agentB := createBusinessReportAgent(t, ctx, officeB.ID, "AgentB", fmt.Sprintf("report_agent_b_%d@test.com", unique), (unique+1)%100000000)

	bookingA := fmt.Sprintf("REPORT-A-%d", unique)
	supplierA := fmt.Sprintf("REPORT-SUP-A-%d", unique)
	adminID := admin.ID
	isOrgAOrganic := orgA.IsOrganic
	reservationAID := seedReservation(t, ctx, s, agentA.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = bookingA
		p.OfficeID = &officeA.ID
		p.OrganizationID = &orgA.ID
		p.IsOrganizationOrganic = &isOrgAOrganic
		p.AdminRefID = &adminID
		p.Broker = "flex"
		p.SupplierCode = supplierA
		p.PickupDate = "2099-01-10"
		p.DropoffDate = "2099-01-14"
		p.DriverTitle = "Ms"
		p.DriverFirstName = "Alice"
		p.DriverLastName = "Report"
		p.CurrencyRate = 4
		p.PurchasePrice = 100
		p.MarkupPercentage = 10
		p.BrokerErpPrice = 20
		p.BtErpPrice = 30
		p.DiscountPercentage = 0
	})

	bookingB := fmt.Sprintf("REPORT-B-%d", unique)
	supplierB := fmt.Sprintf("REPORT-SUP-B-%d", unique)
	isOrgBOrganic := orgB.IsOrganic
	reservationBID := seedReservation(t, ctx, s, agentB.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = bookingB
		p.OfficeID = &officeB.ID
		p.OrganizationID = &orgB.ID
		p.IsOrganizationOrganic = &isOrgBOrganic
		p.Broker = "hertz"
		p.SupplierCode = supplierB
		p.PickupDate = "2099-02-10"
		p.DropoffDate = "2099-02-15"
		p.DriverFirstName = "Bob"
		p.DriverLastName = "Report"
	})
	if err := s.query.CancelReservation(ctx, reservationBID); err != nil {
		t.Fatalf("failed to cancel seeded reservation: %v", err)
	}

	return businessReportSeed{
		orgA:           orgA,
		officeA:        officeA,
		agentA:         agentA,
		admin:          admin,
		bookingA:       bookingA,
		supplierA:      supplierA,
		orgB:           orgB,
		officeB:        officeB,
		agentB:         agentB,
		bookingB:       bookingB,
		supplierB:      supplierB,
		reservationAID: reservationAID,
		reservationBID: reservationBID,
	}
}

func seedBusinessesBalancesReportData(t *testing.T, ctx context.Context, s *Service) businessesBalancesReportSeed {
	t.Helper()
	unique := time.Now().UnixNano()

	organicOrg := createBusinessesBalancesOrg(t, ctx, fmt.Sprintf("Balances Organic Org %d", unique), true, int32(200))
	organicOffice := createBusinessReportOffice(t, ctx, organicOrg.ID, fmt.Sprintf("Balances Organic Office %d", unique))
	organicAgent := createBusinessReportAgent(t, ctx, organicOffice.ID, "BalancesOrganic", fmt.Sprintf("balances_organic_%d@test.com", unique), unique%100000000)
	organic := true

	seedReservation(t, ctx, s, organicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-ORG-OPEN-%d", unique)
		p.OfficeID = &organicOffice.ID
		p.OrganizationID = &organicOrg.ID
		p.IsOrganizationOrganic = &organic
		p.CurrencyRate = 4
		p.PurchasePrice = 100
		p.MarkupPercentage = 10
		p.BrokerErpPrice = 20
		p.BtErpPrice = 30
		p.DiscountPercentage = 0
	})

	organicVoucheredID := seedReservation(t, ctx, s, organicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-ORG-VOUCHERED-%d", unique)
		p.OfficeID = &organicOffice.ID
		p.OrganizationID = &organicOrg.ID
		p.IsOrganizationOrganic = &organic
		p.CurrencyRate = 3
		p.PurchasePrice = 50
		p.MarkupPercentage = 20
		p.BrokerErpPrice = 10
		p.BtErpPrice = 5
		p.DiscountPercentage = 0
	})
	applyVoucherForTest(t, ctx, s, organicVoucheredID, organicAgent.ID, fmt.Sprintf("BALANCES-ORG-VN-%d", unique))

	organicCanceledID := seedReservation(t, ctx, s, organicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-ORG-CANCELED-%d", unique)
		p.OfficeID = &organicOffice.ID
		p.OrganizationID = &organicOrg.ID
		p.IsOrganizationOrganic = &organic
		p.CurrencyRate = 2
		p.PurchasePrice = 80
		p.MarkupPercentage = 25
		p.BrokerErpPrice = 0
		p.BtErpPrice = 20
		p.DiscountPercentage = 0
	})
	cancelReservationForTest(t, ctx, s, organicCanceledID)

	// EUR vouchered — exercises total_eur bucket
	organicEurVoucheredID := seedReservation(t, ctx, s, organicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-ORG-EUR-%d", unique)
		p.OfficeID = &organicOffice.ID
		p.OrganizationID = &organicOrg.ID
		p.IsOrganizationOrganic = &organic
		p.CurrencyCode = "EUR"
		p.CurrencyRate = 1
		p.PurchasePrice = 100
		p.MarkupPercentage = 0
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 0
	})
	applyVoucherForTest(t, ctx, s, organicEurVoucheredID, organicAgent.ID, fmt.Sprintf("BALANCES-ORG-EUR-VN-%d", unique))

	// ILS canceled — exercises total_other_converted bucket (50 * 3.6 = 180)
	organicIlsCanceledID := seedReservation(t, ctx, s, organicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-ORG-ILS-%d", unique)
		p.OfficeID = &organicOffice.ID
		p.OrganizationID = &organicOrg.ID
		p.IsOrganizationOrganic = &organic
		p.CurrencyCode = "ILS"
		p.CurrencyRate = 3.6
		p.PurchasePrice = 50
		p.MarkupPercentage = 0
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 0
	})
	cancelReservationForTest(t, ctx, s, organicIlsCanceledID)

	inorganicOrg := createBusinessesBalancesOrg(t, ctx, fmt.Sprintf("Balances Inorganic Org %d", unique), false, 0)
	inorganicOffice := createBusinessesBalancesOffice(t, ctx, inorganicOrg.ID, fmt.Sprintf("Balances Inorganic Office %d", unique), int32(300))
	inorganicAgent := createBusinessReportAgent(t, ctx, inorganicOffice.ID, "BalancesInorganic", fmt.Sprintf("balances_inorganic_%d@test.com", unique), (unique+1)%100000000)
	inorganic := false

	seedReservation(t, ctx, s, inorganicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-OFFICE-OPEN-%d", unique)
		p.OfficeID = &inorganicOffice.ID
		p.OrganizationID = &inorganicOrg.ID
		p.IsOrganizationOrganic = &inorganic
		p.CurrencyRate = 2
		p.PurchasePrice = 100
		p.MarkupPercentage = 10
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 0
	})

	inorganicVoucheredID := seedReservation(t, ctx, s, inorganicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-OFFICE-VOUCHERED-%d", unique)
		p.OfficeID = &inorganicOffice.ID
		p.OrganizationID = &inorganicOrg.ID
		p.IsOrganizationOrganic = &inorganic
		p.CurrencyRate = 1
		p.PurchasePrice = 200
		p.MarkupPercentage = 10
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 0
	})
	applyVoucherForTest(t, ctx, s, inorganicVoucheredID, inorganicAgent.ID, fmt.Sprintf("BALANCES-OFFICE-VN-%d", unique))

	inorganicCanceledID := seedReservation(t, ctx, s, inorganicAgent.ID, func(p *CreateReservationParams) {
		p.BrokerReservationID = fmt.Sprintf("BALANCES-OFFICE-CANCELED-%d", unique)
		p.OfficeID = &inorganicOffice.ID
		p.OrganizationID = &inorganicOrg.ID
		p.IsOrganizationOrganic = &inorganic
		p.CurrencyRate = 5
		p.PurchasePrice = 30
		p.MarkupPercentage = 10
		p.BrokerErpPrice = 0
		p.BtErpPrice = 0
		p.DiscountPercentage = 0
	})
	cancelReservationForTest(t, ctx, s, inorganicCanceledID)

	return businessesBalancesReportSeed{
		organicOrg:      organicOrg,
		inorganicOffice: inorganicOffice,
	}
}

func createBusinessReportOrg(t *testing.T, ctx context.Context, name string, icountClientID int32) *accounts.OrganizationResponse {
	t.Helper()
	org, err := accounts.CreateOrganization(ctx, accounts.CreateOrganizationParams{
		Name:           name,
		IsOrganic:      true,
		IcountClientID: &icountClientID,
	})
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}
	return org
}

func createBusinessesBalancesOrg(t *testing.T, ctx context.Context, name string, isOrganic bool, icountClientID int32) *accounts.OrganizationResponse {
	t.Helper()
	var icountClientIDPtr *int32
	if isOrganic {
		icountClientIDPtr = &icountClientID
	}
	org, err := accounts.CreateOrganization(ctx, accounts.CreateOrganizationParams{
		Name:           name,
		IsOrganic:      isOrganic,
		IcountClientID: icountClientIDPtr,
	})
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}
	return org
}

func createBusinessReportOffice(t *testing.T, ctx context.Context, orgID int64, name string) *accounts.OfficeResponse {
	t.Helper()
	office, err := accounts.CreateOffice(ctx, accounts.CreateOfficeParams{
		Name:           name,
		OrganizationID: orgID,
	})
	if err != nil {
		t.Fatalf("failed to create office: %v", err)
	}
	return office
}

func createBusinessesBalancesOffice(t *testing.T, ctx context.Context, orgID int64, name string, icountClientID int32) *accounts.OfficeResponse {
	t.Helper()
	office, err := accounts.CreateOffice(ctx, accounts.CreateOfficeParams{
		Name:           name,
		OrganizationID: orgID,
		IcountClientID: &icountClientID,
	})
	if err != nil {
		t.Fatalf("failed to create office: %v", err)
	}
	return office
}

func createBusinessReportAgent(t *testing.T, ctx context.Context, officeID int64, lastName string, email string, phoneSuffix int64) *accounts.CreateAgentResponse {
	t.Helper()
	agent, err := accounts.CreateAgent(ctx, accounts.CreateAgentParams{
		FirstName:   "Report",
		LastName:    lastName,
		Email:       email,
		Password:    "ValidPass123!",
		PhoneNumber: fmt.Sprintf("05%08d", phoneSuffix),
		OfficeID:    officeID,
	})
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}
	return agent
}

func createBusinessReportAdmin(t *testing.T, ctx context.Context, email string) *accounts.CreateAdminResponse {
	t.Helper()
	admin, err := accounts.CreateAdmin(ctx, accounts.CreateAdminParams{
		FirstName: "Report",
		LastName:  "Admin",
		Email:     email,
		Password:  "ValidPass123!",
	})
	if err != nil {
		t.Fatalf("failed to create admin: %v", err)
	}
	return admin
}

func assertBusinessReportBookings(t *testing.T, ctx context.Context, params reports.ReportParams, wantBookings ...string) {
	t.Helper()
	resp, err := GetBusinessReport(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Count != int64(len(wantBookings)) {
		t.Fatalf("expected total %d, got %d", len(wantBookings), resp.Count)
	}
	if len(resp.Reservations) != len(wantBookings) {
		t.Fatalf("expected %d reservations, got %d", len(wantBookings), len(resp.Reservations))
	}

	got := make(map[string]bool, len(resp.Reservations))
	for _, row := range resp.Reservations {
		got[row.BrokerReservationID] = true
	}
	for _, booking := range wantBookings {
		if !got[booking] {
			t.Fatalf("expected booking %q in report response", booking)
		}
	}
}

func applyVoucherForTest(t *testing.T, ctx context.Context, s *Service, reservationID int64, userID int64, voucherNumber string) {
	t.Helper()
	reserv, err := s.query.GetReservationByID(ctx, reservationID)
	if err != nil {
		t.Fatalf("failed to get reservation: %v", err)
	}
	if err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{ID: reservationID, UserID: userID, VoucherNumber: &voucherNumber, CurrencyRate: reserv.CurrencyRate}); err != nil {
		t.Fatalf("failed to apply voucher: %v", err)
	}
}

func cancelReservationForTest(t *testing.T, ctx context.Context, s *Service, reservationID int64) {
	t.Helper()
	if err := s.query.ResolveReservationsPayment(ctx, []int64{reservationID}); err != nil {
		t.Fatalf("failed to resolve reservation payment: %v", err)
	}
	if err := s.query.CancelReservation(ctx, reservationID); err != nil {
		t.Fatalf("failed to cancel reservation: %v", err)
	}
}

func requireBusinessesBalanceRow(t *testing.T, rows []reports.BusinessesBalancesReportRow, entityType reports.BillingEntity, entityID int64) reports.BusinessesBalancesReportRow {
	t.Helper()
	for _, row := range rows {
		if row.BillingEntityType == entityType && row.BillingEntityID == entityID {
			return row
		}
	}
	t.Fatalf("expected businesses balances row for %s %d", entityType, entityID)
	return reports.BusinessesBalancesReportRow{}
}

func assertFloatEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Fatalf("expected %.4f, got %.4f", want, got)
	}
}
