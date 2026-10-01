package price_offer

import (
	"encoding/json"
	"testing"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/availability"
)

func TestFindRenewalPlanMatchesPlanID(t *testing.T) {
	car := broker.CarDetails{Model: "Toyota Yaris", SupplierName: "Avance", Acriss: "EDMR"}

	// Avance lists Standard before Gold, and both share the car.
	plansJSON, err := json.Marshal([]availability.PlanPriceDetails{
		{PlanID: 1, RateQualifier: "vch-quote:E-GQ", SupplierCode: "AVANCE", CarDetails: car, CarPurchasePrice: 77},
		{PlanID: 2, RateQualifier: "fdw0-quote:E-GQ", SupplierCode: "AVANCE", CarDetails: car, CarPurchasePrice: 201},
	})
	if err != nil {
		t.Fatal(err)
	}
	carJSON, err := json.Marshal(car)
	if err != nil {
		t.Fatal(err)
	}

	offer := db.GetPriceOfferByIdRow{PlanID: "2", SupplierCode: "AVANCE", CarDetails: carJSON}
	plan, err := findRenewalPlan(db.AvailablePlansSnapshot{Plans: plansJSON}, offer)
	if err != nil {
		t.Fatalf("findRenewalPlan: %v", err)
	}
	if plan.PlanID != 2 || plan.RateQualifier != "fdw0-quote:E-GQ" {
		t.Errorf("plan = %d %s, want the Gold plan", plan.PlanID, plan.RateQualifier)
	}
}
