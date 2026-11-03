package availability

import (
	"context"
	"sort"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// Sentinel errors returned when location broker codes cannot be resolved.
var (
	errNoBrokerCodeForPickupLocation = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "no broker code found for pickup location", api_errors.ErrorDetails{
		Code:  api_errors.CodeNoBrokerCodeForPickupLocation,
		Field: "pickupLocationId",
	})
	errNoBrokerCodeForDropoffLocation = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "no broker code found for dropoff location", api_errors.ErrorDetails{
		Code:  api_errors.CodeNoBrokerCodeForDropoffLocation,
		Field: "dropoffLocationId",
	})
	errNoCommonBrokerForLocations = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "no common broker found for pickup and dropoff locations", api_errors.ErrorDetails{
		Code: api_errors.CodeNoCommonBrokerForLocations,
	})
)

// availabilityLocationQuery holds broker-specific location IDs, country codes, and canonical location names for a single availability search.
type availabilityLocationQuery struct {
	pickupCountryCode       string
	pickupBrokerLocationID  string
	dropoffBrokerLocationID string
	// pickupStations and dropoffStations are every station the broker has at each location; the
	// broker location IDs above are the first of each. A broker selling several brands may have
	// one station per brand.
	pickupStations      []broker.Station
	dropoffStations     []broker.Station
	pickupLocationName  string
	dropoffLocationName string
}

// availabilityLocations maps each broker to the location query details needed to call its API.
type availabilityLocations map[broker.Name]availabilityLocationQuery

// getLocations resolves pickup and dropoff broker location IDs from the database and returns them grouped by broker.
func getLocations(ctx context.Context, query db.Querier, params SearchAvailabilityParams) (availabilityLocations, error) {
	dropoffLocationRowID := params.DropoffLocationID
	if dropoffLocationRowID == 0 {
		dropoffLocationRowID = params.PickupLocationID
	}

	locs, err := query.GetAllLocationBrokerCodesByLocationIDs(ctx, []int64{params.PickupLocationID, dropoffLocationRowID})
	if err != nil {
		rlog.Error("failed to get broker codes for locations",
			"error", err,
			"pickupLocationId", params.PickupLocationID,
			"dropoffLocationId", dropoffLocationRowID,
		)
		return availabilityLocations{}, api_errors.ErrInternalError
	}

	pickupsByBroker, dropoffsByBroker := createBrokersMap(locs, params, dropoffLocationRowID)
	al := buildAvailabilityLocations(pickupsByBroker, dropoffsByBroker)

	if len(al) == 0 {
		switch {
		case len(pickupsByBroker) == 0:
			return nil, errNoBrokerCodeForPickupLocation
		case len(dropoffsByBroker) == 0:
			return nil, errNoBrokerCodeForDropoffLocation
		default:
			return nil, errNoCommonBrokerForLocations
		}
	}

	return al, nil
}

// brokerLocation holds a broker's stations at a location with the location's country code and canonical name.
type brokerLocation struct {
	stations    []broker.Station
	countryCode string
	name        string
}

// createBrokersMap splits location rows into pickup and dropoff maps keyed by broker, collecting
// every enabled station a broker has at each location, ordered by code.
func createBrokersMap(locs []db.GetAllLocationBrokerCodesByLocationIDsRow, params SearchAvailabilityParams, dropoffLocationRowID int64) (map[db.Broker]brokerLocation, map[db.Broker]brokerLocation) {
	pickupsByBroker := make(map[db.Broker]brokerLocation)
	dropoffsByBroker := make(map[db.Broker]brokerLocation)

	add := func(m map[db.Broker]brokerLocation, loc db.GetAllLocationBrokerCodesByLocationIDsRow) {
		bl := m[loc.Broker]
		bl.stations = append(bl.stations, broker.Station{Code: loc.BrokerLocationID, Brands: loc.Brands})
		bl.countryCode = loc.LocationCountryCode
		bl.name = loc.LocationName
		m[loc.Broker] = bl
	}

	for _, loc := range locs {
		if !loc.Enabled {
			continue
		}

		if loc.LocationID == params.PickupLocationID {
			add(pickupsByBroker, loc)
		}

		if loc.LocationID == dropoffLocationRowID {
			add(dropoffsByBroker, loc)
		}
	}

	for _, m := range []map[db.Broker]brokerLocation{pickupsByBroker, dropoffsByBroker} {
		for _, bl := range m {
			sort.Slice(bl.stations, func(i, j int) bool { return bl.stations[i].Code < bl.stations[j].Code })
		}
	}

	return pickupsByBroker, dropoffsByBroker
}

// buildAvailabilityLocations returns only the brokers that have both a pickup and a dropoff location.
func buildAvailabilityLocations(pickupsByBroker, dropoffsByBroker map[db.Broker]brokerLocation) availabilityLocations {
	al := make(availabilityLocations)

	for brokerName, pickupBrokerLoc := range pickupsByBroker {
		dropoffBrokerLoc, ok := dropoffsByBroker[brokerName]
		if !ok {
			rlog.Info(
				"removing broker from availability search because it does not have both pickup and dropoff locations",
				"broker", brokerName,
			)
			continue
		}

		al[broker.Name(brokerName)] = availabilityLocationQuery{
			pickupBrokerLocationID:  pickupBrokerLoc.stations[0].Code,
			dropoffBrokerLocationID: dropoffBrokerLoc.stations[0].Code,
			pickupStations:          pickupBrokerLoc.stations,
			dropoffStations:         dropoffBrokerLoc.stations,
			pickupCountryCode:       pickupBrokerLoc.countryCode,
			pickupLocationName:      pickupBrokerLoc.name,
			dropoffLocationName:     dropoffBrokerLoc.name,
		}
	}

	return al
}

// extractCountryCode extracts a country code from the availability locations of the first broker in the map. Country code is the same for all locations in the map as it defined in the row of the location and not location broker.
func extractCountryCode(al availabilityLocations) string {
	for _, loc := range al {
		return loc.pickupCountryCode
	}

	return ""
}

// extractPickupLocationName returns the canonical pickup location name from the first broker in the map.
func extractPickupLocationName(al availabilityLocations) string {
	for _, loc := range al {
		return loc.pickupLocationName
	}

	return ""
}

// extractDropoffLocationName returns the canonical dropoff location name from the first broker in the map.
func extractDropoffLocationName(al availabilityLocations) string {
	for _, loc := range al {
		return loc.dropoffLocationName
	}

	return ""
}
