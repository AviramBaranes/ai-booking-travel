package broker

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"sync"

	"encore.dev/rlog"
)

type hertzRequestParams struct {
	brandID  hertzBrand
	planName string
	planCode string
}

// SearchAvailability searches for available vehicles based on the provided search parameters. It returns a slice of AvailableVehicle structs containing details about the available vehicles, or an error if the search fails.
func (h Hertz) SearchAvailability(p SearchAvailabilityParams) ([]AvailableVehicle, error) {
	dayCount, err := calculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
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

			var resp hertzCarAvailabilityResponse
			if err := xml.Unmarshal(body, &resp); err != nil {
				errOnce.Do(func() {
					firstErr = fmt.Errorf("hertz SearchAvailability unmarshal response: %w", err)
				})
				return
			}

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
			sort.Slice(existing.Plans, func(i, j int) bool {
				return existing.Plans[i].Price < existing.Plans[j].Price
			})
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

// mapHertzResponseToAvailableVehicles maps the Hertz API response to a slice of AvailableVehicle structs, extracting relevant details about each available vehicle and its rental plans.
func (h Hertz) mapHertzResponseToAvailableVehicles(p SearchAvailabilityParams, resp hertzCarAvailabilityResponse, brandID hertzBrand, planName string, dayCount int) []AvailableVehicle {
	availableVehicles := make([]AvailableVehicle, 0, len(resp.Cars))

	inclusions := getPlanInclusions(p.CountryCode, planName)
	locationType := "Shuttle"
	if resp.LocationInfo.LocationDetails.AdditionalInfo.ParkLocation.LocationType == "1" {
		locationType = "Airport"
	}

	for _, car := range resp.Cars {
		if car.VehAvailCore.Vehicle.VehMakeModel.Model == "" || car.VehAvailCore.Vehicle.VehMakeModel.Model == "DESCRIPTION NOT AVAILABLE" {
			continue
		}

		parts := strings.SplitN(car.VehAvailCore.Vehicle.VehMakeModel.Model, " ", 2)
		if len(parts) != 2 {
			rlog.Warn("unexpected car model format in Hertz response, skipping vehicle", "car_model", car.VehAvailCore.Vehicle.VehMakeModel.Model)
			continue
		}

		carGroup := parts[0]
		model := parts[1]

		var (
			price               float64
			currency            string
			dropCharge          float64
			dropChargeCurrency  string
			payAtPickup         float64
			payAtPickupCurrency string
		)

		for _, charge := range car.VehAvailCore.RentalRate.Charges {

			switch charge.Purpose {
			case "1":
				if charge.GuaranteedInd {
					price = charge.Amount
					currency = charge.CurrencyCode
				}
			case "2":
				dropCharge = charge.Amount
				dropChargeCurrency = charge.CurrencyCode
			case "23":
				payAtPickup = charge.Amount
				payAtPickupCurrency = charge.CurrencyCode
			}
		}

		if price == 0 {
			continue
		}

		info := getInfo(p.DriverAge, dayCount, p.CountryCode, planName)
		if payAtPickup > 0 {
			info = append([]string{fmt.Sprintf("PayAtPickup:%d:%s", roundToInt(payAtPickup), payAtPickupCurrency)}, info...)
		}

		availableVehicles = append(availableVehicles, AvailableVehicle{
			Broker: h.Name(),
			CarDetails: CarDetails{
				Model:        model,
				CarGroup:     carGroup,
				ImageURL:     hertzImageBaseURL + strings.TrimSpace(car.VehAvailCore.Vehicle.ImageURL),
				SupplierName: mapBrandIDToSupplierName(brandID),
				CarType:      mapCarTypeCodeToCarType(car.VehAvailCore.Vehicle.VehType.CarType),
				Acriss:       car.VehAvailCore.Vehicle.Acriss,
				HasAC:        car.VehAvailCore.Vehicle.HasAC,
				IsAutoGear:   strings.EqualFold(car.VehAvailCore.Vehicle.TransmissionType, "Automatic"),
				Seats:        car.VehAvailCore.Vehicle.Seats,
				Bags:         car.VehAvailCore.Vehicle.Bags,
				Doors:        car.VehAvailCore.Vehicle.Doors,
			},
			Plans: []Plan{
				{
					PlanName:       planName,
					FullPrice:      roundToInt(price),
					Discount:       p.DiscountPercentage,
					Price:          roundToInt(price),
					PlanInclusions: inclusions,
					Info:           info,
					ErpPrice:       roundToInt(getERPPrice(dayCount, p.CountryCode)),
					RateQualifier:  car.VehAvailCore.Reference.ID,
					SupplierCode:   string(brandID),
				},
			},
			LocationDetails: LocationDetails{
				LocationType: locationType,
			},
			PriceDetails: PriceDetails{
				Currency:           currency,
				DropCharge:         roundToInt(dropCharge),
				DropChargeCurrency: dropChargeCurrency,
			},
		})
	}

	return availableVehicles
}

func getERPPrice(dayCount int, countryCode string) float32 {
	const US_ERP_DAY_CHARGE float32 = 2.56
	const CA_ERP_DAY_CHARGE float32 = 5.98

	if countryCode == "US" {
		return US_ERP_DAY_CHARGE * float32(dayCount)
	}
	if countryCode == "CA" {
		return CA_ERP_DAY_CHARGE * float32(dayCount)
	}
	return 0
}

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
