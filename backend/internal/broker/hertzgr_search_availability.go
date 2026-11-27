package broker

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"encore.app/internal/pricing"
	"encore.dev/rlog"
)

const (
	// hertzGRISOCountry is the POS country we send. The guide's samples send "IE"; Autohellas has
	// not said whether the value changes anything for our CCN.
	hertzGRISOCountry = "IL"

	// hertzGRAvailable is the only status we sell. OnRequest cars wait up to 24 hours for
	// Autohellas' reservations desk, which a flow that takes the customer's card up front cannot
	// offer. The request already asks for available cars only; the check guards the response.
	hertzGRAvailable = "Available"

	// hertzGRPayOnArrival is the fee purpose for the amount settled at the station.
	hertzGRPayOnArrival = "23"
	// hertzGRDropOff is the vehicle charge purpose for a one-way fee.
	hertzGRDropOff = "2"
	// hertzGRRentalCharge is the vehicle charge purpose for the rental itself.
	hertzGRRentalCharge = "1"
)

// hertzGRBrandByLetter maps the 4th letter of a rate code, which selects the brand Autohellas quotes,
// onto the brand's OTA vendor code, the code stations are tagged with.
var hertzGRBrandByLetter = map[byte]string{'H': "ZE", 'T': "ZT", 'D': "ZR", 'F': "FF"}

// hertzGRBrandNames are the supplier names customers see, spelled like Flex's so the existing logos
// in the CMS apply.
var hertzGRBrandNames = map[string]string{"ZE": "Hertz", "ZT": "Thrifty", "ZR": "Dollar", "FF": "FireFly"}

// hertzGRNoCarsCodes are the OTA error codes that mean the search has nothing to sell, rather than
// that something is wrong: no cars or no rate for the dates, a pickup or return outside opening
// hours, a one-way rental across a border, or a rental too long for the rate.
var hertzGRNoCarsCodes = map[string]bool{
	"214": true, "221": true, "222": true, "226": true, "227": true,
	"263": true, "280": true, "301": true, "599": true,
}

// hertzGRQuote is one request of a search: a rate code, quoted at the stations its brand works from.
type hertzGRQuote struct {
	rate    HertzGRRate
	planID  int
	brand   string
	pickup  string
	dropoff string
}

// SearchAvailability quotes every rate code concurrently, each at the stations of its brand, and
// merges them into one vehicle per brand and car group, with a plan per rate code.
func (h *HertzGR) SearchAvailability(p SearchAvailabilityParams) (*AvailabilityResponse, error) {
	if len(h.rates) == 0 {
		rlog.Warn("hertz gr has no rate codes configured, returning no vehicles")
		return &AvailabilityResponse{}, nil
	}

	quotes := hertzGRQuotes(h.rates, p)
	if len(quotes) == 0 {
		return &AvailabilityResponse{}, nil
	}

	days, err := CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rental days count: %w", err)
	}

	results := make([]*hertzGRVehAvailRateRS, len(quotes))
	errs := make([]error, len(quotes))
	var (
		stations    map[string]hertzGRStation
		fleetImages map[string]string
	)

	var wg sync.WaitGroup
	for i, q := range quotes {
		wg.Add(1)
		go func(i int, q hertzGRQuote) {
			defer wg.Done()
			results[i], errs[i] = h.quote(p, q)
		}(i, q)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		stations = h.stations()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fleetImages = h.fleetImages()
	}()
	wg.Wait()

	// A brand failing still leaves the other brands to sell, so only fail outright when no quote
	// came back.
	succeeded := 0
	for i, err := range errs {
		if err != nil {
			rlog.Warn("hertz gr availability quote failed", "rateCode", quotes[i].rate.Code, "pickup", quotes[i].pickup, "error", err)
			continue
		}
		succeeded++
	}
	if succeeded == 0 {
		return nil, errors.Join(errs...)
	}

	return buildHertzGRAvailability(p, days, quotes, results, stations, fleetImages), nil
}

// hertzGRQuotes pairs each rate code with the pickup and return stations of its brand. A brand
// without a station at either end is not quoted.
func hertzGRQuotes(rates []HertzGRRate, p SearchAvailabilityParams) []hertzGRQuote {
	quotes := make([]hertzGRQuote, 0, len(rates))

	for i, rate := range rates {
		brand := hertzGRRateBrand(rate.Code)
		if brand == "" {
			rlog.Error("hertz gr rate code has no known brand letter, skipping it", "rateCode", rate.Code)
			continue
		}

		pickup := hertzGRStationFor(p.PickupStations, brand)
		dropoff := hertzGRStationFor(p.DropoffStations, brand)
		if pickup == "" || dropoff == "" {
			continue
		}

		quotes = append(quotes, hertzGRQuote{rate: rate, planID: i + 1, brand: brand, pickup: pickup, dropoff: dropoff})
	}

	return quotes
}

// hertzGRRateBrand returns the brand a rate code quotes, read from its 4th letter.
func hertzGRRateBrand(code string) string {
	if len(code) < 4 {
		return ""
	}
	return hertzGRBrandByLetter[code[3]]
}

// hertzGRStationFor returns the station serving brand among a location's stations.
func hertzGRStationFor(stations []Station, brand string) string {
	for _, s := range stations {
		if slices.Contains(s.Brands, brand) {
			return s.Code
		}
	}
	return ""
}

// quote fetches one rate code's availability. A response that only says there is nothing to sell
// comes back empty rather than as an error.
func (h *HertzGR) quote(p SearchAvailabilityParams, q hertzGRQuote) (*hertzGRVehAvailRateRS, error) {
	rq := hertzGRVehAvailRateRQ{
		Xmlns:   otaNamespace,
		Version: "1.0",
		POS:     hertzGRRequestPOS(),
		Core: hertzGRVehAvailRQCore{
			Status: hertzGRAvailable,
			RentalCore: hertzGRVehRentalCore{
				PickUpDateTime: p.PickupDate + "T" + p.PickupTime + ":00",
				ReturnDateTime: p.DropoffDate + "T" + p.DropoffTime + ":00",
				PickUpLocation: hertzGRLocationCodeRef{LocationCode: q.pickup},
				ReturnLocation: hertzGRLocationCodeRef{LocationCode: q.dropoff},
			},
			RateQualifier: hertzGRRateQualifier{RateQualifier: q.rate.Code},
		},
	}
	if p.DriverAge > 0 {
		rq.Core.DriverType = &hertzGRDriverType{Age: p.DriverAge}
	}

	var rs hertzGRVehAvailRateRS
	if err := h.call("VehicleAvailRate", rq, &rs); err != nil {
		return nil, err
	}

	if len(rs.Errors) > 0 {
		e := rs.Errors[0]
		if hertzGRNoCarsCodes[e.Code] {
			return &hertzGRVehAvailRateRS{}, nil
		}
		return nil, fmt.Errorf("hertz gr availability returned %s: %s", e.Code, e.message())
	}

	return &rs, nil
}

// hertzGRRequestPOS is the POS every request carries: our CCN with the VC as its company name.
func hertzGRRequestPOS() hertzGRPOS {
	return hertzGRPOS{Sources: []hertzGRSource{{
		ISOCountry: hertzGRISOCountry,
		RequestorID: hertzGRRequestorID{
			Type:        "4",
			ID:          hertzGRCreds.ccn,
			CompanyName: &hertzGRCompanyName{Code: hertzGRCreds.vc},
		},
	}}}
}

func (e hertzGRError) message() string {
	if s := strings.TrimSpace(e.ShortText); s != "" {
		return s
	}
	return strings.TrimSpace(e.Text)
}

// buildHertzGRAvailability merges the quotes into vehicles. results is indexed like quotes, and a
// nil entry is a quote that failed. stations holds the feed's details by station code, and
// fleetImages each car group's picture, for cars whose response carries none.
func buildHertzGRAvailability(p SearchAvailabilityParams, days int, quotes []hertzGRQuote, results []*hertzGRVehAvailRateRS, stations map[string]hertzGRStation, fleetImages map[string]string) *AvailabilityResponse {
	vehicles := make(map[string]*AvailableVehicle)
	var order []string
	suppliers := make(map[string]*SupplierInfo)
	var supplierOrder []string

	for i, q := range quotes {
		rs := results[i]
		if rs == nil {
			continue
		}

		name := hertzGRBrandNames[q.brand]
		var planInclusions []string

		for _, va := range rs.VendorAvails {
			locationType := hertzGRLocationType(va.Info.ParkLocation.Location)

			for ai := range va.Avails {
				a := &va.Avails[ai]
				if a.Core.Status != hertzGRAvailable {
					continue
				}

				// A brand rents some groups to young drivers for a surcharge and others not at all.
				youngDriverFee := 0.0
				if p.DriverAge > 0 {
					fee, allowed := hertzGRYoungDriverFee(q.brand, strings.TrimSpace(a.Core.Vehicle.VendorCarType), p.DriverAge, days)
					if !allowed {
						continue
					}
					youngDriverFee = fee
				}

				plan, fees, ok := hertzGRPlan(q, a, days, youngDriverFee)
				if !ok {
					continue
				}

				key := q.brand + ":" + a.Core.Vehicle.VendorCarType
				v, ok := vehicles[key]
				if !ok {
					details := hertzGRCarDetails(&a.Core.Vehicle, name)
					if details.ImageURL == "" {
						details.ImageURL = fleetImages[hertzGRFleetKey(p.CountryCode, q.brand, details.CarGroup)]
					}
					v = &AvailableVehicle{
						Broker:          BrokerHertzGR,
						CarDetails:      details,
						LocationDetails: LocationDetails{LocationType: locationType},
						PriceDetails:    PriceDetails{Currency: a.Core.TotalCharge.CurrencyCode, Fees: fees},
					}
					vehicles[key] = v
					order = append(order, key)
				}
				v.Plans = append(v.Plans, plan)

				// What a rate code includes is the same on every car, so the first one describes it.
				if planInclusions == nil {
					planInclusions = hertzGRInclusions(a)
				}
			}
		}

		sp, ok := suppliers[q.brand]
		if !ok {
			sp = &SupplierInfo{
				Name:               name,
				TermsAndConditions: hertzGRTerms[q.brand],
				PickupDetails:      stations[q.pickup].stationInfo(p.PickupDate),
				DropoffDetails:     stations[q.dropoff].stationInfo(p.DropoffDate),
			}
			suppliers[q.brand] = sp
			supplierOrder = append(supplierOrder, q.brand)
		}
		if planInclusions != nil {
			sp.Inclusions = append(sp.Inclusions, Inclusions{ProductName: q.rate.PlanName, ProductInclusions: planInclusions})
		}
	}

	if len(order) == 0 {
		return &AvailabilityResponse{}
	}

	out := make([]AvailableVehicle, 0, len(order))
	for _, key := range order {
		out = append(out, *vehicles[key])
	}

	// One entry per brand, named like its plans' SupplierName: plans whose supplier has no matching
	// entry are dropped before pricing.
	infos := make([]SupplierInfo, 0, len(supplierOrder))
	for _, brand := range supplierOrder {
		infos = append(infos, *suppliers[brand])
	}

	return &AvailabilityResponse{AvailableVehicles: out, SuppliersInfo: infos}
}

// hertzGRPlan prices one car under one rate code. youngDriverFee is the surcharge for the whole
// rental, paid at the station. It returns false for a car with nothing to charge, which a malformed
// response could produce.
func hertzGRPlan(q hertzGRQuote, a *hertzGRVehAvail, days int, youngDriverFee float64) (Plan, Fees, bool) {
	total := a.Core.TotalCharge
	currency := total.CurrencyCode

	// The grand total includes what the customer pays at the station, such as a one-way fee; we
	// collect the rest.
	payOnArrival := 0.0
	for _, f := range a.Core.Fees {
		if f.Purpose == hertzGRPayOnArrival {
			payOnArrival += f.Amount
		}
	}
	price := total.EstimatedTotalAmount - payOnArrival
	if price <= 0 {
		rlog.Warn("hertz gr car has no prepaid amount, skipping it", "rateCode", q.rate.Code, "group", a.Core.Vehicle.VendorCarType)
		return Plan{}, Fees{}, false
	}

	var fees Fees
	dropOff := 0.0
	for _, c := range a.Core.RentalRate.VehicleCharges {
		if c.Purpose == hertzGRDropOff && !c.IncludedInRate {
			dropOff += c.Amount
		}
	}
	if dropOff > 0 {
		fees.DropCharge = pricing.RoundToInt(dropOff)
		fees.DropChargeCurrency = currency
	}
	if youngDriverFee > 0 {
		fees.YoungDriverFee = pricing.RoundToInt(youngDriverFee)
		fees.YoungDriverFeeCurrency = currency
	}

	// The excess and the deposit come from Autohellas' terms; the response carries neither.
	group := strings.TrimSpace(a.Core.Vehicle.VendorCarType)
	excess := hertzGRExcess[q.brand][group]

	var info []string
	// The fees tell the customer only about a one-way fee, so anything else due at the station is
	// spelled out.
	if payOnArrival > 0 && pricing.RoundToInt(payOnArrival) != fees.DropCharge {
		info = append(info, fmt.Sprintf("Pay on arrival: %.2f %s", payOnArrival, currency))
	}
	info = append(info, hertzGRFuelPolicy)
	if excess > 0 {
		info = append(info, fmt.Sprintf("Excess: %d %s", excess, currency))
	}
	if days == 1 {
		info = append(info, hertzGROneDayMileage)
	}

	return Plan{
		PlanID:                 q.planID,
		PlanName:               q.rate.PlanName,
		Price:                  price,
		BrokerErpPrice:         0,
		ChargedErpPriceWithVat: float64(days) * q.rate.ErpDayCharge,
		Info:                   info,
		RateQualifier:          hertzGRRateQualifierFor(q.rate.Code, a.Core.Vehicle.VendorCarType),
		SupplierName:           hertzGRBrandNames[q.brand],
		SupplierCode:           q.brand,
		PickupLocationCode:     q.pickup,
		DropoffLocationCode:    q.dropoff,
		Deposit:                excess,
		DepositCurrency:        currency,
		Excess:                 excess,
		ExcessCurrency:         currency,
		TheftExcess:            excess,
		TheftExcessCurrency:    currency,
	}, fees, true
}

// hertzGRRateQualifierFor identifies a plan at booking. Booking finds a plan by rate qualifier,
// supplier code and plan id, and a rate code alone is shared by every car of its brand, so the car
// group is appended. Booking needs both halves anyway: the rate code to price the reservation and
// the group to reserve.
func hertzGRRateQualifierFor(rateCode, group string) string {
	return rateCode + ":" + group
}

func hertzGRCarDetails(v *hertzGRVehicle, supplierName string) CarDetails {
	acriss := strings.TrimSpace(v.Code)
	short := acriss
	if len(short) > 4 {
		short = short[:4]
	}

	return CarDetails{
		Model:        normalizeModelName(strings.TrimSpace(v.MakeModel.Name)),
		CarGroup:     strings.TrimSpace(v.VendorCarType),
		ImageURL:     hertzGRImageURL(v.PictureURL),
		SupplierName: supplierName,
		CarType:      mapCarTypeCodeToCarType(v.VehType.VehicleCategory),
		Acriss:       short,
		FullAcriss:   acriss,
		HasAC:        v.AirConditionInd || acrissHasAC(acriss),
		IsAutoGear:   v.TransmissionType == "Automatic" || isAcrissShowsAutoGear(acriss),
		IsElectric:   isElectric(acriss),
		Seats:        v.PassengerQuantity,
		Bags:         v.BaggageQuantity,
		Doors:        v.VehType.DoorCount,
	}
}

// hertzGRImageURL serves car pictures over https: some come from Autohellas' own image host over
// http, which the site cannot load.
func hertzGRImageURL(u string) string {
	u = strings.TrimSpace(u)
	if rest, ok := strings.CutPrefix(u, "http://"); ok {
		return "https://" + rest
	}
	return u
}

// hertzGRLocationType maps where the car park is (OTA's Vehicle Where At Facility code) onto the
// values the results page badges.
func hertzGRLocationType(parkLocation string) string {
	switch parkLocation {
	case "1":
		return "Airport"
	case "2", "3":
		return "Shuttle"
	default:
		return "City"
	}
}

// hertzGRInclusions lists what a rate code includes: the charges and coverages marked as included
// in the rate, other than the rental itself, and the mileage.
func hertzGRInclusions(a *hertzGRVehAvail) []string {
	incs := make([]string, 0, 4)

	for _, c := range a.Core.RentalRate.VehicleCharges {
		if c.IncludedInRate && c.Purpose != hertzGRRentalCharge {
			if d := strings.TrimSpace(c.Description); d != "" {
				incs = append(incs, d)
			}
		}
	}
	for _, pc := range a.Info.PricedCoverages {
		if !pc.Charge.IncludedInRate {
			continue
		}
		if d := strings.TrimSpace(pc.Charge.Description); d != "" {
			incs = append(incs, d)
		} else if code := strings.TrimSpace(pc.Coverage.Code); code != "" {
			incs = append(incs, code)
		}
	}

	rd := a.Core.RentalRate.RateDistance
	if rd.Unlimited {
		incs = append(incs, "Unlimited mileage")
	} else if rd.Quantity > 0 {
		incs = append(incs, fmt.Sprintf("%d %s included", rd.Quantity, strings.ToLower(rd.DistUnitName)))
	}

	return incs
}
