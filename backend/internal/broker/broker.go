package broker

// Name represents the name of a broker, such as "flex" or "hertz".
type Name string

const (
	BrokerFlex  Name = "flex"
	BrokerHertz Name = "hertz"
)

// secrets hold broker-specific secrets such as API credentials.
var secrets struct {
	// Flex secrets:
	flexAgentCode string
	flexPassword  string

	// Hertz secrets:
	hertzAgentDutyCode string
	hertzVendorNumber  string
	hertzCodeContext   string
}

const (
	flexBaseURL  = "http://www.flexibleautos.com/horizon/horizonxml.asmx"
	hertzBaseURL = "https://vv.xnet.hertz.com/DirectLinkWEB/handlers/DirectLinkHandler?id=ota2007a"
)

// LocationSearcher provides location listing capabilities for a broker.
type LocationSearcher interface {
	Name() Name
	GetLocationsPage(cursor string) (LocationPage, error)
}

// AvailabilitySearcher provides vehicle availability search capabilities for a broker.
type AvailabilitySearcher interface {
	Name() Name
	SearchAvailability(params SearchAvailabilityParams) (*AvailabilityResponse, error)
}

// AvailabilityResponse is the response structure for the AvailabilitySearcher interface
type AvailabilityResponse struct {
	AvailableVehicles []AvailableVehicle
	SuppliersInfo     []SupplierInfo
}

// Booker provides booking capabilities for a broker.
type Booker interface {
	Name() Name
	Book(p BookingParams) (BookingResponse, error)
}

// Canceler provides booking cancellation capabilities for a broker.
type Canceler interface {
	Name() Name
	Cancel(bookingID, lastName, supplierCode string) error
}

// VoucherProvider provides voucher generation capabilities for a broker.
type VoucherProvider interface {
	Name() Name
	GenerateVoucher(data *VoucherData) (string, error)
}

// VoucherData contains the details required to generate a voucher for a booking.
type VoucherData struct {
	ReservationNum     string
	CustomerName       string
	Telephone          string
	FlightNumber       string
	Issuer             string
	BillingAcc         string
	BookingReferenceID string
	Supplier           string
	PickupLoc          string
	PickupBranch       string
	PickupPhone        string
	PickupDate         string
	PickupTime         string
	DropoffLoc         string
	DropoffBranch      string
	DropoffPhone       string
	DropoffDate        string
	DropoffTime        string
	DropoffFee         string
	CarGroupDesc       string
	LeadModel          string
	RateCode           string
	ITCode             string
	Passengers         int
	Suitcases          int
	PrepaidIncludes    []string
	OptionalServices   []string
}

// Broker composes all broker capabilities into a single interface.
type Broker interface {
	LocationSearcher
	AvailabilitySearcher
	Booker
	Canceler
	VoucherProvider
}

// LocationPage represents a page of locations returned by a broker, including the list of locations and a cursor for the next page.
type LocationPage struct {
	Locations []Location
	NextPage  string
}

// Location represents a car rental location, including its ID, name, country, city, country code, and IATA code.
type Location struct {
	ID          string
	Name        string
	Country     string
	City        string
	CountryCode string
	Iata        string
}

// SearchAvailabilityParams represents the parameters required to search for available vehicles, such as pickup and return dates, location, and car preferences.
type SearchAvailabilityParams struct {
	CountryCode     string
	PickupLocation  string
	DropoffLocation string
	PickupDate      string
	DropoffDate     string
	PickupTime      string
	DropoffTime     string
	DriverAge       int
}

// AvailableVehicle represents a vehicle that is available for rent, including details about the car, the rental plans, add-ons, location details, and price details.
type AvailableVehicle struct {
	Broker          Name
	CarDetails      CarDetails
	Plans           []Plan
	LocationDetails LocationDetails
	PriceDetails    PriceDetails
}

// PriceDetails represents the pricing details of a rental, including the currency, drop charge, and drop charge currency.
type PriceDetails struct {
	Fees     Fees   `json:"fees"`
	Currency string `json:"currency"`
}

type Fees struct {
	DropCharge             int    `json:"dropCharge"`
	DropChargeCurrency     string `json:"dropChargeCurrency"`
	YoungDriverFee         int    `json:"youngDriverFee"`
	YoungDriverFeeCurrency string `json:"youngDriverFeeCurrency"`
}

// LocationDetails holds the locationType of a rental location
type LocationDetails struct {
	LocationType string `json:"locationType"`
}

// CarDetails represents the details of a car available for rent, including its ID, model, car group, image URL, supplier, car type, car size, ACRISS code, whether it has AC and auto gear, and the number of seats, bags, and doors.
type CarDetails struct {
	Model        string `json:"model"`
	CarGroup     string `json:"carGroup"`
	ImageURL     string `json:"imageUrl"`
	SupplierName string `json:"supplierName"`
	CarType      string `json:"carType"`
	Acriss       string `json:"acriss"`
	HasAC        bool   `json:"hasAC"`
	IsAutoGear   bool   `json:"isAutoGear"`
	IsElectric   bool   `json:"isElectric"`
	Seats        int    `json:"seats"`
	Bags         int    `json:"bags"`
	Doors        int    `json:"doors"`
}

// Plan represents a rental plan, including its ID, name, description, full price, discount, and other pricing details.
type Plan struct {
	PlanID                 int
	PlanName               string
	Price                  float64
	BrokerErpPrice         float64
	ChargedErpPriceWithVat float64
	Info                   []string
	RateQualifier          string
	SupplierName           string
	SupplierCode           string
}

// AddOn represents an additional service or product that can be added to a rental, including its ID, name, price, allowed quantity, and rental period.
type AddOn struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	Currency        string `json:"currency"`
	AllowedQuantity int    `json:"allowedQuantity"`
	Period          string `json:"period"`
}

// SupplierInfo represents information about a supplier, including its name, code, and any additional details.
type SupplierInfo struct {
	Name               string                   `json:"name"`
	AddOns             []AddOn                  `json:"addOns"`
	PlanInclusions     []string                 `json:"inclusions"`
	TermsAndConditions []TermsAndConditionsItem `json:"termsAndConditions"`
	PickupDetails      StationInfo              `json:"pickupDetails"`
	DropoffDetails     StationInfo              `json:"dropoffDetails"`
}

// StationInfo represents information about a rental station, including its location info, address, phone number, and opening hours.
type StationInfo struct {
	LocationInfo string             `json:"locationInfo"`
	Address      string             `json:"address"`
	PhoneNumber  string             `json:"phoneNumber"`
	OpeningHours []OpeningHoursItem `json:"openingHours"`
}

// OpeningHoursItem represents the opening hours for a specific day, including the day of the week and the opening and closing times.
type OpeningHoursItem struct {
	Day       string `json:"day"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
}

// TermsAndConditionsItem represents a single term or condition associated with a supplier, including its title and HTML content.
type TermsAndConditionsItem struct {
	Title       string `json:"title"`
	HtmlContent string `json:"htmlContent"`
}

// BookingResponse represents the response received after booking a rental
type BookingResponse struct {
	ConfirmationNumber string `json:"confirmationNumber"`
}

// BookingParams represents the parameters required to book a rental
type BookingParams struct {
	// snapshot plan details:
	RateQualifier   string
	SupplierCode    string
	Acriss          string
	PlanID          string
	PickupLocation  string
	DropoffLocation string

	// request:
	IncludeERP      bool
	SelectedAddOns  []SelectAddOn
	DriverTitle     string
	DriverFirstName string
	DriverLastName  string
	FlightNumber    string

	// snapshot columns:
	DriverAge   string
	PickupDate  string
	DropoffDate string
	PickupTime  string
	DropoffTime string
	CountryCode string
}

// SelectAddOn represents an add-on selected by the user during the booking process, including its ID and quantity.
type SelectAddOn struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}
