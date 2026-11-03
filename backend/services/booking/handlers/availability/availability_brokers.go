package availability

import (
	"sync"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.dev/rlog"
)

func searchAvailabilityAcrossBrokers(cfg *AvailableVehiclesConfig, p SearchAvailabilityParams, locs availabilityLocations, isAgent bool) (*broker.AvailabilityResponse, error) {
	vs := make([]broker.AvailableVehicle, 0)
	si := make([]broker.SupplierInfo, 0)
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

			b, err := getBrokerByName(cfg, bn)
			if err != nil {
				rlog.Error("failed to get broker by name", "error", err, "broker", bn)
				errOnce.Do(func() {
					firstErr = err
				})
				return
			}

			result, err := searchCars(b, p, loc, isAgent)
			if err != nil {
				errOnce.Do(func() {
					firstErr = err
				})
				return
			}

			// Two brokers can sell a supplier of the same name, so supplier info is told apart by
			// broker.
			for i := range result.SuppliersInfo {
				result.SuppliersInfo[i].Broker = bn
			}

			mu.Lock()
			vs = append(vs, result.AvailableVehicles...)
			si = append(si, result.SuppliersInfo...)
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

	return &broker.AvailabilityResponse{
		AvailableVehicles: vs,
		SuppliersInfo:     si,
	}, nil
}

// searchCars calls the given broker to search for available vehicles with the supplied location and date parameters.
func searchCars(b broker.AvailabilitySearcher, p SearchAvailabilityParams, loc availabilityLocationQuery, isAgent bool) (*broker.AvailabilityResponse, error) {
	plID, dlID := loc.pickupBrokerLocationID, loc.dropoffBrokerLocationID
	vs, err := b.SearchAvailability(broker.SearchAvailabilityParams{
		CountryCode:          loc.pickupCountryCode,
		PickupLocation:       plID,
		DropoffLocation:      dlID,
		PickupStations:       loc.pickupStations,
		DropoffStations:      loc.dropoffStations,
		PickupTime:           p.PickupTime,
		DropoffTime:          p.DropoffTime,
		PickupDate:           p.PickupDate,
		DropoffDate:          p.DropoffDate,
		DriverAge:            p.DriverAge,
		IsRestrictionApplied: !isAgent,
	})
	if err != nil {
		rlog.Error("failed to search availability", "error", err, "broker", b.Name(), "pickupLocation", plID, "dropoffLocation", dlID)
		return nil, api_errors.ErrInternalError
	}

	return vs, nil
}

// getBrokerByName returns an initialised Broker for the given broker name.
func getBrokerByName(cfg *AvailableVehiclesConfig, name broker.Name) (broker.AvailabilitySearcher, error) {
	switch name {
	case broker.BrokerFlex:
		return broker.NewFlexWithErpCfg(cfg.FlexErpDayCharge()), nil
	case broker.BrokerHertz:
		return broker.NewHertzWithCharges(cfg.HertzErpDayChargeUS(), cfg.HertzErpDayChargeCA()), nil
	case broker.BrokerAvance:
		return broker.NewAvance(broker.AvanceConfig{
			BaseURL:                cfg.AvanceBaseURL(),
			StandardErpDayCharge:   cfg.AvanceErpDayChargeStandard(),
			ZeroExcessErpDayCharge: cfg.AvanceErpDayChargeZeroExcess(),
		}), nil
	case broker.BrokerHertzGR:
		return broker.NewHertzGR(broker.HertzGRConfig{
			BaseURL: cfg.HertzGRBaseURL(),
			Rates:   hertzGRRates(cfg),
			LogXML:  cfg.HertzGRLogXML(),
		}), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}

// hertzGRRates converts the configured rate codes into the broker's own type, which keeps
// internal/broker free of Encore config.
func hertzGRRates(cfg *AvailableVehiclesConfig) []broker.HertzGRRate {
	configured := cfg.HertzGRRates()
	rates := make([]broker.HertzGRRate, len(configured))
	for i, r := range configured {
		rates[i] = broker.HertzGRRate{Code: r.Code, PlanName: r.PlanName, ErpDayCharge: r.ErpDayCharge}
	}
	return rates
}
