package broker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"encore.app/internal/pricing"
	"encore.dev/rlog"
)

const (
	avanceSupplierName = "Avance"
	avanceSupplierCode = "AVANCE"

	// avanceAvailable is the only availability we sell. ONREQUEST results — one-way rentals,
	// out-of-hours pickups, and some groups at some stations — wait for manual confirmation by
	// Avance, which a booking flow that takes the customer's card up front cannot offer.
	avanceAvailable = "AVAILABLE"

	// avanceMaxQuantity caps the extras Avance lets us book in quantity. Their terms allow at most
	// three additional drivers.
	avanceMaxQuantity = 3
)

// The two products (CDPs) Avance activated for us.
const (
	avanceCDPStandard   = "VCH"
	avanceCDPZeroExcess = "FDW0"
)

// avanceProduct is one Wheelsys product (CDP), presented as one plan on every vehicle.
type avanceProduct struct {
	cdp      string
	planID   int
	planName string
}

// avanceProducts maps each product to its plan. The plan names reuse the Standard/Gold pair the
// frontend already translates.
var avanceProducts = []avanceProduct{
	{cdp: avanceCDPStandard, planID: 1, planName: "Standard"},
	{cdp: avanceCDPZeroExcess, planID: 2, planName: "Gold"},
}

// avanceAddOnIDs gives the extras we sell a numeric id: broker.AddOn ids are ints and Avance uses
// string codes. The 9000 block keeps clear of Flex's supplier-assigned ids, since both brokers'
// add-ons are looked up by id in the same CMS image gallery. Every other option, including the FDW
// waiver on the standard product, is deliberately not offered.
var avanceAddOnIDs = map[string]int{"CS": 9001, "BS": 9002, "SNC": 9003, "ADD": 9004}

// avanceAddOnCodes fixes the order add-ons are listed in.
var avanceAddOnCodes = []string{"CS", "BS", "SNC", "ADD"}

var avanceFuelPolicies = map[string]string{
	"SL": "Fuel policy: return with the same level",
	"FF": "Fuel policy: full to full",
	"FE": "Fuel policy: full to empty",
	"FH": "Fuel policy: full to half",
}

// SearchAvailability quotes both products concurrently and merges them into one vehicle per car
// group, each carrying a Standard and a Gold plan.
func (a *Avance) SearchAvailability(p SearchAvailabilityParams) (*AvailabilityResponse, error) {
	days, err := CalculateDaysCount(p.PickupDate, p.PickupTime, p.DropoffDate, p.DropoffTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rental days count: %w", err)
	}

	quotes := make([]*avanceQuoteXML, len(avanceProducts))
	errs := make([]error, len(avanceProducts))
	var ref avanceReference

	var wg sync.WaitGroup
	for i, prod := range avanceProducts {
		wg.Add(1)
		go func(i int, prod avanceProduct) {
			defer wg.Done()
			quotes[i], errs[i] = a.quote(p, prod.cdp)
		}(i, prod)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ref = a.reference()
	}()
	wg.Wait()

	// One product failing still leaves a sellable plan on every car, so only fail outright when
	// neither quote came back.
	succeeded := 0
	for i, err := range errs {
		if err != nil {
			rlog.Warn("avance price-quote failed", "cdp", avanceProducts[i].cdp, "error", err)
			continue
		}
		succeeded++
	}
	if succeeded == 0 {
		return nil, errors.Join(errs...)
	}

	return buildAvanceAvailability(p, days, a.erpDayCharges, quotes, ref), nil
}

// quote fetches and parses one product's price quote.
func (a *Avance) quote(p SearchAvailabilityParams, cdp string) (*avanceQuoteXML, error) {
	q := url.Values{}
	q.Set("DATE_FROM", formatDate(p.PickupDate))
	q.Set("TIME_FROM", p.PickupTime)
	q.Set("DATE_TO", formatDate(p.DropoffDate))
	q.Set("TIME_TO", p.DropoffTime)
	q.Set("PICKUP_STATION", p.PickupLocation)
	q.Set("RETURN_STATION", p.DropoffLocation)
	q.Set("CDP", cdp)
	if p.DriverAge > 0 {
		// Wheelsys applies Avance's age rules itself: groups a young driver may not take are
		// dropped, and where they may, the young driver fee comes back as mandatory.
		q.Set("DriverAge", strconv.Itoa(p.DriverAge))
	}

	body, err := a.get("price-quote", q)
	if err != nil {
		return nil, err
	}

	return parseAvanceQuote(body)
}

func parseAvanceQuote(body []byte) (*avanceQuoteXML, error) {
	var resp avanceQuoteXML
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("avance price-quote unmarshal response: %w", err)
	}

	if len(resp.Errors) > 0 {
		e := resp.Errors[0]
		return nil, fmt.Errorf("avance price-quote returned %s: %s", e.Code, strings.TrimSpace(e.Message))
	}

	return &resp, nil
}

// buildAvanceAvailability merges the per-product quotes into vehicles. quotes is indexed like
// avanceProducts, and a nil entry is a product whose quote failed. erpDayCharges is keyed by CDP.
func buildAvanceAvailability(p SearchAvailabilityParams, days int, erpDayCharges map[string]float64, quotes []*avanceQuoteXML, ref avanceReference) *AvailabilityResponse {
	vehicles := make(map[string]*AvailableVehicle)
	var (
		order      []string
		inclusions []Inclusions
		addOns     []AddOn
	)

	pickup := ref.stations[p.PickupLocation]

	for i, prod := range avanceProducts {
		q := quotes[i]
		if q == nil {
			continue
		}

		var productInclusions []string
		for ci := range q.Categories {
			c := &q.Categories[ci]
			if c.Availability != avanceAvailable {
				continue
			}

			v, ok := vehicles[c.Code]
			if !ok {
				v = &AvailableVehicle{
					Broker:          BrokerAvance,
					CarDetails:      avanceCarDetails(c),
					LocationDetails: LocationDetails{LocationType: avanceLocationType(pickup.stationType)},
					PriceDetails:    PriceDetails{Currency: q.Currency, Fees: avanceFees(q, c)},
				}
				vehicles[c.Code] = v
				order = append(order, c.Code)
			}
			v.Plans = append(v.Plans, avancePlan(prod, q, c, days, erpDayCharges[prod.cdp]))

			// What a product includes is the same on every group, so the first one describes it.
			if productInclusions == nil {
				productInclusions = avanceInclusions(c, ref)
			}
			if addOns == nil {
				addOns = avanceAddOns(q, c, ref)
			}
		}

		if productInclusions != nil {
			inclusions = append(inclusions, Inclusions{ProductName: prod.planName, ProductInclusions: productInclusions})
		}
	}

	if len(order) == 0 {
		return &AvailabilityResponse{}
	}

	out := make([]AvailableVehicle, 0, len(order))
	for _, code := range order {
		out = append(out, *vehicles[code])
	}

	return &AvailabilityResponse{
		AvailableVehicles: out,
		// Exactly one entry, named like every plan's SupplierName: plans whose supplier has no
		// matching entry are dropped before pricing.
		SuppliersInfo: []SupplierInfo{{
			Name:               avanceSupplierName,
			AddOns:             addOns,
			Inclusions:         inclusions,
			TermsAndConditions: []TermsAndConditionsItem{},
			PickupDetails:      pickup.info,
			DropoffDetails:     ref.stations[p.DropoffLocation].info,
		}},
	}
}

func avancePlan(prod avanceProduct, q *avanceQuoteXML, c *avanceCategoryXML, days int, erpDayCharge float64) Plan {
	excess := avanceWithVAT(c.Excess, c)

	// The theft waiver keeps its own excess even on the zero-excess product, which only waives
	// collision damage.
	theftExcess := 0
	if thw := c.option("THW"); thw != nil && thw.Exc != nil {
		theftExcess = avanceWithVAT(*thw.Exc, c)
	}

	return Plan{
		PlanID:   prod.planID,
		PlanName: prod.planName,
		// prepaidamount is what we collect at booking and what Avance invoices us, VAT inclusive.
		// The rest of totalrate is settled by the customer at the station.
		Price:                  avanceMinorToUnits(c.PrepaidAmount),
		BrokerErpPrice:         0,
		ChargedErpPriceWithVat: float64(days) * erpDayCharge,
		Info:                   avanceInfo(q, c, excess, theftExcess),
		RateQualifier:          avanceRateQualifier(q.ID, c.Code),
		SupplierName:           avanceSupplierName,
		SupplierCode:           avanceSupplierCode,
		Deposit:                avanceWithVAT(avanceDeposit(c), c),
		DepositCurrency:        q.Currency,
		Excess:                 excess,
		ExcessCurrency:         q.Currency,
		TheftExcess:            theftExcess,
		TheftExcessCurrency:    q.Currency,
	}
}

// avanceRateQualifier identifies a plan at booking. Booking finds a plan by rate qualifier,
// supplier code and plan id, and a quote id alone is shared by every car in the quote, so the car
// group is appended. Booking needs both halves anyway: the quote id to lock the price and the group
// to reserve. The separator is a colon because group codes contain hyphens, as in E-GQ.
func avanceRateQualifier(quoteID, group string) string {
	return quoteID + ":" + group
}

// avanceDeposit returns the deposit, excluding VAT and in minor units, held by the waiver that sets
// the category's excess. It needs a fallback because on some FDW0 groups the FDW waiver carries a
// deposit but no excess attribute; the fallback skips zero deposits, which the theft waiver reports
// on some FDW0 groups.
func avanceDeposit(c *avanceCategoryXML) int {
	fallback := 0
	for _, o := range c.Options {
		if !o.Inclusive || !o.Mandatory || o.Dep == nil {
			continue
		}
		if o.Exc != nil && *o.Exc == c.Excess {
			return *o.Dep
		}
		if fallback == 0 && *o.Dep > 0 {
			fallback = *o.Dep
		}
	}
	return fallback
}

// avanceFees returns the charges the customer settles at the station. A charge already inside
// prepaidamount is left out, or the customer would be told to pay it twice.
func avanceFees(q *avanceQuoteXML, c *avanceCategoryXML) Fees {
	var f Fees

	if c.OneWayCharge > 0 && !q.OWPrepaid {
		f.DropCharge = pricing.RoundToInt(avanceMinorToUnits(c.OneWayCharge))
		f.DropChargeCurrency = q.Currency
	}

	if ydr := c.option("YDR"); ydr != nil && ydr.Mandatory && !ydr.Prepaid && ydr.Rate > 0 {
		f.YoungDriverFee = pricing.RoundToInt(avanceMinorToUnits(ydr.Rate))
		f.YoungDriverFeeCurrency = q.Currency
	}

	return f
}

func avanceCarDetails(c *avanceCategoryXML) CarDetails {
	acriss := strings.TrimSpace(c.Acriss)
	short := acriss
	if len(short) > 4 {
		short = short[:4]
	}

	return CarDetails{
		Model:        normalizeModelName(strings.TrimSpace(c.Model)),
		CarGroup:     c.Code,
		ImageURL:     c.ImageURL,
		SupplierName: avanceSupplierName,
		CarType:      c.Category,
		Acriss:       short,
		FullAcriss:   acriss,
		HasAC:        acrissHasAC(acriss),
		IsAutoGear:   isAcrissShowsAutoGear(acriss),
		IsElectric:   isElectric(acriss),
		Seats:        c.Pax,
		// Avance reports small bags as zero on every group; suitcases is the real luggage count.
		Bags:  c.Suitcases,
		Doors: c.Doors,
	}
}

// avanceLocationType maps Avance's station type onto the values the results page badges.
func avanceLocationType(stationType string) string {
	if _, ok := avanceAirportStationTypes[stationType]; ok {
		return "Airport"
	}
	if stationType == "Airport Shuttle" {
		return "Shuttle"
	}
	return stationType
}

// avanceInfo lists a plan's terms. The search response carries no excess of its own, so the figures
// are spelled out here, VAT included.
func avanceInfo(q *avanceQuoteXML, c *avanceCategoryXML, excess, theftExcess int) []string {
	info := make([]string, 0, 4)

	if fp, ok := avanceFuelPolicies[q.FuelPolicy]; ok {
		info = append(info, fp)
	}

	info = append(info, fmt.Sprintf("Excess: %d %s", excess, q.Currency))
	if theftExcess != excess {
		info = append(info, fmt.Sprintf("Theft excess: %d %s", theftExcess, q.Currency))
	}

	// Avance's terms for the zero-excess program, which the API does not express.
	if fdw := c.option("FDW"); fdw != nil && fdw.Inclusive {
		info = append(info, "Damage to the underside, interior, wheels and glass is not covered")
	}

	return info
}

// avanceInclusions lists what a product includes, named the way Avance names its options.
func avanceInclusions(c *avanceCategoryXML, ref avanceReference) []string {
	incs := make([]string, 0, 6)

	for _, o := range c.Options {
		if o.Inclusive && o.ChargeType == "I" {
			incs = append(incs, avanceOptionName(o.Code, ref))
		}
	}

	if c.Unlimited {
		incs = append(incs, "Unlimited mileage")
	} else if c.IncludedKm > 0 {
		incs = append(incs, fmt.Sprintf("%d km included", c.IncludedKm))
	}

	if add := c.option(avanceAdditionalDriver); add != nil && add.FirstFree {
		incs = append(incs, "First additional driver free")
	}

	return incs
}

// avanceAddOns returns the extras we sell. Their rates are whole-rental totals already capped by
// Avance, the same on every group and product, and settled at the station.
func avanceAddOns(q *avanceQuoteXML, c *avanceCategoryXML, ref avanceReference) []AddOn {
	addOns := make([]AddOn, 0, len(avanceAddOnCodes))

	for _, code := range avanceAddOnCodes {
		o := c.option(code)
		if o == nil {
			continue
		}

		qty := 1
		if ref.options[code].quant {
			qty = avanceMaxQuantity
		}

		// The add-on counts only paid drivers, which booking relies on to add the free first one.
		// A quote without the free driver would be mis-booked, so it isn't offered at all.
		if code == avanceAdditionalDriver {
			if !o.FirstFree {
				continue
			}
			qty = avanceMaxQuantity - 1
		}

		addOns = append(addOns, AddOn{
			ID:              avanceAddOnIDs[code],
			Name:            avanceOptionName(code, ref),
			Price:           pricing.RoundToInt(avanceMinorToUnits(o.Rate)),
			Currency:        q.Currency,
			AllowedQuantity: qty,
			Period:          "Per Rental",
		})
	}

	return addOns
}

func avanceOptionName(code string, ref avanceReference) string {
	if o, ok := ref.options[code]; ok && o.name != "" {
		return o.name
	}
	return code
}

func (c *avanceCategoryXML) option(code string) *avanceOptionXML {
	for i := range c.Options {
		if c.Options[i].Code == code {
			return &c.Options[i]
		}
	}
	return nil
}

// avanceWithVAT converts an amount quoted excluding VAT, in minor units, into whole currency units
// including VAT. Avance quotes excess and deposits net of VAT, while the customer faces the gross.
func avanceWithVAT(minor int, c *avanceCategoryXML) int {
	gross := float64(minor) * (1 + float64(c.Taxes.Tax1Rate)/10000)
	return pricing.RoundToInt(gross / 100)
}

func avanceMinorToUnits(minor int) float64 {
	return float64(minor) / 100
}
