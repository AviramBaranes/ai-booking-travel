package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"encore.app/internal/pricing"
	"encore.dev/rlog"
)

// SearchAvailability searches for available vehicles based on the provided search parameters. It returns a slice of AvailableVehicle structs containing details about the available vehicles, or an error if the search fails.
func (f *Flex) SearchAvailability(p SearchAvailabilityParams) (*AvailabilityResponse, error) {
	form := url.Values{}
	form.Set("SIPP", "")
	form.Set("SupplierCode", "")

	productID := "1"
	if p.CountryCode == "US" || p.CountryCode == "CA" {
		productID = "1,3" // Include both "Inclusive" and "Gold" products for US and CA markets
	}
	form.Set("ProductID", productID)
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

	if resp.Count == 0 {
		rlog.Info("no cars found in CarAvailability response", "pickup_location", p.PickupLocation, "dropoff_location", p.DropoffLocation, "pickup_date", p.PickupDate, "dropoff_date", p.DropoffDate)
		return &AvailabilityResponse{}, nil
	}

	rlog.Info("CarAvailability response", "return_code", resp.ReturnCode, "error_message", resp.ErrorMessage, "car_count", len(resp.Cars), "supplier_details_count", len(resp.SupplierDetails))

	if len(resp.SupplierDetails) == 0 {
		return &AvailabilityResponse{}, fmt.Errorf("no supplier details found in CarAvailability response for pickup_location=%s dropoff_location=%s pickup_date=%s dropoff_date=%s", p.PickupLocation, p.DropoffLocation, p.PickupDate, p.DropoffDate)
	}

	supplierDetailsMap := createSupplierMap(resp.SupplierDetails)

	carsMap := make(map[string]AvailableVehicle)
	for _, c := range resp.Cars {
		s, ok := flexSupplierMap[c.SupplierCode]
		if !ok {
			rlog.Warn("unknown supplier code in CarAvailability response, skipping vehicle", "supplier_code", c.SupplierCode)
			continue
		}

		supplierDetails, ok := supplierDetailsMap[c.Supplier]
		if !ok {
			rlog.Warn("no supplier details found for supplier in CarAvailability response, using empty details", "supplier_name", s.name)
			continue
		}

		plans := f.getPlans(c, dayCount, supplierDetails)
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

		carDetails := flexCarToBrokerCar(c, s.name)

		ydFee, ydFeeCurrency := f.getYoungDriverFee(c.Information)

		car := AvailableVehicle{
			Broker:     BrokerFlex,
			CarDetails: carDetails,
			Plans:      plans,
			LocationDetails: LocationDetails{
				LocationType: supplierDetails.PickUpDetails.LocationType,
			},
			PriceDetails: PriceDetails{
				Currency: c.Currency,
				Fees: Fees{
					DropCharge:             pricing.RoundToInt(parseFloat(c.DropCharge)),
					DropChargeCurrency:     c.DropChargeCurrency,
					YoungDriverFee:         ydFee,
					YoungDriverFeeCurrency: ydFeeCurrency,
				},
			},
		}

		if car.LocationDetails.LocationType != "Airport" && car.LocationDetails.LocationType != "Shuttle" && car.LocationDetails.LocationType != "City" {
			rlog.Warn("unexpected location type in CarAvailability response, expected 'Airport', 'Shuttle', or 'City'", "location_type", car.LocationDetails.LocationType)
			continue
		}

		carsMap[carMapID] = car
	}

	out := make([]AvailableVehicle, 0, len(resp.Cars))
	for _, car := range carsMap {
		out = append(out, car)
	}

	return &AvailabilityResponse{
		AvailableVehicles: out,
		SuppliersInfo:     f.parseSupplierDetails(resp.SupplierDetails),
	}, nil
}

// getYoungDriverFee returns the young driver fee and its currency for the given driver age, rental length, and market.
// Each info item may be a comma-separated list of key:value:unit segments, e.g.:
//
//	"YoungDriverFee:29:$"
//	"MANDATORY CHARGES - OneWay:149.99:EUR,YoungDriverFee:150.00:EUR,moreinfo:10:$"
func (f *Flex) getYoungDriverFee(info []string) (int, string) {
	const prefix = "YoungDriverFee:"
	for _, item := range info {
		if !strings.Contains(item, prefix) {
			continue
		}
		// Split by comma to handle items that bundle multiple key:value:unit segments.
		for _, segment := range strings.Split(item, ",") {
			segment = strings.TrimSpace(segment)
			if !strings.HasPrefix(segment, prefix) {
				continue
			}

			// segment is now "YoungDriverFee:fee_amount:fee_currency"
			parts := strings.SplitN(segment, ":", 3)
			if len(parts) != 3 {
				return 0, ""
			}
			fee, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return 0, ""
			}
			return int(fee), parts[2]
		}
	}

	return 0, ""
}

// parseAddons returns a map of supplier name to addOn slice
func parseAddons(s flexSupplierDetails) []AddOn {
	addOns := make([]AddOn, 0, len(s.AvailableExtras))
	for _, e := range s.AvailableExtras {
		if e.MaxAmount == 0 {
			continue
		}
		addOns = append(addOns, AddOn{
			ID:              e.ExtraID,
			Name:            e.Name,
			Price:           int(e.Price),
			AllowedQuantity: e.MaxAmount,
			Period:          e.Period,
			Currency:        e.Currency,
		})
	}

	return addOns
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
func (f *Flex) getInsuranceExtraCost(days int) float64 {
	return float64(days) * f.erpDayCharge
}

// getPlans returns the list of plans for a given car
func (f *Flex) getPlans(c flexCar, dayCount int, supplierDetails flexSupplierDetails) []Plan {
	plans := make([]Plan, 0, len(c.Costs))
	for _, p := range c.Costs {
		planID, ok := flexProductMap[p.Product]
		if !ok {
			planID = 1
		}

		price := parseFloat(p.Price)
		if price == 0 {
			rlog.Warn("plan price is zero in CarAvailability response, skipping plan", "car_name", c.Name, "product", p.Product)
			continue
		}

		plans = append(plans, Plan{
			PlanID:                 planID,
			PlanName:               p.Product,
			Price:                  price,
			BrokerErpPrice:         parseFloat(c.ERP),
			ChargedErpPriceWithVat: f.getInsuranceExtraCost(dayCount),
			Info:                   c.Information,
			RateQualifier:          c.RateQualifier,
			SupplierName:           supplierDetails.Supplier,
			SupplierCode:           c.SupplierCode,
		})
	}

	return plans
}

// flexCarToBrokerCar converts a flexCar to a CarDetails struct for the broker response.
func flexCarToBrokerCar(c flexCar, supplierName string) CarDetails {
	seats, _ := strconv.Atoi(c.Passenger)
	doors, _ := strconv.Atoi(c.Doors)
	bags, _ := strconv.Atoi(c.Luggage)

	acriss := c.Code
	if len(acriss) > 4 {
		acriss = acriss[:4]
	}

	return CarDetails{
		Model:        normalizeModelName(c.Name),
		ImageURL:     c.URL,
		SupplierName: supplierName,
		CarType:      c.CarType,
		Acriss:       acriss,
		FullAcriss:   c.Code,
		HasAC:        strings.HasPrefix(c.IsAirCon, "Y"),
		IsAutoGear:   strings.HasPrefix(c.IsAutomatic, "Y"),
		IsElectric:   isElectric(c.Code),
		Seats:        seats,
		Doors:        doors,
		Bags:         bags,
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

func parseFloat(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

func (f *Flex) parseSupplierDetails(supplierDetails []flexSupplierDetails) []SupplierInfo {
	supplierInfo := make([]SupplierInfo, 0, len(supplierDetails))
	for _, s := range supplierDetails {
		supplierInfo = append(supplierInfo, SupplierInfo{
			Name:               s.Supplier,
			AddOns:             parseAddons(s),
			PlanInclusions:     f.parseInclusions(s),
			TermsAndConditions: f.parseTerms(s.Terms),
			PickupDetails:      f.parseStation(s.PickUpDetails),
			DropoffDetails:     f.parseStation(s.DropOffDetails),
		})
	}
	return supplierInfo
}

func (f *Flex) parseInclusions(s flexSupplierDetails) []string {
	var planInclusions []string
	inclusionsMap := make(map[string]struct{})
	for _, inc := range s.Inclusions {
		raw := strings.Split(inc.Inclusion, ";")
		for _, inclusion := range raw {
			if trimmed := strings.TrimSpace(inclusion); trimmed != "" {
				if _, exists := inclusionsMap[trimmed]; !exists {
					planInclusions = append(planInclusions, trimmed)
					inclusionsMap[trimmed] = struct{}{}
				}
			}
		}
	}

	return planInclusions
}

func (f *Flex) parseTerms(terms []flexTerms) []TermsAndConditionsItem {
	parsedTerms := make([]TermsAndConditionsItem, 0, len(terms))
	for _, t := range terms {
		parsedTerms = append(parsedTerms, TermsAndConditionsItem{
			Title:       t.Header,
			HtmlContent: t.Paragraph,
		})
	}

	return parsedTerms
}

func (f *Flex) parseStation(d flexLocationDetails) StationInfo {
	address := d.Address1
	if d.Address2 != "" {
		address += ", " + d.Address2
	}
	if d.Address3 != "" {
		address += ", " + d.Address3
	}
	return StationInfo{
		LocationInfo: d.LocationInformation,
		Address:      address,
		PhoneNumber:  d.Phone,
		OpeningHours: parseOpeningHours(d.OpeningHours),
	}
}

func parseOpeningHours(oh flexOpeningHours) []OpeningHoursItem {
	return []OpeningHoursItem{
		{Day: "Sunday", OpenTime: oh.SunOpen, CloseTime: oh.SunClose},
		{Day: "Monday", OpenTime: oh.MonOpen, CloseTime: oh.MonClose},
		{Day: "Tuesday", OpenTime: oh.TueOpen, CloseTime: oh.TueClose},
		{Day: "Wednesday", OpenTime: oh.WedOpen, CloseTime: oh.WedClose},
		{Day: "Thursday", OpenTime: oh.ThuOpen, CloseTime: oh.ThuClose},
		{Day: "Friday", OpenTime: oh.FriOpen, CloseTime: oh.FriClose},
		{Day: "Saturday", OpenTime: oh.SatOpen, CloseTime: oh.SatClose},
	}
}
