package availability

import (
	"testing"

	"encore.app/internal/broker"
)

func TestDropDuplicateCars(t *testing.T) {
	car := func(id int, b broker.Name, supplier, acriss, model, currency, rq string, price int) AvailableVehicle {
		return AvailableVehicle{
			ID:           id,
			Broker:       b,
			CarDetails:   broker.CarDetails{SupplierName: supplier, Acriss: acriss, Model: model},
			PriceDetails: broker.PriceDetails{Currency: currency},
			Plans:        []Plan{{PlanID: 1, Price: price, RateQualifier: rq, SupplierCode: supplier}},
		}
	}
	plan := func(v AvailableVehicle) PlanPriceDetails {
		return PlanPriceDetails{Broker: v.Broker, PlanID: 1, RateQualifier: v.Plans[0].RateQualifier, SupplierCode: v.Plans[0].SupplierCode}
	}
	ids := func(a availabilityArtifacts) []int {
		var out []int
		for _, v := range a.availableCars {
			out = append(out, v.ID)
		}
		return out
	}
	rates := map[string]float64{"EUR": 4, "USD": 3.5}

	t.Run("keeps the cheaper of the same car from two brokers", func(t *testing.T) {
		flex := car(0, broker.BrokerFlex, "Hertz AutoHellas", "MDMR", "Hyundai i10 or similar", "EUR", "flex-rq", 100)
		ours := car(1, broker.BrokerHertzGR, "Hertz", "MDMR", "HYUNDAI  I10", "EUR", "ITRH296:A", 95)
		got := dropDuplicateCars(availabilityArtifacts{
			availableCars: []AvailableVehicle{flex, ours},
			plansDetails:  []PlanPriceDetails{plan(flex), plan(ours)},
		}, rates)

		if want := []int{1}; len(ids(got)) != 1 || ids(got)[0] != want[0] {
			t.Errorf("cars = %v, want %v", ids(got), want)
		}
		if len(got.plansDetails) != 1 || got.plansDetails[0].Broker != broker.BrokerHertzGR {
			t.Errorf("snapshot plans = %+v, want only the kept car's", got.plansDetails)
		}
	})

	t.Run("compares prices across currencies", func(t *testing.T) {
		eur := car(0, broker.BrokerFlex, "Hertz", "MDMR", "i10", "EUR", "a", 100)    // 400
		usd := car(1, broker.BrokerHertzGR, "Hertz", "MDMR", "i10", "USD", "b", 110) // 385
		got := dropDuplicateCars(availabilityArtifacts{availableCars: []AvailableVehicle{eur, usd}}, rates)
		if len(ids(got)) != 1 || ids(got)[0] != 1 {
			t.Errorf("cars = %v, want the dollar car, cheaper once converted", ids(got))
		}
	})

	t.Run("keeps cars that differ in brand, ACRISS or model", func(t *testing.T) {
		got := dropDuplicateCars(availabilityArtifacts{availableCars: []AvailableVehicle{
			car(0, broker.BrokerFlex, "Hertz", "MDMR", "i10", "EUR", "a", 100),
			car(1, broker.BrokerHertzGR, "Thrifty", "MDMR", "i10", "EUR", "b", 90),
			car(2, broker.BrokerHertzGR, "Hertz", "EDMR", "i10", "EUR", "c", 90),
			car(3, broker.BrokerHertzGR, "Hertz", "MDMR", "Picanto", "EUR", "d", 90),
		}}, rates)
		if len(ids(got)) != 4 {
			t.Errorf("cars = %v, want all four", ids(got))
		}
	})

	t.Run("never compares cars of the same broker", func(t *testing.T) {
		got := dropDuplicateCars(availabilityArtifacts{availableCars: []AvailableVehicle{
			car(0, broker.BrokerFlex, "Hertz", "MDMR", "i10", "EUR", "a", 100),
			car(1, broker.BrokerFlex, "Hertz", "MDMR", "i10", "EUR", "b", 90),
		}}, rates)
		if len(ids(got)) != 2 {
			t.Errorf("cars = %v, want both", ids(got))
		}
	})
}
