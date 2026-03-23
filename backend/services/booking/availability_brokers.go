package booking

import (
	"sync"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.dev/rlog"
)

func searchAvailabilityAcrossBrokers(p SearchAvailabilityRequest, locs availabilityLocations) ([]broker.AvailableVehicle, error) {
	vs := make([]broker.AvailableVehicle, 0)
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errOnce  sync.Once
		firstErr error
	)
	for bn, loc := range locs {
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

			result, err := searchCars(b, p, loc.pickupBrokerLocationID, loc.dropoffBrokerLocationID, loc.pickupCountryCode)
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
		if len(vs) == 0 {
			rlog.Error("search availability across brokers failed", "error", firstErr)
			return nil, firstErr
		}
		rlog.Warn("search availability across brokers completed with partial results due to some errors", "error", firstErr)
	}

	return vs, nil
}

// searchCars calls the given broker to search for available vehicles with the supplied location and date parameters.
func searchCars(b broker.AvailabilitySearcher, params SearchAvailabilityRequest, plID, dlID, countryCode string) ([]broker.AvailableVehicle, error) {
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

// getBrokerByName returns an initialised Broker for the given broker name.
func getBrokerByName(name broker.Name) (broker.AvailabilitySearcher, error) {
	switch name {
	case broker.BrokerFlex:
		f := broker.NewFlexWithErpCfg(avCfg.FlexErpDayCharge())
		return &f, nil
	case broker.BrokerHertz:
		return broker.NewHertz(avCfg.HertzErpDayChargeUS(), avCfg.HertzErpDayChargeCA()), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}
