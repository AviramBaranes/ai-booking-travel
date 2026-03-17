package broker

// flexCarAvailabilityResponse is the top-level XML response for the Flex CarAvailability API.
type flexCarAvailabilityResponse struct {
	Cars            []flexCar             `xml:"CarSet>Car"`
	SupplierDetails []flexSupplierDetails `xml:"SupplierInfo>Details"`
}

// flexCar represents a single available vehicle option with pricing and specs.
type flexCar struct {
	Status             string     `xml:"Status"`
	Category           string     `xml:"Category"`
	Name               string     `xml:"Name"`
	Code               string     `xml:"Code"`
	URL                string     `xml:"URL"`
	Luggage            string     `xml:"Luggage"`
	Passenger          string     `xml:"Passenger"`
	Doors              string     `xml:"Doors"`
	Currency           string     `xml:"Currency"`
	TotalCharge        float64    `xml:"TotalCharge"`
	Costs              []flexCost `xml:"Costs>Cost"`
	RateQualifier      string     `xml:"RateQualifier"`
	IsAirCon           string     `xml:"IsAirCon"`
	IsAutomatic        string     `xml:"IsAutomatic"`
	CarType            string     `xml:"CarType"`
	CarDescription     string     `xml:"CarDescription"`
	SupplierCode       string     `xml:"SupplierCode"`
	Supplier           string     `xml:"Supplier"`
	DropCharge         float64    `xml:"DropCharge"`
	DropChargeCurrency string     `xml:"DropChargeCurrency"`
	ERP                float64    `xml:"ERP"`
	Information        []string   `xml:"Information>string"`
}

// flexCost is a single cost line item (e.g. base rate, insurance) within a car's pricing breakdown.
type flexCost struct {
	Product string  `xml:"Product"`
	Price   float64 `xml:"Price"`
}

// flexSupplierDetails holds supplier-specific info: inclusions, policies, location details, and extras.
type flexSupplierDetails struct {
	Supplier        string              `xml:"Supplier"`
	Inclusions      []flexInclusion     `xml:"Inclusions>Inclusion"`
	FuelPolicy      string              `xml:"FuelPolicy"`
	ExcessPolicy    string              `xml:"ExcessPolicy"`
	PickUpDetails   flexLocationDetails `xml:"PickUpDetails"`
	DropOffDetails  flexLocationDetails `xml:"DropOffDetails"`
	AvailableExtras []flexExtraDetail   `xml:"AvailableExtras>ExtraDetails"`
	Terms           []flexTerms         `xml:"TermsAndConditions>TandCs"`
}

type flexTerms struct {
	Header    string `xml:"Header"`
	Paragraph string `xml:"Paragraph"`
}

// flexInclusion describes a single included product/coverage in the rental (e.g. CDW, theft protection).
type flexInclusion struct {
	Product   string `xml:"Product"`
	Inclusion string `xml:"Inclusion"`
}

// flexExtraDetail represents an optional add-on (e.g. GPS, child seat) available for purchase.
type flexExtraDetail struct {
	Name         string  `xml:"Name"`
	SupplierCode string  `xml:"SupplierCode"`
	ExtraID      int     `xml:"ExtraID"`
	Price        float64 `xml:"Price"`
	Currency     string  `xml:"Currency"`
	Period       string  `xml:"Period"`
	MaxAmount    int     `xml:"MaxAmount"`
	Information  string  `xml:"Information"`
}

// flexLocationDetails is the pick-up or drop-off location with address, phone, and opening hours.
type flexLocationDetails struct {
	LocationType        string           `xml:"LocationType"`
	LocationInformation string           `xml:"LocationInformation"`
	Address1            string           `xml:"Address1"`
	Address2            string           `xml:"Address2"`
	Address3            string           `xml:"Address3"`
	Phone               string           `xml:"PhoneNo"`
	OpeningHours        flexOpeningHours `xml:"OpeningHours"`
}

// flexOpeningHours holds the weekly opening/closing times for a location.
type flexOpeningHours struct {
	MonOpen  string `xml:"Mon_Open"`
	MonClose string `xml:"Mon_Close"`
	TueOpen  string `xml:"Tue_Open"`
	TueClose string `xml:"Tue_Close"`
	WedOpen  string `xml:"Wed_Open"`
	WedClose string `xml:"Wed_Close"`
	ThuOpen  string `xml:"Thu_Open"`
	ThuClose string `xml:"Thu_Close"`
	FriOpen  string `xml:"Fri_Open"`
	FriClose string `xml:"Fri_Close"`
	SatOpen  string `xml:"Sat_Open"`
	SatClose string `xml:"Sat_Close"`
	SunOpen  string `xml:"Sun_Open"`
	SunClose string `xml:"Sun_Close"`
}
