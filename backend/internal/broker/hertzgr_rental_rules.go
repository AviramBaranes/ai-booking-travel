package broker

import "math"

// The rental rules below come from Autohellas' "Autohellas Greece 2026 Terms & Conditions" workbook,
// one set per brand (OTA vendor code). The availability response carries none of them. Update them
// with hertzGRTerms when Autohellas updates its terms.

// hertzGRVAT is the Greek VAT the workbook's "+ VAT" charges are subject to.
const hertzGRVAT = 0.24

// hertzGRFuelPolicy and hertzGROneDayMileage summarise the terms of every brand as plan terms.
const (
	hertzGRFuelPolicy    = "Fuel policy: return with the same level on 1-2 day rentals; from 3 days, the fuel in the tank is prepaid and the fuel returned is refunded"
	hertzGROneDayMileage = "Mileage: 150 km included on one-day rentals"
)

// hertzGRMinDriverAge is the youngest driver any brand rents to.
const hertzGRMinDriverAge = 19

// hertzGRExcess is the excess per car group, per brand. The workbook sets the collision (CDW)
// excess, the theft (TPC) excess and the deposit held at pickup to the same amount, in EUR.
var hertzGRExcess = map[string]map[string]int{
	"ZE": hertzGRGroupAmounts(map[int][]string{
		650:  {"A", "B", "J", "Z", "AE", "ZE"},
		800:  {"C", "D", "F", "K", "M", "P", "R", "T", "BE"},
		1000: {"E", "EM", "N", "NM", "O", "U", "V", "Y", "YL", "1", "5", "8", "DE", "FE", "EE"},
		1300: {"9", "SE", "SY", "9F", "V6"},
		3000: {"3", "YP", "9L"},
	}),
	"ZT": hertzGRGroupAmounts(map[int][]string{
		850:  {"A", "B", "J"},
		1000: {"C", "D", "F", "K", "M", "T"},
		1300: {"8", "E", "N", "NM", "O", "U", "YL"},
	}),
	"ZR": hertzGRGroupAmounts(map[int][]string{
		950:  {"A", "B", "J"},
		1100: {"C", "D", "F", "K", "M"},
		1500: {"O", "N", "NM"},
	}),
}

// hertzGRYoungDriverBand is the car groups a driver of the band's ages may rent, with each group's
// daily surcharge excluding VAT (0 for none). A group missing from the band is not available.
type hertzGRYoungDriverBand struct {
	minAge, maxAge int
	daily          map[string]float64
}

// hertzGRYoungDrivers are the bands below 25, per brand. From 25 every group is available without
// a surcharge.
var hertzGRYoungDrivers = map[string][]hertzGRYoungDriverBand{
	"ZE": {
		{19, 20, hertzGRGroupRates(map[float64][]string{10: {"A", "B", "J", "Z", "AE", "ZE", "V2"}})},
		{21, 24, hertzGRGroupRates(map[float64][]string{
			0:  {"A", "B", "Z", "J", "V2", "AE"},
			10: {"C", "D", "F", "K", "M", "P", "T", "R", "BE"},
			15: {"E", "EM", "N", "U", "V", "DE", "NM", "O", "V3", "V4", "Y", "EE", "FE"},
		})},
	},
	"ZT": {
		{19, 20, hertzGRGroupRates(map[float64][]string{10: {"A", "B", "J"}})},
		{21, 24, hertzGRGroupRates(map[float64][]string{
			0:  {"A", "B", "J"},
			10: {"C", "D", "F", "K", "M", "T"},
			15: {"E", "N", "U", "NM", "O"},
		})},
	},
	"ZR": {
		{19, 20, hertzGRGroupRates(map[float64][]string{10: {"A", "B", "J"}})},
		{21, 24, hertzGRGroupRates(map[float64][]string{
			0:  {"A", "B", "J"},
			10: {"C", "D", "F", "K", "M"},
			15: {"N", "NM", "O"},
		})},
	},
}

// hertzGRYoungDriverFee returns the young driver surcharge for the whole rental, VAT included, and
// whether the brand rents the car group to a driver of that age at all.
func hertzGRYoungDriverFee(brand, group string, age, days int) (float64, bool) {
	if age < hertzGRMinDriverAge {
		return 0, false
	}
	for _, band := range hertzGRYoungDrivers[brand] {
		if age < band.minAge || age > band.maxAge {
			continue
		}
		daily, ok := band.daily[group]
		if !ok {
			return 0, false
		}
		return math.Round(daily*(1+hertzGRVAT)*float64(days)*100) / 100, true
	}
	return 0, true
}

func hertzGRGroupAmounts(byAmount map[int][]string) map[string]int {
	out := make(map[string]int)
	for amount, groups := range byAmount {
		for _, g := range groups {
			out[g] = amount
		}
	}
	return out
}

func hertzGRGroupRates(byRate map[float64][]string) map[string]float64 {
	out := make(map[string]float64)
	for rate, groups := range byRate {
		for _, g := range groups {
			out[g] = rate
		}
	}
	return out
}
