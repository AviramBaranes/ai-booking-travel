package booking

import (
	"context"
	"sort"
	"sync"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// SearchAvailabilityRequest represents the request for searching availability of vehicles.
type SearchAvailabilityRequest struct {
	PickupLocationID  int64  `query:"pickupLocationId" validate:"required"`
	DropoffLocationID int64  `query:"dropoffLocationId" validate:"omitempty"`
	PickupTime        string `query:"pickupTime" validate:"required,datetime=15:04"`
	DropoffTime       string `query:"dropoffTime" validate:"required,datetime=15:04"`
	PickupDate        string `query:"pickupDate" validate:"required,datetime=2006-01-02"`
	DropoffDate       string `query:"dropoffDate" validate:"required,datetime=2006-01-02"`
	DriverAge         int    `query:"driverAge" validate:"required,gte=18"`
	CouponCode        string `query:"couponCode" validate:"omitempty"`
}

func (p SearchAvailabilityRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// SearchAvailabilityResponse represents the response for searching availability of vehicles.
type SearchAvailabilityResponse struct {
	AvailableVehicles []AvailableVehicle `json:"availableVehicles"`
}

// SearchAvailability handles the http request for searching availability of vehicles.
// encore:api public method=GET path=/booking/availability
func (s *Service) SearchAvailability(ctx context.Context, params SearchAvailabilityRequest) (*SearchAvailabilityResponse, error) {
	availableLocs, err := getLocations(ctx, s.query, params)
	if err != nil {
		return nil, err
	}

	vs := make([]broker.AvailableVehicle, 0)
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errOnce  sync.Once
		firstErr error
	)
	for bn, loc := range availableLocs {
		wg.Add(1)
		go func(loc availabilityLocationQuery, bn broker.Name) {
			defer wg.Done()

			b, err := getBrokerByName(bn)
			if err != nil {
				rlog.Error("failed to get broker by name", "error", err, "broker", bn)
				errOnce.Do(func() {
					firstErr = err
				})
				return
			}

			result, err := searchCars(b, params, loc.pickupBrokerLocationID, loc.dropoffBrokerLocationID, loc.pickupCountryCode)
			if err != nil {
				errOnce.Do(func() {
					firstErr = err
				})
				return
			}

			mu.Lock()
			vs = append(vs, result...)
			mu.Unlock()
		}(loc, bn)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Plans[0].Price < vs[j].Plans[0].Price
	})

	return &SearchAvailabilityResponse{
		AvailableVehicles: []AvailableVehicle{},
	}, nil
}

func searchCars(b broker.Broker, params SearchAvailabilityRequest, plID, dlID, countryCode string) ([]broker.AvailableVehicle, error) {
	vs, err := b.SearchAvailability(broker.SearchAvailabilityParams{
		CountryCode:        countryCode,
		PickupLocation:     plID,
		DropoffLocation:    dlID,
		PickupTime:         params.PickupTime,
		DropoffTime:        params.DropoffTime,
		PickupDate:         params.PickupDate,
		DropoffDate:        params.DropoffDate,
		DriverAge:          params.DriverAge,
		DiscountPercentage: 0, // TODO: apply coupon code to get discount percentage
	})
	if err != nil {
		rlog.Error("failed to search availability", "error", err, "broker", b.Name(), "pickupLocation", plID, "dropoffLocation", dlID)
		return nil, api_errors.ErrInternalError
	}

	return vs, nil
}

func getBrokerByName(name broker.Name) (broker.Broker, error) {
	switch name {
	case broker.BrokerFlex:
		f := broker.NewFlex()
		return &f, nil
	case broker.BrokerHertz:
		return broker.NewHertz(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}

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

type availabilityLocationQuery struct {
	pickupCountryCode       string
	dropoffCountryCode      string
	pickupBrokerLocationID  string
	dropoffBrokerLocationID string
}

type availabilityLocations map[broker.Name]availabilityLocationQuery

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

type brokerLocation struct {
	locationID  string
	countryCode string
}

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
			dropoffCountryCode:      dropoffID.countryCode,
		}
	}

	return al
}

type planPriceDetails struct {
	planID                string
	rateQualifier         string
	supplierCode          string
	originalPrice         float64
	priceWithMarkup       float64
	priceAfterDiscount    float64
	erpOriginalPrice      float64
	erpPriceWithMarkup    float64
	erpPriceWithVAT       float64
	erpPriceAfterDiscount float64
}
