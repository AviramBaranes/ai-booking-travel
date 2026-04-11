package broker

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"

	"encore.app/internal/pricing"
	"encore.dev/rlog"
)

// hertzRequestParams identifies one Hertz plan lookup to perform.
type hertzRequestParams struct {
	brandID  hertzBrand
	planName string
	planCode string
}

// SearchAvailability fetches and merges Hertz vehicle availability results.
func (h Hertz) SearchAvailability(p SearchAvailabilityParams) ([]AvailableVehicle, error) {
	dayCount, err := CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		return nil, fmt.Errorf("calculate days count %w", err)
	}
	requestsParams := make([]hertzRequestParams, 0, 4)
	for brand, plans := range hertzSupplierPlansMap[p.CountryCode] {
		requestsParams = append(requestsParams, hertzRequestParams{
			brandID:  brand,
			planCode: plans.Standard,
			planName: "Standard",
		})
		requestsParams = append(requestsParams, hertzRequestParams{
			brandID:  brand,
			planCode: plans.Gold,
			planName: "Gold",
		})
	}

	av := make([]AvailableVehicle, 0)

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errOnce  sync.Once
		firstErr error
	)
	for _, rp := range requestsParams {
		wg.Add(1)
		go func(rp hertzRequestParams) {
			defer wg.Done()
			xmlReq, err := h.buildSearchAvailabilityRequest(p, rp.brandID, rp.planCode, dayCount)
			if err != nil {
				errOnce.Do(func() {
					firstErr = fmt.Errorf("building xml request %w", err)
				})
				return
			}

			body, err := h.sendXMLRequest(xmlReq)
			if err != nil {
				errOnce.Do(func() {
					firstErr = fmt.Errorf("sending xml request %w", err)
				})
				return
			}

			var raw hertzCarAvailabilityXML
			if err := xml.Unmarshal(body, &raw); err != nil {
				errOnce.Do(func() {
					firstErr = fmt.Errorf("hertz SearchAvailability unmarshal response: %w", err)
				})
				return
			}
			if len(raw.Errors) > 0 {
				errOnce.Do(func() {
					firstErr = fmt.Errorf("hertz SearchAvailability response errors: %+v", raw.Errors[0].ShortText)
				})
				return
			}
			resp := raw.toResponse()

			availableVehicles := h.mapHertzResponseToAvailableVehicles(p, resp, rp.brandID, rp.planName, dayCount)

			mu.Lock()
			av = append(av, availableVehicles...)
			mu.Unlock()
		}(rp)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	avMap := make(map[string]AvailableVehicle)
	for _, v := range av {
		key := fmt.Sprintf("%s-%s-%s", v.CarDetails.Model, v.CarDetails.SupplierName, v.CarDetails.Acriss)
		if existing, ok := avMap[key]; ok {
			existing.Plans = append(existing.Plans, v.Plans...)
			avMap[key] = existing
		} else {
			avMap[key] = v
		}
	}

	out := make([]AvailableVehicle, 0, len(avMap))
	for _, v := range avMap {
		out = append(out, v)
	}

	return out, nil
}

const (
	hertzImageBaseURL = "https://images.hertz.com/vehicles/152x88/"
)

// mapHertzResponseToAvailableVehicles converts a Hertz response into available vehicles.
func (h Hertz) mapHertzResponseToAvailableVehicles(p SearchAvailabilityParams, resp hertzCarAvailabilityResponse, brandID hertzBrand, planName string, dayCount int) []AvailableVehicle {
	availableVehicles := make([]AvailableVehicle, 0, len(resp.Cars))

	inclusions := getPlanInclusions(p.CountryCode, planName)
	locationType := "Shuttle"
	if resp.LocationType == "1" {
		locationType = "Airport"
	}

	for _, car := range resp.Cars {
		if car.Model == "" || car.Model == "DESCRIPTION NOT AVAILABLE" {
			continue
		}

		parts := strings.SplitN(car.Model, " ", 2)
		if len(parts) != 2 {
			rlog.Warn("unexpected car model format in Hertz response, skipping vehicle", "car_model", car.Model)
			continue
		}

		carGroup := parts[0]
		model := parts[1]

		chargeDetails := extractHertzChargeDetails(car.Charges)

		if chargeDetails.price == 0 {
			continue
		}

		info := getInfo(p.DriverAge, dayCount, p.CountryCode, planName)
		if chargeDetails.payAtPickup > 0 {
			info = append([]string{fmt.Sprintf("PayAtPickup:%d:%s", pricing.RoundToInt(chargeDetails.payAtPickup), chargeDetails.payAtPickupCurrency)}, info...)
		}

		availableVehicles = append(availableVehicles, AvailableVehicle{
			Broker: h.Name(),
			CarDetails: CarDetails{
				Model:        normalizeModelName(model),
				CarGroup:     carGroup,
				ImageURL:     hertzImageBaseURL + strings.TrimSpace(car.ImageURL),
				SupplierName: mapBrandIDToSupplierName(brandID),
				CarType:      mapCarTypeCodeToCarType(car.CarType),
				Acriss:       car.Acriss,
				HasAC:        car.HasAC,
				IsAutoGear:   car.TransmissionType == "Automatic",
				IsElectric:   isElectric(car.Acriss),
				Seats:        car.Seats,
				Bags:         car.Bags,
				Doors:        car.Doors,
			},
			Plans: []Plan{
				{
					PlanName:               planName,
					Price:                  chargeDetails.price,
					PlanInclusions:         inclusions,
					Info:                   info,
					BrokerErpPrice:         0,
					ChargedErpPriceWithVat: h.getERPPrice(dayCount, p.CountryCode),
					RateQualifier:          car.ID,
					SupplierCode:           string(brandID),
				},
			},
			LocationDetails: LocationDetails{
				LocationType: locationType,
			},
			PriceDetails: PriceDetails{
				Currency:           chargeDetails.currency,
				DropCharge:         pricing.RoundToInt(chargeDetails.dropCharge),
				DropChargeCurrency: chargeDetails.dropChargeCurrency,
			},
		})
	}

	return availableVehicles
}

// getERPPrice returns the ERP surcharge for the market and rental length.
func (h Hertz) getERPPrice(dayCount int, countryCode string) int {
	if countryCode == "US" {
		return h.usErpDayCharge * dayCount
	}
	if countryCode == "CA" {
		return h.caErpDayCharge * dayCount
	}
	return 0
}

// getPlanInclusions returns a copy of the inclusions for a Hertz plan.
func getPlanInclusions(countryCode, planName string) []string {
	var incs []string
	if countryCode == "US" {
		incs = hertzUSAInclusionsMap[planName]
	}
	if countryCode == "CA" {
		incs = hertzCanadaInclusionsMap[planName]
	}

	out := make([]string, len(incs))
	copy(out, incs)
	return out
}

// getInfo builds plan notes for the driver, market, and rental length.
func getInfo(driverAge, dayCount int, countryCode, planName string) []string {
	info, _ := hertzTermsMap[planName]
	infoCopy := make([]string, len(info))
	copy(infoCopy, info)

	if driverAge < 25 {
		if countryCode == "US" {
			extraCharge := 29 * dayCount
			infoCopy = append(infoCopy, fmt.Sprintf("YoungDriverFee:%d:$", extraCharge))
		}
		if countryCode == "CA" {
			extraCharge := 15 * dayCount
			infoCopy = append(infoCopy, fmt.Sprintf("YoungDriverFee:%d:CAD", extraCharge))
		}
	}

	return infoCopy
}

// hertzChargeDetails holds the price-related values parsed from Hertz charges.
type hertzChargeDetails struct {
	price               float64
	currency            string
	dropCharge          float64
	dropChargeCurrency  string
	payAtPickup         float64
	payAtPickupCurrency string
}

// extractHertzChargeDetails pulls pricing fields from Hertz charge entries.
func extractHertzChargeDetails(charges []hertzCharge) hertzChargeDetails {
	var details hertzChargeDetails

	for _, charge := range charges {
		switch charge.Purpose {
		case "1":
			if charge.GuaranteedInd {
				details.price = charge.Amount
				details.currency = charge.CurrencyCode
			}
		case "2":
			details.dropCharge = charge.Amount
			details.dropChargeCurrency = charge.CurrencyCode
		case "23":
			details.payAtPickup = charge.Amount
			details.payAtPickupCurrency = charge.CurrencyCode
		}
	}

	return details
}
