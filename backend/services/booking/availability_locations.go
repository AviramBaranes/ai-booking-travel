package booking

import (
	"context"

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

// availabilityLocationQuery holds broker-specific location IDs and country codes for a single availability search.
type availabilityLocationQuery struct {
	pickupCountryCode       string
	pickupBrokerLocationID  string
	dropoffBrokerLocationID string
}

// availabilityLocations maps each broker to the location query details needed to call its API.
type availabilityLocations map[broker.Name]availabilityLocationQuery

// getLocations resolves pickup and dropoff broker location IDs from the database and returns them grouped by broker.
func getLocations(ctx context.Context, query db.Querier, params SearchAvailabilityRequest) (availabilityLocations, error) {
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

// brokerLocation pairs a broker-specific location ID with the location's country code.
type brokerLocation struct {
	locationID  string
	countryCode string
}

// createBrokersMap splits location rows into pickup and dropoff maps keyed by broker.
func createBrokersMap(locs []db.GetAllLocationBrokerCodesByLocationIDsRow, params SearchAvailabilityRequest, dropoffLocationRowID int64) (map[db.Broker]brokerLocation, map[db.Broker]brokerLocation) {
	pickupsByBroker := make(map[db.Broker]brokerLocation)
	dropoffsByBroker := make(map[db.Broker]brokerLocation)

	for _, loc := range locs {
		if !loc.Enabled {
			continue
		}

		if loc.LocationID == params.PickupLocationID {
			pickupsByBroker[loc.Broker] = brokerLocation{
				locationID:  loc.BrokerLocationID,
				countryCode: loc.LocationCountryCode,
			}
		}

		if loc.LocationID == dropoffLocationRowID {
			dropoffsByBroker[loc.Broker] = brokerLocation{
				locationID:  loc.BrokerLocationID,
				countryCode: loc.LocationCountryCode,
			}
		}
	}

	return pickupsByBroker, dropoffsByBroker
}

// buildAvailabilityLocations returns only the brokers that have both a pickup and a dropoff location.
func buildAvailabilityLocations(pickupsByBroker, dropoffsByBroker map[db.Broker]brokerLocation) availabilityLocations {
	al := make(availabilityLocations)

	for brokerName, pickupID := range pickupsByBroker {
		dropoffID, ok := dropoffsByBroker[brokerName]
		if !ok {
			rlog.Info(
				"removing broker from availability search because it does not have both pickup and dropoff locations",
				"broker", brokerName,
			)
			continue
		}

		al[broker.Name(brokerName)] = availabilityLocationQuery{
			pickupBrokerLocationID:  pickupID.locationID,
			dropoffBrokerLocationID: dropoffID.locationID,
			pickupCountryCode:       pickupID.countryCode,
		}
	}

	return al
}
