package availability

import (
	"context"
	"math/rand/v2"
	"sort"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/lang"
	"encore.app/internal/pricing"
	"encore.app/services/accounts"
	auth "encore.app/services/accounts"
	"encore.dev/rlog"
)

// availabilityArtifacts groups the processed availability results: the API-facing vehicles and the stored plan details.
type availabilityArtifacts struct {
	availableCars []AvailableVehicle
	plansDetails  []PlanPriceDetails
}

const (
	minLiveViewers       = 25
	maxLiveViewers       = 112
	minRemainingCount    = 1
	maxRemainingCount    = 4
	maxVehiclesWithFlags = 20
	topChoiceTag         = "top choice"
)

// buildAvailabilityArtifacts applies markup, coupon discounts, and currency data to produce the final response vehicles and plan snapshots.
func (s *AvailabilityService) buildAvailabilityArtifacts(ctx context.Context, p SearchAvailabilityParams, locs availabilityLocations, avResp *broker.AvailabilityResponse, couponDiscount float64) (availabilityArtifacts, error) {
	artifacts := availabilityArtifacts{
		availableCars: make([]AvailableVehicle, 0, len(avResp.AvailableVehicles)),
		plansDetails:  make([]PlanPriceDetails, 0, len(avResp.AvailableVehicles)*2), //most cars have 1-2 plans
	}

	daysCount, err := broker.CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		rlog.Error("failed to calculate rental days count", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	markupProviders := s.getMarkupProviderMap(ctx, locs, daysCount, p.PickupDate, extractCarGroups(avResp.AvailableVehicles))

	grossMarkupResp, err := accounts.GetUserMarkupGross(ctx)
	if err != nil {
		rlog.Error("failed to get user gross markup", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	authData := auth.GetAuthData()
	isAgent := authData != nil && (authData.Role == auth.UserRoleAgent)
	currenciesMap, err := s.buildCurrencyMap(ctx)
	if err != nil {
		rlog.Error("failed to build currency map", "error", err)
		return artifacts, api_errors.ErrInternalError
	}

	spMap := make(map[string]*broker.SupplierInfo)
	translatedSuppliers := make(map[string]bool)
	for i, sp := range avResp.SuppliersInfo {
		spMap[sp.Name] = &avResp.SuppliersInfo[i]
		translatedSuppliers[sp.Name] = false
	}

	for i, v := range avResp.AvailableVehicles {
		mp, ok := markupProviders[v.Broker]
		if !ok {
			rlog.Warn("no markup provider found for broker, skipping applying markup", "broker", v.Broker)
			mp = NewMarkupProviderFromCfg(s.cfg.MarkUpGross(), s.cfg.MarkUpNet())
		}

		av := AvailableVehicle{
			ID:              i,
			CarDetails:      v.CarDetails,
			Broker:          v.Broker,
			LocationDetails: v.LocationDetails,
			PriceDetails:    v.PriceDetails,
		}
		avPlans := make([]Plan, 0, len(v.Plans))

		for _, p := range v.Plans {
			sp, ok := spMap[p.SupplierName]
			if !ok {
				rlog.Warn("no supplier info found for supplier name, skipping plan details enrichment", "supplierName", p.SupplierName)
				continue
			}

			markupPercentage := mp.GetMarkup(isAgent)
			if markupPercentage <= 0 {
				rlog.Warn("calculated car price with markup is less than or equal to 0, skipping plan", "carGroup", v.CarDetails.CarGroup, "brand", p.SupplierCode)
				continue
			}
			cr, ok := currenciesMap[v.PriceDetails.Currency]
			if !ok {
				rlog.Warn("no currency rate found for currency, skipping plan", "currency", v.PriceDetails.Currency)
				continue
			}

			brokerLoc, ok := locs[v.Broker]
			if !ok {
				rlog.Warn("no location data found for broker, skipping plan", "broker", v.Broker)
				continue
			}

			if lang.FromContext(ctx, "en") == "he" {
				if !translatedSuppliers[sp.Name] {
					for i, inc := range sp.Inclusions {
						sp.Inclusions[i].ProductInclusions = s.translatePlanDetails(ctx, inc.ProductInclusions)
					}
					translatedSuppliers[sp.Name] = true
				}
				p.Info = s.translatePlanDetails(ctx, p.Info)
			}

			var incs []string
			for _, prd := range sp.Inclusions {
				if prd.ProductName == p.PlanName {
					incs = prd.ProductInclusions
					break
				}
			}
			pd := PlanPriceDetails{
				PlanID:                 p.PlanID,
				RateQualifier:          p.RateQualifier,
				SupplierCode:           p.SupplierCode,
				Broker:                 v.Broker,
				PickupLocationCode:     brokerLoc.pickupBrokerLocationID,
				DropoffLocationCode:    brokerLoc.dropoffBrokerLocationID,
				CurrencyCode:           v.PriceDetails.Currency,
				CurrencyRate:           cr,
				DiscountPercentage:     couponDiscount,
				CarPurchasePrice:       p.Price,
				MarkupPercentage:       markupPercentage,
				SupplierErpPrice:       p.BrokerErpPrice,
				ChargedERPPriceWithVat: p.ChargedErpPriceWithVat,
				CarDetails:             v.CarDetails,
				Inclusions:             incs,
				AvailableAddOns:        sp.AddOns,
				Fees:                   v.PriceDetails.Fees,
			}

			artifacts.plansDetails = append(artifacts.plansDetails, pd)

			carPriceWithMarkup := pricing.ApplyMarkup(p.Price, markupPercentage)
			erpWithMarkup := pricing.ApplyMarkup(p.BrokerErpPrice, markupPercentage)
			discountedErp := pricing.CalculateDiscountedPrice(erpWithMarkup, couponDiscount)
			discountedCarPrice := pricing.CalculateDiscountedPrice(carPriceWithMarkup, couponDiscount) // no discount on charged erp

			// with gross:
			carPriceWithMarkup = pricing.ApplyMarkup(carPriceWithMarkup, grossMarkupResp.GrossMarkup)
			discountedCarPrice = pricing.ApplyMarkup(discountedCarPrice, grossMarkupResp.GrossMarkup)
			discountedErp = pricing.ApplyMarkup(discountedErp, grossMarkupResp.GrossMarkup)
			chargedErpPriceWithVat := pricing.ApplyMarkup(p.ChargedErpPriceWithVat, grossMarkupResp.GrossMarkup)

			avPlan := Plan{
				PlanID:        p.PlanID,
				PlanName:      p.PlanName,
				FullPrice:     pricing.RoundToInt(carPriceWithMarkup),
				Discount:      pricing.RoundToInt(couponDiscount),
				Price:         pricing.RoundToInt(discountedCarPrice),
				ErpPrice:      pricing.RoundToInt(discountedErp + chargedErpPriceWithVat), // no discount on charged erp
				Info:          p.Info,
				RateQualifier: p.RateQualifier,
				SupplierName:  sp.Name,
				SupplierCode:  p.SupplierCode,
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

	assignRandomBookingSignals(artifacts.availableCars)

	return artifacts, nil
}

func assignRandomBookingSignals(vs []AvailableVehicle) {
	if len(vs) == 0 {
		return
	}

	selectedCount := randomSignalVehicleCount(len(vs))
	selectedIndices := randomVehicleIndices(len(vs), selectedCount)

	for _, idx := range selectedIndices {
		vs[idx].Signals = newRandomBookingSignals()
	}
}

func randomSignalVehicleCount(totalVehicles int) int {
	maxSelected := min(totalVehicles, maxVehiclesWithFlags)
	return randomIntInRange(1, maxSelected)
}

func randomVehicleIndices(totalVehicles, selectedCount int) []int {
	indices := make([]int, totalVehicles)
	for i := 0; i < totalVehicles; i++ {
		indices[i] = i
	}

	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	return indices[:selectedCount]
}

func newRandomBookingSignals() *BookingSignals {
	return &BookingSignals{
		LiveViewers:    randomIntInRange(minLiveViewers, maxLiveViewers),
		RemainingCount: randomIntInRange(minRemainingCount, maxRemainingCount),
		Tags:           []string{topChoiceTag},
	}
}

func randomIntInRange(minValue, maxValue int) int {
	return rand.IntN(maxValue-minValue+1) + minValue
}

// sortPlansByPrice sorts the plans in-place by their price in ascending order.
func sortPlansByPrice(plans []Plan) {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Price < plans[j].Price
	})
}

// getMarkupProviderMap returns a MarkupProvider for each broker present in the availability locations.
func (s *AvailabilityService) getMarkupProviderMap(ctx context.Context, locs availabilityLocations, rentalDays int, pickupDate string, carGroups []string) map[broker.Name]MarkupProvider {
	markupProviderMap := make(map[broker.Name]MarkupProvider)
	for brokerName := range locs {
		provider, err := NewMarkupProviderFromDB(ctx, s.query, locs[brokerName].pickupCountryCode, string(brokerName))
		if err != nil {
			rlog.Error("initializing flex markup provider: %w", err)
			provider = NewMarkupProviderFromCfg(s.cfg.MarkUpGross(), s.cfg.MarkUpNet())
		}
		markupProviderMap[brokerName] = provider
	}

	return markupProviderMap
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

// sortAvailableVehiclesByCheapestPlan sorts the available vehicles in-place by the price of their cheapest plan.
func sortAvailableVehiclesByCheapestPlan(vs []AvailableVehicle) {
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Plans[0].Price < vs[j].Plans[0].Price
	})
}

// buildCurrencyMap query all currencies and returns a map of currency code to currency rates
func (s *AvailabilityService) buildCurrencyMap(ctx context.Context) (map[string]float64, error) {
	currencyMap := make(map[string]float64)
	rows, err := s.query.ListCurrencies(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		currencyMap[r.CurrencyIsoName] = dbadapters.NumericToFloat64(r.Rate)
	}

	return currencyMap, nil
}

// translatePlanDetails translates the plan details using the service's translator, returning the original detail if no translation is found after inserting it to db.
func (s *AvailabilityService) translatePlanDetails(ctx context.Context, details []string) []string {
	translatedDetails := make([]string, len(details))
	for i, detail := range details {
		template, values := s.tc.NormalizeSentence(detail)
		if translated, exists := s.tc.GetVerified(template); exists {
			translatedDetails[i] = s.tc.InsertValuesToSentence(translated, values)
		} else {
			translatedDetails[i] = detail
			if !s.tc.Exists(template) {
				_, err := s.query.InsertBrokerTranslation(ctx, template)
				if err != nil {
					rlog.Error("failed to insert missing translation to db", "template", template, "error", err)
				}
				s.tc.AddKnown(template)
			}
		}
	}
	return translatedDetails
}
