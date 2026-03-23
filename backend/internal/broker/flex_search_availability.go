package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"encore.dev/rlog"
)

// SearchAvailability searches for available vehicles based on the provided search parameters. It returns a slice of AvailableVehicle structs containing details about the available vehicles, or an error if the search fails.
func (f Flex) SearchAvailability(p SearchAvailabilityParams) ([]AvailableVehicle, error) {
	form := url.Values{}
	form.Set("SIPP", "")
	form.Set("SupplierCode", "")
	form.Set("ProductID", "1")
	form.Set("Language", "UK")
	form.Set("AdditionalParameters", "Timeout=15000")
	form.Set("PickupLocationID", p.PickupLocation)
	form.Set("DropoffLocationID", p.DropoffLocation)
	form.Set("PickupDate", formatDate(p.PickupDate))
	form.Set("DropoffDate", formatDate(p.DropoffDate))
	form.Set("PickUpTime", p.PickupTime)
	form.Set("DropoffTime", p.DropoffTime)
	form.Set("DriversAge", strconv.Itoa(p.DriverAge))

	dayCount, err := CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rental days count: %w", err)
	}

	body, err := f.postForm("CarAvailability", form)
	if err != nil {
		return nil, err
	}

	var resp flexCarAvailabilityResponse

	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("flex CarAvailability unmarshal response: %w", err)
	}

	if resp.ReturnCode != 0 {
		return nil, fmt.Errorf("CarAvailability API returned error code %d with message: %s", resp.ReturnCode, resp.ErrorMessage)
	}

	if len(resp.Cars) == 0 {
		rlog.Info("no cars found in CarAvailability response", "pickup_location", p.PickupLocation, "dropoff_location", p.DropoffLocation, "pickup_date", p.PickupDate, "dropoff_date", p.DropoffDate)
		return []AvailableVehicle{}, nil
	}

	if len(resp.SupplierDetails) == 0 {
		return []AvailableVehicle{}, fmt.Errorf("no supplier details found in CarAvailability response for pickup_location=%s dropoff_location=%s pickup_date=%s dropoff_date=%s", p.PickupLocation, p.DropoffLocation, p.PickupDate, p.DropoffDate)
	}

	addOnsMap := createAddOnMap(resp.SupplierDetails)
	supplierDetailsMap := createSupplierMap(resp.SupplierDetails)

	carsMap := make(map[string]AvailableVehicle)
	for _, c := range resp.Cars {
		s, ok := flexSupplierMap[c.SupplierCode]
		if !ok {
			rlog.Warn("unknown supplier code in CarAvailability response, skipping vehicle", "supplier_code", c.SupplierCode)
			continue
		}

		supplierDetails, ok := supplierDetailsMap[s.name]
		if !ok {
			rlog.Warn("no supplier details found for supplier in CarAvailability response, using empty details", "supplier_name", s.name)
			continue
		}

		addOns, ok := addOnsMap[s.name]
		if !ok {
			rlog.Warn("no add-ons found for supplier in CarAvailability response, using empty add-on list", "supplier_name", s.name)
			addOns = []AddOn{}
		}

		plans := f.getPlans(c, dayCount, supplierDetails, p.DiscountPercentage)
		if len(plans) == 0 {
			rlog.Warn("no valid plans found for car in CarAvailability response, skipping vehicle", "car_name", c.Name)
			continue
		}

		carMapID := fmt.Sprintf("%s-%s-%s", c.Name, s.code, c.Code)
		if car, ok := carsMap[carMapID]; ok {
			car.Plans = append(car.Plans, plans...)
			carsMap[carMapID] = car
			continue
		}

		carDetails, err := flexCarToBrokerCar(c, s.name)
		if err != nil {
			rlog.Warn("failed to convert flex car to broker car details, skipping vehicle", "error", err)
			continue
		}

		car := AvailableVehicle{
			Broker:          BrokerFlex,
			CarDetails:      carDetails,
			Plans:           plans,
			AddOns:          addOns,
			LocationDetails: getLocationDetails(supplierDetails),
			PriceDetails: PriceDetails{
				Currency:           c.Currency,
				DropCharge:         roundToInt(c.DropCharge),
				DropChargeCurrency: c.DropChargeCurrency,
			},
		}

		carsMap[carMapID] = car
	}

	out := make([]AvailableVehicle, 0, len(resp.Cars))
	for _, car := range carsMap {
		out = append(out, car)
	}

	return out, nil
}

// createAddOnMap returns a map of supplier name to addOn slice
func createAddOnMap(suppliers []flexSupplierDetails) map[string][]AddOn {
	addOnMap := make(map[string][]AddOn)
	for _, s := range suppliers {
		addOns := make([]AddOn, 0, len(s.AvailableExtras))
		for _, e := range s.AvailableExtras {
			addOns = append(addOns, AddOn{
				ID:              e.ExtraID,
				Name:            e.Name,
				Price:           int(e.Price),
				AllowedQuantity: e.MaxAmount,
				Period:          e.Period,
				Currency:        e.Currency,
			})
		}
		addOnMap[s.Supplier] = addOns
	}

	return addOnMap
}

// createSupplierMap maps Flex supplierInfo by name
func createSupplierMap(suppliers []flexSupplierDetails) map[string]flexSupplierDetails {
	supplierMap := make(map[string]flexSupplierDetails)
	for _, s := range suppliers {
		supplierMap[s.Supplier] = s
	}

	return supplierMap
}

// flexProductMap maps flex product name to its ids
var flexProductMap = map[string]int{
	"Inclusive":            1,
	"Inclusive GPS":        2,
	"Gold":                 3,
	"Gold GPS":             4,
	"Young Driver Package": 10,
}

// getInsuranceExtraCost calculates the extra insurance cost based on the number of rental days, using a fixed daily rate.
func (f Flex) getInsuranceExtraCost(days int) int {
	return days * f.erpDayCharge
}

// getPlans returns the list of plans for a given car
func (f Flex) getPlans(c flexCar, dayCount int, supplierDetails flexSupplierDetails, discount int) []Plan {
	plans := make([]Plan, 0, len(c.Costs))
	for _, p := range c.Costs {
		planID, ok := flexProductMap[p.Product]
		if !ok {
			planID = 1
		}

		var planInclusions []string
		for _, inc := range supplierDetails.Inclusions {
			if inc.Product == p.Product {
				planInclusions = strings.Split(inc.Inclusion, ";")
				break
			}
		}

		plans = append(plans, Plan{
			PlanID:                 planID,
			PlanName:               p.Product,
			PlanInclusions:         planInclusions,
			Price:                  p.Price,
			BrokerErpPrice:         c.ERP,
			ChargedErpPriceWithVat: f.getInsuranceExtraCost(dayCount),
			Info:                   c.Information,
			RateQualifier:          c.RateQualifier,
			SupplierCode:           c.SupplierCode,
		})
	}

	return plans
}

// flexCarToBrokerCar converts a flexCar to a CarDetails struct for the broker response.
func flexCarToBrokerCar(c flexCar, supplierName string) (CarDetails, error) {
	seats, err := strconv.Atoi(c.Passenger)
	if err != nil {
		return CarDetails{}, fmt.Errorf("flexCarToBrokerCar parse seats got %s: %w", c.Passenger, err)
	}

	doors, err := strconv.Atoi(c.Doors)
	if err != nil {
		return CarDetails{}, fmt.Errorf("flexCarToBrokerCar parse doors got %s: %w", c.Doors, err)
	}

	bags, err := strconv.Atoi(c.Luggage)
	if err != nil {
		return CarDetails{}, fmt.Errorf("flexCarToBrokerCar parse bags got %s: %w", c.Luggage, err)
	}

	return CarDetails{
		Model:        c.Name,
		ImageURL:     c.URL,
		SupplierName: supplierName,
		CarType:      c.CarType,
		Acriss:       c.Code,
		HasAC:        strings.HasPrefix(c.IsAirCon, "Y"),
		IsAutoGear:   strings.HasPrefix(c.IsAutomatic, "Y"),
		Seats:        seats,
		Doors:        doors,
		Bags:         bags,
	}, nil
}

// getLocationDetails extracts location details from the supplier details
func getLocationDetails(s flexSupplierDetails) LocationDetails {
	var dc string
	const DELIVERY_COLLECTION_HEADER = "DELIVERY AND COLLECTION"
	for _, term := range s.Terms {
		if strings.HasPrefix(term.Header, DELIVERY_COLLECTION_HEADER) {
			dc = term.Paragraph
			break
		}
	}

	return LocationDetails{
		DeliveryCollection:  dc,
		PickupBranchAddress: s.PickUpDetails.Address1,
		ReturnBranchAddress: s.DropOffDetails.Address1,
		PickupBranchPhone:   s.PickUpDetails.Phone,
		ReturnBranchPhone:   s.DropOffDetails.Phone,
		LocationType:        s.PickUpDetails.LocationType,
		PickupNotes:         s.PickUpDetails.LocationInformation,
	}
}

// formatDate formats a date string from "2006-01-02" to "02/01/2006"
func formatDate(dateStr string) string {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return dateStr
	}
	return parts[2] + "/" + parts[1] + "/" + parts[0]
}
