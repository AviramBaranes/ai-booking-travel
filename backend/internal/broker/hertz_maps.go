package broker

type hertzBrand string

const (
	hertzBrandDollar  hertzBrand = "ZR"
	hertzBrandThrifty hertzBrand = "ZT"
)

var hertzSupplierPlansMap = map[string]map[hertzBrand]struct {
	Standard string
	Gold     string
}{
	"US": {
		hertzBrandThrifty: {
			Standard: "IT1006124LGA",
			Gold:     "IT1006124DSL",
		},
		hertzBrandDollar: {
			Standard: "IT1006125LGA",
			Gold:     "IT1006125DSL",
		},
	},

	"CA": {
		hertzBrandThrifty: {
			Standard: "IT1006126",
			Gold:     "IT1006127",
		},
		hertzBrandDollar: {
			Standard: "IT1006128",
			Gold:     "IT1006129",
		},
	},
}

func mapBrandIDToSupplierName(brandID hertzBrand) string {
	switch brandID {
	case hertzBrandDollar:
		return "Dollar"
	case hertzBrandThrifty:
		return "Thrifty"
	default:
		return "Unknown"
	}
}

func mapCarTypeCodeToCarType(code string) string {
	switch code {
	case "1":
		return "car"
	case "2":
		return "van"
	case "3":
		return "suv"
	case "4":
		return "convertible"
	case "7":
		return "limousine"
	case "8":
		return "station wagon"
	case "9":
		return "pickup"
	case "10":
		return "motorhome"
	case "11":
		return "all-terrain"
	case "12":
		return "recreational"
	case "13":
		return "sport"
	case "14":
		return "special"
	case "15":
		return "pickup extended cab"
	case "16":
		return "regular cab pickup"
	case "17":
		return "special offer"
	case "18":
		return "coupe"
	case "19":
		return "monospace"
	case "20":
		return "2 wheel vehicle"
	case "21":
		return "roadster"
	case "22":
		return "crossover"
	case "23":
		return "commercial van/truck"
	default:
		return "unknown"
	}
}

var hertzUSAInclusionsMap = map[string][]string{
	"Standard": {
		"Unlimited mileage",
		"Super cover",
		"Loss damage waiver",
		"Liability insurance supplement",
		"Up to 4 additional Drivers age 25 & over",
		"Uninsured motor protection",
		"Homeland security fee",
		"Local tax",
		"Customer facility charge",
	},
	"Gold": {
		"Unlimited mileage",
		"Super cover",
		"Loss damage waiver",
		"Liability insurance supplement",
		"Up to 4 additional Drivers age 25 & over",
		"Uninsured motor protection",
		"Homeland security fee",
		"Customer facility charge",
		"Local tax",
		"First fuel tank option",
	},
}

var hertzCanadaInclusionsMap = map[string][]string{
	"Standard": {
		"Unlimited mileage",
		"Loss damage waiver",
		"Liability insurance supplement",
		"Uninsured motor protection",
		"Homeland security fee",
		"Local tax",
		"Customer facility charge",
	},
	"Gold": {
		"Additional Driver",
		"Unlimited mileage",
		"Loss damage waiver",
		"Liability insurance supplement",
		"Uninsured motor protection",
		"Homeland security fee",
		"Customer facility charge",
		"Local tax",
		"First fuel tank option",
	},
}

var hertzTermsMap = map[string][]string{
	"Gold": {
		"FUEL - FULL TO EMPTY",
		"EXCESS - ZERO",
		"CANCELLATION - NO FEES APPLY",
		"NO SHOW - NO FEES APPLY",
	},
	"Standard": {
		"FUEL - FULL TO FULL",
		"EXCESS - ZERO",
		"CANCELLATION - NO FEES APPLY",
		"NO SHOW - NO FEES APPLY",
	},
}
