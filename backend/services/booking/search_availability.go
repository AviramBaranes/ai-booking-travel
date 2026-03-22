package booking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.app/services/auth"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/config"
	"encore.dev/rlog"
)

type AvailableVehiclesConfig struct {
	HertzErpDayChargeUS config.Int
	HertzErpDayChargeCA config.Int
	FlexErpDayCharge    config.Int
	MarkUpGross         config.Float64
	MarkUpNet           config.Float64
}

var avCfg = config.Load[*AvailableVehiclesConfig]()

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
	planDetailsID     int64
	AvailableVehicles []AvailableVehicle `json:"availableVehicles"`
}

// SearchAvailability handles the http request for searching availability of vehicles.
// encore:api public method=GET path=/booking/availability
func (s *Service) SearchAvailability(ctx context.Context, p SearchAvailabilityRequest) (*SearchAvailabilityResponse, error) {
	availableLocs, err := getLocations(ctx, s.query, p)
	if err != nil {
		return nil, err
	}

	coupon, err := s.query.FindCouponByCode(ctx, p.CouponCode)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to find coupon by code", "error", err, "code", p.CouponCode)
		return nil, api_errors.ErrInternalError
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
		return nil, firstErr
	}

	if len(vs) == 0 {
		return &SearchAvailabilityResponse{AvailableVehicles: []AvailableVehicle{}}, nil
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Plans[0].Price < vs[j].Plans[0].Price
	})

	daysCount, err := broker.CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	mp, err := getMarkupProviderMap(ctx, availableLocs, s.query, daysCount, p.PickupDate, extractCarGroups(vs))
	if err != nil {
		rlog.Error("failed to get markup provider map", "error", err)
		return nil, api_errors.ErrInternalError
	}

	authData := auth.GetAuthData()
	isAgent := authData != nil && (authData.Role == auth.UserRoleAgent)
	currenciesMap, err := buildCurrencyMap(s.query)
	if err != nil {
		rlog.Error("failed to build currency map", "error", err)
		return nil, api_errors.ErrInternalError
	}
	artifacts := buildAvailabilityArtifacts(vs, mp, coupon.Discount, isAgent, currenciesMap)
	plansDetailsRowID, err := storePlansDetails(ctx, s.query, artifacts.plansDetails)
	if err != nil {
		rlog.Error("failed to store plans details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &SearchAvailabilityResponse{
		planDetailsID:     plansDetailsRowID,
		AvailableVehicles: artifacts.availableCars,
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
		f := broker.NewFlexWithErpCfg(avCfg.FlexErpDayCharge())
		return &f, nil
	case broker.BrokerHertz:
		return broker.NewHertz(avCfg.HertzErpDayChargeUS(), avCfg.HertzErpDayChargeCA()), nil
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

func getMarkupProviderMap(ctx context.Context, locs availabilityLocations, q db.Querier, rentalDays int, pickupDate string, carGroups []string) (map[broker.Name]MarkupProvider, error) {
	markupProviderMap := make(map[broker.Name]MarkupProvider)
	for brokerName := range locs {
		var provider MarkupProvider
		switch brokerName {
		case broker.BrokerFlex:
			provider = NewFlexMarkupProvider(avCfg.MarkUpGross(), avCfg.MarkUpNet())
		case broker.BrokerHertz:
			hp, err := NewHertzMarkupProvider(ctx, q, locs[brokerName].pickupCountryCode, pickupDate, rentalDays, carGroups)
			if err != nil {
				return nil, fmt.Errorf("initializing hertz markup provider: %w", err)
			}
			provider = hp
		default:
			return nil, fmt.Errorf("unsupported broker: %s", brokerName)
		}
		markupProviderMap[broker.Name(brokerName)] = provider
	}

	return markupProviderMap, nil
}

func extractCarGroups(vs []broker.AvailableVehicle) []string {
	carGroupSet := make(map[string]struct{})
	for _, v := range vs {
		for _, plan := range v.Plans {
			carGroupSet[plan.SupplierCode] = struct{}{}
		}
	}

	carGroups := make([]string, 0, len(carGroupSet))
	for cg := range carGroupSet {
		carGroups = append(carGroups, cg)
	}

	return carGroups
}

type planPriceDetails struct {
	planID                    int
	carModel                  string
	broker                    broker.Name
	rateQualifier             string
	supplierCode              string
	currencyCode              string
	currencyRate              float64
	carPurchasePrice          float64
	carSellPriceWithVat       int //rounded
	carPurchasePriceWithErp   float64
	carSellPriceWithErpAndVat int //rounded
	DiscountPercentage        int
	chargedERPPriceWithVat    int //rounded
}

type availabilityArtifacts struct {
	availableCars []AvailableVehicle
	plansDetails  []planPriceDetails
}

func buildAvailabilityArtifacts(
	availableVehicles []broker.AvailableVehicle,
	markupProviderMap map[broker.Name]MarkupProvider,
	couponDiscount int32,
	isAgent bool,
	currenciesMap map[string]float64) availabilityArtifacts {
	artifacts := availabilityArtifacts{
		availableCars: make([]AvailableVehicle, 0, len(availableVehicles)),
		plansDetails:  make([]planPriceDetails, 0, len(availableVehicles)*2), //most cars have 1-2 plans
	}

	for _, v := range availableVehicles {
		mp, ok := markupProviderMap[v.Broker]
		if !ok {
			rlog.Warn("no markup provider found for broker, skipping applying markup", "broker", v.Broker)
			mp = NewFlexMarkupProvider(avCfg.MarkUpGross(), avCfg.MarkUpNet()) // default to flex markup provider with 0% markup
		}

		av := AvailableVehicle{
			CarDetails:      v.CarDetails,
			Broker:          v.Broker,
			AddOns:          v.AddOns,
			LocationDetails: v.LocationDetails,
			PriceDetails:    v.PriceDetails,
		}
		avPlans := make([]Plan, len(v.Plans))

		for _, p := range v.Plans {
			carPriceWithMarkup := mp.CalculateMarkup(isAgent, p.Price, v.CarDetails.CarGroup, v.CarDetails.SupplierName)
			carPriceWithErpWithMarkup := carPriceWithMarkup
			if p.BrokerErpPrice > 0 {
				carPriceWithErpWithMarkup = mp.CalculateMarkup(isAgent, p.Price+p.BrokerErpPrice, v.CarDetails.CarGroup, v.CarDetails.SupplierName)
			}

			pd := planPriceDetails{
				planID:                    p.PlanID,
				rateQualifier:             p.RateQualifier,
				supplierCode:              p.SupplierCode,
				carModel:                  v.CarDetails.Model,
				broker:                    v.Broker,
				carPurchasePrice:          p.Price,
				carSellPriceWithVat:       calculateDiscountedPrice(carPriceWithMarkup, couponDiscount),
				carPurchasePriceWithErp:   p.Price + p.BrokerErpPrice,
				carSellPriceWithErpAndVat: calculateDiscountedPrice(carPriceWithErpWithMarkup, couponDiscount),
				chargedERPPriceWithVat:    p.ChargedErpPriceWithVat,
				DiscountPercentage:        int(couponDiscount),
				currencyCode:              v.PriceDetails.Currency,
				currencyRate:              currenciesMap[v.PriceDetails.Currency],
			}

			artifacts.plansDetails = append(artifacts.plansDetails, pd)

			avPlan := Plan{
				PlanID:         p.PlanID,
				PlanName:       p.PlanName,
				FullPrice:      roundToInt(carPriceWithMarkup),
				Discount:       int(couponDiscount),
				Price:          calculateDiscountedPrice(carPriceWithMarkup, couponDiscount),
				ErpPrice:       roundToInt(carPriceWithErpWithMarkup-carPriceWithMarkup) + p.ChargedErpPriceWithVat,
				PlanInclusions: p.PlanInclusions,
				Info:           p.Info,
				RateQualifier:  p.RateQualifier,
				SupplierCode:   p.SupplierCode,
			}
			avPlans = append(avPlans, avPlan)
		}

		av.Plans = avPlans
		artifacts.availableCars = append(artifacts.availableCars, av)
	}

	return artifacts
}

// storePlansDetails stores the given plan details in the database and returns the ID of the inserted snapshot.
func storePlansDetails(ctx context.Context, q db.Querier, plans []planPriceDetails) (int64, error) {
	plansJson, err := json.Marshal(plans)
	if err != nil {
		return 0, fmt.Errorf("marshaling plans details: %w", err)
	}

	ID, err := q.InsertAvailablePlansSnapshot(ctx, plansJson)
	if err != nil {
		return 0, fmt.Errorf("inserting available plans snapshot: %w", err)
	}

	return ID, nil
}

// roundToInt rounds a number to the nearest integer.
func roundToInt[T float32 | float64](f T) int {
	return int(math.Round(float64(f)))
}

// calculateDiscountedPrice calculates the price after applying a discount percentage.
func calculateDiscountedPrice(priceBeforeDesc float64, discountPercentage int32) int {
	discountAmount := priceBeforeDesc * float64(discountPercentage) / 100
	return roundToInt(priceBeforeDesc - discountAmount)
}

// buildCurrencyMap query all currencies and returns a map of currency code to currency rates
func buildCurrencyMap(q db.Querier) (map[string]float64, error) {
	currencyMap := make(map[string]float64)
	rows, err := q.ListCurrencies(context.Background())
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		currencyMap[r.CurrencyIsoName] = db.NumericToFloat64(r.Rate)
	}

	return currencyMap, nil
}
