package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	authpkg "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

// --- Helpers ---

func priceOfferInvalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

func priceOfferAuthContext(userID int32) context.Context {
	uid := auth.UID(strconv.Itoa(int(userID)))
	return auth.WithContext(context.Background(), uid, &authpkg.AuthData{
		UserID: userID,
		Role:   authpkg.UserRoleAgent,
	})
}

var locSeqCounter int64

func nextLocSeq() int64 {
	locSeqCounter++
	return locSeqCounter
}

// uniqueLocName returns a unique location name so tests can seed many locations
// without violating the country_code+name unique index.
func uniqueLocName(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), nextLocSeq())
}

// seedPriceOfferLocation inserts a location and returns its stringified ID
// (which is what is stored in price_offers.pickup_location_id / dropoff_location_id)
// alongside the location name for assertions on JOIN-derived fields.
func seedPriceOfferLocation(t *testing.T, q *db.Queries, namePrefix string) (string, string) {
	t.Helper()
	ctx := context.Background()
	city := "TestCity"
	name := uniqueLocName(namePrefix)
	loc, err := q.InsertLocation(ctx, db.InsertLocationParams{
		Country:     "TestCountry",
		CountryCode: "TC",
		City:        &city,
		Name:        name,
	})
	if err != nil {
		t.Fatalf("failed to seed location: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteLocationByID(ctx, loc.ID) })
	return strconv.FormatInt(loc.ID, 10), loc.Name
}

// seedSnapshot inserts an available_plans_snapshots row with the given plans
// and registers cleanup. Returns the snapshot ID.
func seedSnapshot(t *testing.T, q *db.Queries, plans []planPriceDetails) int64 {
	t.Helper()
	ctx := context.Background()

	plansJSON, err := json.Marshal(plans)
	if err != nil {
		t.Fatalf("failed to marshal plans: %v", err)
	}

	id, err := q.InsertAvailablePlansSnapshot(ctx, db.InsertAvailablePlansSnapshotParams{
		Plans:       plansJSON,
		DriverAge:   "30",
		PickupDate:  db.DateFromString("2026-08-01"),
		PickupTime:  "08:00",
		ReturnDate:  db.DateFromString("2026-08-05"),
		ReturnTime:  "10:00",
		CountryCode: "US",
	})
	if err != nil {
		t.Fatalf("failed to seed snapshot: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteSnapshotByID(ctx, id) })
	return id
}

// defaultPlan returns a fully populated plan that can be used to seed a snapshot.
func defaultPlan(pickupLocCode, dropoffLocCode string) planPriceDetails {
	return planPriceDetails{
		PlanID:                 1,
		RateQualifier:          "RQ1",
		SupplierCode:           "SUP1",
		Broker:                 broker.BrokerFlex,
		PickupLocationCode:     pickupLocCode,
		DropoffLocationCode:    dropoffLocCode,
		CurrencyCode:           "USD",
		CurrencyRate:           3.65,
		DiscountPercentage:     10,
		CarPurchasePrice:       100,
		SupplierErpPrice:       10,
		MarkupPercentage:       50,
		ChargedERPPriceWithVat: 15,
		CarDetails: broker.CarDetails{
			Model:        "Toyota Corolla",
			CarGroup:     "Economy",
			ImageURL:     "https://example.com/car.png",
			SupplierName: "Hertz",
			CarType:      "Sedan",
			Acriss:       "CDMR",
			HasAC:        true,
			IsAutoGear:   true,
			IsElectric:   false,
			Seats:        5,
			Bags:         2,
			Doors:        4,
		},
		Inclusions: []string{"Unlimited Mileage", "Collision Damage Waiver"},
	}
}

// validCreatePriceOfferParams returns valid params for CreatePriceOffer pointing to the
// supplied plan and snapshot ID.
func validCreatePriceOfferParams(snapshotID int64, plan planPriceDetails) CreatePriceOfferParams {
	return CreatePriceOfferParams{
		SnapshotID:          snapshotID,
		RateQualifier:       plan.RateQualifier,
		SupplierCode:        plan.SupplierCode,
		IncludeERP:          true,
		Name:                "Test Offer",
		OfferedCurrencyCode: "USD",
		OfferedPrice:        200,
	}
}

// snapshotNotFoundErr returns the expected snapshot-not-found error for assertions.
func snapshotNotFoundErr() error {
	return api_errors.NewErrorWithDetail(errs.NotFound, "Snapshot not found",
		api_errors.ErrorDetails{Code: api_errors.CodeSnapshotNotFound})
}

// planNotFoundErr returns the expected plan-not-found error for assertions.
func planNotFoundErr() error {
	return api_errors.NewErrorWithDetail(errs.NotFound, "Plan not found",
		api_errors.ErrorDetails{Code: api_errors.CodePlanNotFound})
}

func int32Ptr(v int32) *int32 { return &v }

// --- CreatePriceOffer ---

func TestCreatePriceOffer(t *testing.T) {
	const agentID int32 = 200001
	ctx := priceOfferAuthContext(agentID)
	q := testQuerier()

	pickupID, _ := seedPriceOfferLocation(t, q, "create-pickup")
	dropoffID, _ := seedPriceOfferLocation(t, q, "create-dropoff")
	plan := defaultPlan(pickupID, dropoffID)
	snapshotID := seedSnapshot(t, q, []planPriceDetails{plan})

	t.Run("rejects missing snapshot id", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.SnapshotID = 0
		api_errors.AssertApiError(t, priceOfferInvalidValueErr("snapshotId"), p.Validate())
	})

	t.Run("rejects missing rate qualifier", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.RateQualifier = ""
		api_errors.AssertApiError(t, priceOfferInvalidValueErr("rateQualifier"), p.Validate())
	})

	t.Run("rejects missing supplier code", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.SupplierCode = ""
		api_errors.AssertApiError(t, priceOfferInvalidValueErr("supplierCode"), p.Validate())
	})

	t.Run("accepts valid params", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns snapshot not found", func(t *testing.T) {
		p := validCreatePriceOfferParams(999999999, plan)
		_, err := CreatePriceOffer(ctx, p)
		api_errors.AssertApiError(t, snapshotNotFoundErr(), err)
	})

	t.Run("returns plan not found when rate qualifier doesn't match", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.RateQualifier = "DOES-NOT-EXIST"
		_, err := CreatePriceOffer(ctx, p)
		api_errors.AssertApiError(t, planNotFoundErr(), err)
	})

	t.Run("creates offer and persists all fields", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.Name = "Persisted Offer"
		p.OfferedCurrencyCode = "EUR"
		p.OfferedPrice = 250
		p.IncludeERP = true

		resp, err := CreatePriceOffer(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero offer id")
		}
		if resp.Token == "" {
			t.Fatal("expected non-empty token")
		}

		row, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
			ID:      resp.ID,
			AgentID: agentID,
		})
		if err != nil {
			t.Fatalf("failed to fetch created offer: %v", err)
		}

		if row.Name != "Persisted Offer" {
			t.Errorf("name: got %q, want %q", row.Name, "Persisted Offer")
		}
		if row.AgentID != agentID {
			t.Errorf("agent id: got %d, want %d", row.AgentID, agentID)
		}
		if row.SupplierCode != plan.SupplierCode {
			t.Errorf("supplier code: got %q, want %q", row.SupplierCode, plan.SupplierCode)
		}
		if row.CurrencyCode != plan.CurrencyCode {
			t.Errorf("currency code: got %q, want %q", row.CurrencyCode, plan.CurrencyCode)
		}
		if row.OfferedCurrencyCode != "EUR" {
			t.Errorf("offered currency: got %q, want %q", row.OfferedCurrencyCode, "EUR")
		}
		if row.OfferedPrice != 250 {
			t.Errorf("offered price: got %d, want %d", row.OfferedPrice, 250)
		}
		if string(row.Status) != "open" {
			t.Errorf("status: got %q, want %q", row.Status, "open")
		}
		if row.BtErpPrice != int32(plan.ChargedERPPriceWithVat) {
			t.Errorf("bt erp price: got %d, want %d", row.BtErpPrice, plan.ChargedERPPriceWithVat)
		}
		if db.NumericToFloat64(row.BrokerErpPrice) == 0 {
			t.Error("broker erp price should be non-zero when IncludeERP=true")
		}
		// Total = round((100*1.5 + 10*1.5) * 0.9) + 15 = round(148.5)+15 = 149+15 = 164
		if row.TotalPrice != 164 {
			t.Errorf("total price: got %d, want %d", row.TotalPrice, 164)
		}
		if row.PickupLocationID != pickupID {
			t.Errorf("pickup location id: got %q, want %q", row.PickupLocationID, pickupID)
		}
		if row.DropoffLocationID != dropoffID {
			t.Errorf("dropoff location id: got %q, want %q", row.DropoffLocationID, dropoffID)
		}
	})

	t.Run("excludes erp when IncludeERP=false", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.IncludeERP = false

		resp, err := CreatePriceOffer(ctx, p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		row, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
			ID:      resp.ID,
			AgentID: agentID,
		})
		if err != nil {
			t.Fatalf("failed to fetch created offer: %v", err)
		}
		if row.BtErpPrice != 0 {
			t.Errorf("bt erp price: got %d, want 0", row.BtErpPrice)
		}
		if db.NumericToFloat64(row.BrokerErpPrice) != 0 {
			t.Errorf("broker erp price: got %v, want 0", db.NumericToFloat64(row.BrokerErpPrice))
		}
		// Total = round(100*1.5 * 0.9) = 135.
		if row.TotalPrice != 135 {
			t.Errorf("total price: got %d, want 135", row.TotalPrice)
		}
	})
}

// --- GetClientPriceOffer ---

func TestGetClientPriceOffer(t *testing.T) {
	const agentID int32 = 200002
	ctx := priceOfferAuthContext(agentID)
	q := testQuerier()

	pickupID, pickupName := seedPriceOfferLocation(t, q, "client-pickup")
	dropoffID, dropoffName := seedPriceOfferLocation(t, q, "client-dropoff")
	plan := defaultPlan(pickupID, dropoffID)
	snapshotID := seedSnapshot(t, q, []planPriceDetails{plan})

	t.Run("returns 404 for invalid token", func(t *testing.T) {
		_, err := GetClientPriceOffer(context.Background(), "not-a-uuid")
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns 404 for non-existent token", func(t *testing.T) {
		_, err := GetClientPriceOffer(context.Background(), "00000000-0000-0000-0000-000000000000")
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	createParams := validCreatePriceOfferParams(snapshotID, plan)
	createParams.Name = "Client Offer"
	createParams.OfferedCurrencyCode = "EUR"
	createParams.OfferedPrice = 300
	created, err := CreatePriceOffer(ctx, createParams)
	if err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}

	t.Run("returns offer with public fields only", func(t *testing.T) {
		resp, err := GetClientPriceOffer(context.Background(), created.Token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID != created.ID {
			t.Errorf("id: got %d, want %d", resp.ID, created.ID)
		}
		if resp.Name != "Client Offer" {
			t.Errorf("name: got %q, want %q", resp.Name, "Client Offer")
		}
		if resp.Status != "open" {
			t.Errorf("status: got %q, want %q", resp.Status, "open")
		}
		if resp.PickupLocationName != pickupName {
			t.Errorf("pickup location: got %q, want %q", resp.PickupLocationName, pickupName)
		}
		if resp.DropoffLocationName != dropoffName {
			t.Errorf("dropoff location: got %q, want %q", resp.DropoffLocationName, dropoffName)
		}
		if resp.PickupDate != "2026-08-01" {
			t.Errorf("pickup date: got %q, want 2026-08-01", resp.PickupDate)
		}
		if resp.ReturnDate != "2026-08-05" {
			t.Errorf("return date: got %q, want 2026-08-05", resp.ReturnDate)
		}
		if resp.PickupTime != "08:00" {
			t.Errorf("pickup time: got %q, want 08:00", resp.PickupTime)
		}
		if resp.DropoffTime != "10:00" {
			t.Errorf("dropoff time: got %q, want 10:00", resp.DropoffTime)
		}
		if resp.DriverAge != "30" {
			t.Errorf("driver age: got %q, want 30", resp.DriverAge)
		}
		if resp.CurrencyCode != "EUR" {
			t.Errorf("currency code: got %q, want EUR", resp.CurrencyCode)
		}
		if resp.TotalPrice != 300 {
			t.Errorf("total price: got %d, want 300", resp.TotalPrice)
		}
		if !resp.IsErpIncluded {
			t.Error("expected IsErpIncluded=true")
		}
		if resp.CarDetails.Model != "Toyota Corolla" {
			t.Errorf("car model: got %q, want Toyota Corolla", resp.CarDetails.Model)
		}
		if len(resp.PlanInclusions) != 2 {
			t.Errorf("inclusions: got %d, want 2", len(resp.PlanInclusions))
		}
		if resp.CreatedAt == "" {
			t.Error("expected non-empty createdAt")
		}
	})
}

// --- GetAgentPriceOffer ---

func TestGetAgentPriceOffer(t *testing.T) {
	const agentID int32 = 200003
	ctx := priceOfferAuthContext(agentID)
	q := testQuerier()

	pickupID, pickupName := seedPriceOfferLocation(t, q, "agent-pickup")
	dropoffID, dropoffName := seedPriceOfferLocation(t, q, "agent-dropoff")
	plan := defaultPlan(pickupID, dropoffID)
	snapshotID := seedSnapshot(t, q, []planPriceDetails{plan})

	t.Run("returns 404 for non-existent id", func(t *testing.T) {
		_, err := GetAgentPriceOffer(ctx, 99999999)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	createParams := validCreatePriceOfferParams(snapshotID, plan)
	createParams.Name = "Agent Offer"
	createParams.IncludeERP = true
	createParams.OfferedCurrencyCode = "USD"
	createParams.OfferedPrice = 220
	created, err := CreatePriceOffer(ctx, createParams)
	if err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}

	t.Run("returns 404 for offer of another agent", func(t *testing.T) {
		otherCtx := priceOfferAuthContext(agentID + 1)
		_, err := GetAgentPriceOffer(otherCtx, created.ID)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns offer with internal pricing breakdown", func(t *testing.T) {
		resp, err := GetAgentPriceOffer(ctx, created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.ID != created.ID {
			t.Errorf("id: got %d, want %d", resp.ID, created.ID)
		}
		if resp.Token != created.Token {
			t.Errorf("token: got %q, want %q", resp.Token, created.Token)
		}
		if resp.Name != "Agent Offer" {
			t.Errorf("name: got %q, want %q", resp.Name, "Agent Offer")
		}
		if resp.Status != "open" {
			t.Errorf("status: got %q, want open", resp.Status)
		}
		if resp.SupplierCode != plan.SupplierCode {
			t.Errorf("supplier code: got %q, want %q", resp.SupplierCode, plan.SupplierCode)
		}
		if resp.CurrencyCode != plan.CurrencyCode {
			t.Errorf("currency code: got %q, want %q", resp.CurrencyCode, plan.CurrencyCode)
		}
		if resp.OfferedCurrencyCode != "USD" {
			t.Errorf("offered currency: got %q, want USD", resp.OfferedCurrencyCode)
		}
		if resp.OfferedPrice != 220 {
			t.Errorf("offered price: got %d, want 220", resp.OfferedPrice)
		}
		if resp.PickupLocationName != pickupName {
			t.Errorf("pickup location: got %q, want %q", resp.PickupLocationName, pickupName)
		}
		if resp.DropoffLocationName != dropoffName {
			t.Errorf("dropoff location: got %q, want %q", resp.DropoffLocationName, dropoffName)
		}

		// Pricing math:
		// car_full = round(100 * 1.5) = 150
		// erp_full = round(10 * 1.5 + 15) = 30
		// total   = round((150+15) * 0.9) + 15 = round(148.5)+15 = 149+15 = 164
		if resp.CarFullPrice != 150 {
			t.Errorf("car full price: got %d, want 150", resp.CarFullPrice)
		}
		if resp.ErpPrice != 30 {
			t.Errorf("erp price: got %d, want 30", resp.ErpPrice)
		}
		if resp.TotalPrice != 164 {
			t.Errorf("total price: got %d, want 164", resp.TotalPrice)
		}
		if resp.CarDetails.Model != "Toyota Corolla" {
			t.Errorf("car model: got %q, want Toyota Corolla", resp.CarDetails.Model)
		}
		if len(resp.PlanInclusions) != 2 {
			t.Errorf("inclusions: got %d, want 2", len(resp.PlanInclusions))
		}
		if resp.CreatedAt == "" {
			t.Error("expected non-empty createdAt")
		}
	})

	t.Run("zero erp prices when IncludeERP=false", func(t *testing.T) {
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.IncludeERP = false
		p.Name = "No ERP"
		c, err := CreatePriceOffer(ctx, p)
		if err != nil {
			t.Fatalf("failed to create offer: %v", err)
		}

		resp, err := GetAgentPriceOffer(ctx, c.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.CarFullPrice != 150 {
			t.Errorf("car full price: got %d, want 150", resp.CarFullPrice)
		}
		if resp.ErpPrice != 0 {
			t.Errorf("erp price: got %d, want 0", resp.ErpPrice)
		}
		// total = round(150 * 0.9) = 135
		if resp.TotalPrice != 135 {
			t.Errorf("total price: got %d, want 135", resp.TotalPrice)
		}
	})
}

// --- ListPriceOffers ---

func TestListPriceOffers(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		t.Run("rejects page 0", func(t *testing.T) {
			p := ListPriceOffersRequest{Page: 0}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("page"), p.Validate())
		})

		t.Run("rejects negative page", func(t *testing.T) {
			p := ListPriceOffersRequest{Page: -1}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("page"), p.Validate())
		})

		t.Run("accepts valid params", func(t *testing.T) {
			p := ListPriceOffersRequest{Page: 1}
			if err := p.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	})

	const agentID int32 = 200004
	ctx := priceOfferAuthContext(agentID)
	q := testQuerier()

	pickupID, _ := seedPriceOfferLocation(t, q, "list-pickup")
	dropoffID, _ := seedPriceOfferLocation(t, q, "list-dropoff")
	plan := defaultPlan(pickupID, dropoffID)
	snapshotID := seedSnapshot(t, q, []planPriceDetails{plan})

	createOffer := func(t *testing.T, name string) *PriceOfferResponse {
		t.Helper()
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.Name = name
		resp, err := CreatePriceOffer(ctx, p)
		if err != nil {
			t.Fatalf("failed to create offer %q: %v", name, err)
		}
		return resp
	}

	alice := createOffer(t, "Alice Offer")
	time.Sleep(10 * time.Millisecond) // ensure created_at differs
	bob := createOffer(t, "Bob Offer")
	time.Sleep(10 * time.Millisecond)
	charlie := createOffer(t, "Charlie Offer")

	t.Run("returns all offers for agent ordered by created_at DESC", func(t *testing.T) {
		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 3 {
			t.Fatalf("total: got %d, want 3", resp.Total)
		}
		if len(resp.PriceOffers) != 3 {
			t.Fatalf("len(priceOffers): got %d, want 3", len(resp.PriceOffers))
		}
		if resp.PriceOffers[0].ID != charlie.ID {
			t.Errorf("position 0 id: got %d, want %d (Charlie)", resp.PriceOffers[0].ID, charlie.ID)
		}
		if resp.PriceOffers[1].ID != bob.ID {
			t.Errorf("position 1 id: got %d, want %d (Bob)", resp.PriceOffers[1].ID, bob.ID)
		}
		if resp.PriceOffers[2].ID != alice.ID {
			t.Errorf("position 2 id: got %d, want %d (Alice)", resp.PriceOffers[2].ID, alice.ID)
		}
	})

	t.Run("summary fields populated correctly", func(t *testing.T) {
		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Name: "Alice", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.PriceOffers) != 1 {
			t.Fatalf("expected 1 offer, got %d", len(resp.PriceOffers))
		}
		got := resp.PriceOffers[0]
		if got.ID != alice.ID {
			t.Errorf("id: got %d, want %d", got.ID, alice.ID)
		}
		if got.Name != "Alice Offer" {
			t.Errorf("name: got %q, want Alice Offer", got.Name)
		}
		if got.Status != "open" {
			t.Errorf("status: got %q, want open", got.Status)
		}
		if got.PickupDate != "2026-08-01" {
			t.Errorf("pickup date: got %q, want 2026-08-01", got.PickupDate)
		}
		if got.ReturnDate != "2026-08-05" {
			t.Errorf("return date: got %q, want 2026-08-05", got.ReturnDate)
		}
		if got.CurrencyCode != "USD" {
			t.Errorf("currency code: got %q, want USD", got.CurrencyCode)
		}
		if got.OfferedPrice != 200 {
			t.Errorf("offered price: got %d, want 200", got.OfferedPrice)
		}
		if got.CreatedAt == "" {
			t.Error("expected non-empty createdAt")
		}
	})

	t.Run("filters by name (case-insensitive ILIKE)", func(t *testing.T) {
		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Name: "bob", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 1 {
			t.Fatalf("total: got %d, want 1", resp.Total)
		}
		if len(resp.PriceOffers) != 1 || resp.PriceOffers[0].ID != bob.ID {
			t.Fatalf("expected only Bob offer, got %+v", resp.PriceOffers)
		}
	})

	t.Run("filters by status", func(t *testing.T) {
		// Update Charlie to declined.
		if err := UpdatePriceOffer(ctx, charlie.ID, UpdatePriceOfferParams{
			Status: strPtr("declined"),
		}); err != nil {
			t.Fatalf("failed to update offer: %v", err)
		}

		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Status: "declined", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 1 {
			t.Fatalf("total: got %d, want 1", resp.Total)
		}
		if len(resp.PriceOffers) != 1 || resp.PriceOffers[0].ID != charlie.ID {
			t.Fatalf("expected only Charlie offer, got %+v", resp.PriceOffers)
		}
		if resp.PriceOffers[0].Status != "declined" {
			t.Errorf("status: got %q, want declined", resp.PriceOffers[0].Status)
		}

		// Open status filter excludes Charlie.
		resp, err = ListPriceOffers(ctx, ListPriceOffersRequest{Status: "open", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("total: got %d, want 2", resp.Total)
		}
		for _, o := range resp.PriceOffers {
			if o.ID == charlie.ID {
				t.Error("Charlie (declined) should not appear when filtering by open")
			}
		}
	})

	t.Run("paginates results", func(t *testing.T) {
		// Limit is 8, so all offers fit on a single page; page=2 should be empty.
		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.PriceOffers) != 0 {
			t.Errorf("expected 0 offers on page 2, got %d", len(resp.PriceOffers))
		}
		if resp.Total != 3 {
			t.Errorf("total should still reflect all matches; got %d, want 3", resp.Total)
		}
	})

	t.Run("non-matching filter returns empty with zero total", func(t *testing.T) {
		resp, err := ListPriceOffers(ctx, ListPriceOffersRequest{Name: "zzznonexistent", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 0 {
			t.Fatalf("total: got %d, want 0", resp.Total)
		}
		if len(resp.PriceOffers) != 0 {
			t.Errorf("len: got %d, want 0", len(resp.PriceOffers))
		}
	})

	t.Run("isolates offers by agent", func(t *testing.T) {
		otherCtx := priceOfferAuthContext(agentID + 100)
		resp, err := ListPriceOffers(otherCtx, ListPriceOffersRequest{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 0 {
			t.Errorf("expected 0 offers for other agent, got total=%d", resp.Total)
		}
	})
}

// --- UpdatePriceOffer ---

func TestUpdatePriceOffer(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		t.Run("rejects invalid status", func(t *testing.T) {
			p := UpdatePriceOfferParams{Status: strPtr("invalid")}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("status"), p.Validate())
		})

		t.Run("rejects blank name", func(t *testing.T) {
			p := UpdatePriceOfferParams{Name: strPtr("   ")}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("name"), p.Validate())
		})

		t.Run("rejects non-3-letter currency code", func(t *testing.T) {
			p := UpdatePriceOfferParams{OfferedCurrencyCode: strPtr("US")}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("offeredCurrencyCode"), p.Validate())
		})

		t.Run("rejects lowercase currency code", func(t *testing.T) {
			p := UpdatePriceOfferParams{OfferedCurrencyCode: strPtr("usd")}
			api_errors.AssertApiError(t, priceOfferInvalidValueErr("offeredCurrencyCode"), p.Validate())
		})

		t.Run("accepts empty params (all nil)", func(t *testing.T) {
			if err := (UpdatePriceOfferParams{}).Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})

		t.Run("accepts all valid fields", func(t *testing.T) {
			p := UpdatePriceOfferParams{
				Status:              strPtr("booked"),
				Name:                strPtr("Updated"),
				OfferedCurrencyCode: strPtr("EUR"),
				OfferedPrice:        int32Ptr(500),
			}
			if err := p.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	})

	const agentID int32 = 200005
	ctx := priceOfferAuthContext(agentID)
	q := testQuerier()

	pickupID, _ := seedPriceOfferLocation(t, q, "update-pickup")
	dropoffID, _ := seedPriceOfferLocation(t, q, "update-dropoff")
	plan := defaultPlan(pickupID, dropoffID)
	snapshotID := seedSnapshot(t, q, []planPriceDetails{plan})

	createOffer := func(t *testing.T, name string) *PriceOfferResponse {
		t.Helper()
		p := validCreatePriceOfferParams(snapshotID, plan)
		p.Name = name
		resp, err := CreatePriceOffer(ctx, p)
		if err != nil {
			t.Fatalf("failed to create offer: %v", err)
		}
		return resp
	}

	t.Run("succeeds for non-existent id (no-op)", func(t *testing.T) {
		// :exec returns no error when zero rows are affected.
		err := UpdatePriceOffer(ctx, 99999999, UpdatePriceOfferParams{Status: strPtr("booked")})
		if err != nil {
			t.Fatalf("expected no error for non-existent id, got %v", err)
		}
	})

	t.Run("updates all provided fields", func(t *testing.T) {
		offer := createOffer(t, "Original")

		err := UpdatePriceOffer(ctx, offer.ID, UpdatePriceOfferParams{
			Status:              strPtr("booked"),
			Name:                strPtr("Updated Name"),
			OfferedCurrencyCode: strPtr("EUR"),
			OfferedPrice:        int32Ptr(999),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		row, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
			ID:      offer.ID,
			AgentID: agentID,
		})
		if err != nil {
			t.Fatalf("failed to fetch offer: %v", err)
		}
		if string(row.Status) != "booked" {
			t.Errorf("status: got %q, want booked", row.Status)
		}
		if row.Name != "Updated Name" {
			t.Errorf("name: got %q, want Updated Name", row.Name)
		}
		if row.OfferedCurrencyCode != "EUR" {
			t.Errorf("offered currency: got %q, want EUR", row.OfferedCurrencyCode)
		}
		if row.OfferedPrice != 999 {
			t.Errorf("offered price: got %d, want 999", row.OfferedPrice)
		}
	})

	t.Run("only updates provided fields (nil leaves existing)", func(t *testing.T) {
		offer := createOffer(t, "Partial Original")

		// Update only the status.
		if err := UpdatePriceOffer(ctx, offer.ID, UpdatePriceOfferParams{
			Status: strPtr("declined"),
		}); err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		row, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
			ID:      offer.ID,
			AgentID: agentID,
		})
		if err != nil {
			t.Fatalf("failed to fetch offer: %v", err)
		}
		if string(row.Status) != "declined" {
			t.Errorf("status: got %q, want declined", row.Status)
		}
		if row.Name != "Partial Original" {
			t.Errorf("name should be unchanged: got %q", row.Name)
		}
		if row.OfferedCurrencyCode != "USD" {
			t.Errorf("offered currency should be unchanged: got %q", row.OfferedCurrencyCode)
		}
		if row.OfferedPrice != 200 {
			t.Errorf("offered price should be unchanged: got %d", row.OfferedPrice)
		}
	})

	t.Run("does not update offers of other agents", func(t *testing.T) {
		offer := createOffer(t, "Mine")

		otherCtx := priceOfferAuthContext(agentID + 1)
		if err := UpdatePriceOffer(otherCtx, offer.ID, UpdatePriceOfferParams{
			Name: strPtr("Hijacked"),
		}); err != nil {
			t.Fatalf("expected no error (no-op), got %v", err)
		}

		row, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
			ID:      offer.ID,
			AgentID: agentID,
		})
		if err != nil {
			t.Fatalf("failed to fetch offer: %v", err)
		}
		if row.Name != "Mine" {
			t.Errorf("name should be unchanged by other agent: got %q, want Mine", row.Name)
		}
	})
}
