package booking

import (
	"context"
	"fmt"
	"math"
	"sort"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/auth"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// availabilityArtifacts groups the processed availability results: the API-facing vehicles and the stored plan details.
type availabilityArtifacts struct {
	availableCars []AvailableVehicle
	plansDetails  []planPriceDetails
}

// buildAvailabilityArtifacts applies markup, coupon discounts, and currency data to produce the final response vehicles and plan snapshots.
func (s *Service) buildAvailabilityArtifacts(ctx context.Context, p SearchAvailabilityRequest, locs availabilityLocations, rawVehicles []broker.AvailableVehicle, couponDiscount int) (availabilityArtifacts, error) {
	artifacts := availabilityArtifacts{
		availableCars: make([]AvailableVehicle, 0, len(rawVehicles)),
		plansDetails:  make([]planPriceDetails, 0, len(rawVehicles)*2), //most cars have 1-2 plans
	}

	daysCount, err := broker.CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		rlog.Error("failed to calculate rental days count", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	markupProviders, err := getMarkupProviderMap(ctx, locs, s.query, daysCount, p.PickupDate, extractCarGroups(rawVehicles))
	if err != nil {
		rlog.Error("failed to get markup provider map", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	authData := auth.GetAuthData()
	isAgent := authData != nil && (authData.Role == auth.UserRoleAgent)
	currenciesMap, err := buildCurrencyMap(ctx, s.query)
	if err != nil {
		rlog.Error("failed to build currency map", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	for _, v := range rawVehicles {
		mp, ok := markupProviders[v.Broker]
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
		avPlans := make([]Plan, 0, len(v.Plans))

		for _, p := range v.Plans {
			carPriceWithMarkup := mp.CalculateMarkup(isAgent, p.Price, v.CarDetails.CarGroup, p.SupplierCode)
			if carPriceWithMarkup <= 0 {
				rlog.Warn("calculated car price with markup is less than or equal to 0, skipping plan", "carGroup", v.CarDetails.CarGroup, "brand", p.SupplierCode)
				continue
			}
			cr, ok := currenciesMap[v.PriceDetails.Currency]
			if !ok {
				rlog.Warn("no currency rate found for currency, skipping plan", "currency", v.PriceDetails.Currency)
				continue
			}

			carPriceWithErpWithMarkup := carPriceWithMarkup
			if p.BrokerErpPrice > 0 {
				carPriceWithErpWithMarkup = mp.CalculateMarkup(isAgent, p.Price+p.BrokerErpPrice, v.CarDetails.CarGroup, p.SupplierCode)
			}

			pd := planPriceDetails{
				PlanID:                    p.PlanID,
				RateQualifier:             p.RateQualifier,
				SupplierCode:              p.SupplierCode,
				CarModel:                  v.CarDetails.Model,
				Broker:                    v.Broker,
				CarPurchasePrice:          p.Price,
				CarSellPriceWithVat:       calculateDiscountedPrice(carPriceWithMarkup, couponDiscount),
				CarPurchasePriceWithErp:   p.Price + p.BrokerErpPrice,
				CarSellPriceWithErpAndVat: calculateDiscountedPrice(carPriceWithErpWithMarkup, couponDiscount),
				ChargedERPPriceWithVat:    p.ChargedErpPriceWithVat,
				DiscountPercentage:        int(couponDiscount),
				CurrencyCode:              v.PriceDetails.Currency,
				CurrencyRate:              cr,
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

		if len(avPlans) == 0 {
			continue
		}

		sortPlansByPrice(avPlans)
		av.Plans = avPlans
		artifacts.availableCars = append(artifacts.availableCars, av)
	}

	return artifacts, nil
}

// sortPlansByPrice sorts the plans in-place by their price in ascending order.
func sortPlansByPrice(plans []Plan) {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Price < plans[j].Price
	})
}

// getMarkupProviderMap returns a MarkupProvider for each broker present in the availability locations.
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

// extractCarGroups returns a deduplicated slice of car group codes from the given vehicles.
func extractCarGroups(vs []broker.AvailableVehicle) []string {
	carGroupSet := make(map[string]struct{})
	for _, v := range vs {
		carGroupSet[v.CarDetails.CarGroup] = struct{}{}
	}

	carGroups := make([]string, 0, len(carGroupSet))
	for cg := range carGroupSet {
		carGroups = append(carGroups, cg)
	}

	return carGroups
}

// calculateDiscountedPrice calculates the price after applying a discount percentage.
func calculateDiscountedPrice(priceBeforeDesc float64, discountPercentage int) int {
	discountAmount := priceBeforeDesc * float64(discountPercentage) / 100
	return roundToInt(priceBeforeDesc - discountAmount)
}

// sortAvailableVehiclesByCheapestPlan sorts the available vehicles in-place by the price of their cheapest plan.
func sortAvailableVehiclesByCheapestPlan(vs []AvailableVehicle) {
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Plans[0].Price < vs[j].Plans[0].Price
	})
}

// buildCurrencyMap query all currencies and returns a map of currency code to currency rates
func buildCurrencyMap(ctx context.Context, q db.Querier) (map[string]float64, error) {
	currencyMap := make(map[string]float64)
	rows, err := q.ListCurrencies(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		currencyMap[r.CurrencyIsoName] = db.NumericToFloat64(r.Rate)
	}

	return currencyMap, nil
}

// roundToInt rounds a number to the nearest integer.
func roundToInt[T float32 | float64](f T) int {
	return int(math.Round(float64(f)))
}
