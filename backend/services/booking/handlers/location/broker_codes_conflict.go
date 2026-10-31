package location

import (
	"slices"

	"encore.app/services/booking/db"
)

// brokerCodesConflict reports whether two stations cannot belong to the same location. A broker
// without brands holds one station per location. A broker with brands holds one per brand, since a
// brand may work from its own desk at the same airport, so its stations conflict only when they
// share a brand.
func brokerCodesConflict(brokerA db.Broker, brandsA []string, brokerB db.Broker, brandsB []string) bool {
	if brokerA != brokerB {
		return false
	}
	if len(brandsA) == 0 || len(brandsB) == 0 {
		return true
	}

	return slices.ContainsFunc(brandsA, func(b string) bool { return slices.Contains(brandsB, b) })
}

// stationConflict returns the code of another station among a location's codes that the given
// station cannot sit beside, or "" when there is none.
func stationConflict(codes []db.GetAllLocationBrokerCodesByLocationIDsRow, dbBroker db.Broker, brokerLocationID string, brands []string) string {
	for _, c := range codes {
		if c.BrokerLocationID == brokerLocationID && c.Broker == dbBroker {
			continue
		}
		if brokerCodesConflict(c.Broker, c.Brands, dbBroker, brands) {
			return c.BrokerLocationID
		}
	}

	return ""
}

// brandsOrEmpty keeps a station's brands a JSON array for the admin client.
func brandsOrEmpty(brands []string) []string {
	if brands == nil {
		return []string{}
	}

	return brands
}
