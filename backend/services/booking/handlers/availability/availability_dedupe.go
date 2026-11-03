package availability

import (
	"strings"

	"encore.app/internal/broker"
)

// brandAliases maps a supplier name onto the brand it sells, where a broker names the supplier
// differently from the brand: Flex sells Autohellas, the Hertz franchisee in Greece, as "Hertz
// AutoHellas".
var brandAliases = map[string]string{
	"hertz autohellas": "hertz",
}

// carIdentity is what makes two brokers' cars the same car: the brand, the ACRISS code and the
// model.
type carIdentity struct {
	brand  string
	acriss string
	model  string
}

func identityOf(v AvailableVehicle) carIdentity {
	brand := strings.ToLower(strings.TrimSpace(v.CarDetails.SupplierName))
	if alias, ok := brandAliases[brand]; ok {
		brand = alias
	}

	model := strings.ToLower(strings.Join(strings.Fields(v.CarDetails.Model), " "))
	model = strings.TrimSuffix(model, " or similar")

	return carIdentity{
		brand:  brand,
		acriss: strings.ToUpper(strings.TrimSpace(v.CarDetails.Acriss)),
		model:  model,
	}
}

// dropDuplicateCars keeps the cheaper of two cars that different brokers sell as the same car (same
// brand, ACRISS and model), so Flex's Autohellas cars and Hertz GR's own don't both show. Cars of
// one broker are never compared. Prices are compared by each car's cheapest plan, converted with
// the search's currency rates. The dropped cars' plans are removed from the snapshot too.
func dropDuplicateCars(a availabilityArtifacts, rates map[string]float64) availabilityArtifacts {
	cheapest := make(map[carIdentity]int) // identity -> index in a.availableCars of the car kept so far
	drop := make(map[int]bool)

	for i, v := range a.availableCars {
		if len(v.Plans) == 0 {
			continue
		}

		id := identityOf(v)
		j, seen := cheapest[id]
		if !seen {
			cheapest[id] = i
			continue
		}

		kept := a.availableCars[j]
		if kept.Broker == v.Broker {
			// The same broker selling a car twice is two offers, not a duplicate.
			continue
		}

		if lowestPrice(v, rates) < lowestPrice(kept, rates) {
			drop[j] = true
			cheapest[id] = i
		} else {
			drop[i] = true
		}
	}

	if len(drop) == 0 {
		return a
	}

	dropped := make(map[planKey]bool)
	cars := make([]AvailableVehicle, 0, len(a.availableCars)-len(drop))
	for i, v := range a.availableCars {
		if !drop[i] {
			cars = append(cars, v)
			continue
		}
		for _, p := range v.Plans {
			dropped[planKey{v.Broker, p.RateQualifier, p.SupplierCode, p.PlanID}] = true
		}
	}

	plans := make([]PlanPriceDetails, 0, len(a.plansDetails))
	for _, pd := range a.plansDetails {
		if !dropped[planKey{pd.Broker, pd.RateQualifier, pd.SupplierCode, pd.PlanID}] {
			plans = append(plans, pd)
		}
	}

	return availabilityArtifacts{availableCars: cars, plansDetails: plans}
}

// planKey identifies a plan within a snapshot, the way booking finds it, plus the broker.
type planKey struct {
	broker        broker.Name
	rateQualifier string
	supplierCode  string
	planID        int
}

// lowestPrice returns the car's cheapest plan price in a common currency.
func lowestPrice(v AvailableVehicle, rates map[string]float64) float64 {
	lowest := v.Plans[0].Price
	for _, p := range v.Plans[1:] {
		lowest = min(lowest, p.Price)
	}
	return float64(lowest) * rates[v.PriceDetails.Currency]
}
