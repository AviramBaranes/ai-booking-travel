package broker

import (
	"encoding/xml"
)

// ----- Hertz Availability Types -----

// hertzSearchAvailabilityReq is the top-level OTA XML request envelope for vehicle availability.
type hertzSearchAvailabilityReq struct {
	XMLName        xml.Name                    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_VehAvailRateRQ"`
	XmlnsXsi       string                      `xml:"xmlns:xsi,attr"`
	SchemaLocation string                      `xml:"xsi:schemaLocation,attr"`
	Version        string                      `xml:"Version,attr"`
	POS            hertzPOS                    `xml:"POS"`
	VehAvailRQCore hertzSearchAvailabilityCore `xml:"VehAvailRQCore"`
	VehAvailRQInfo hertzSearchAvailabilityInfo `xml:"VehAvailRQInfo"`
}

// hertzPOS represents the Point of Sale element, identifying the requesting agent.
type hertzPOS struct {
	Source []hertzSource `xml:"Source"`
}

// hertzSource identifies the source of the request, including country and agent duty code.
type hertzSource struct {
	ISOCountry    string           `xml:"ISOCountry,attr,omitempty"`
	AgentDutyCode string           `xml:"AgentDutyCode,attr,omitempty"`
	RequestorID   hertzRequestorID `xml:"RequestorID"`
}

// hertzRequestorID carries the requestor type, ID, and optional company name.
type hertzRequestorID struct {
	Type        string            `xml:"Type,attr"`
	ID          string            `xml:"ID,attr"`
	CompanyName *hertzCompanyName `xml:"CompanyName,omitempty"`
}

// hertzCompanyName holds a company code and its context (e.g. IATA, ERSP).
type hertzCompanyName struct {
	Code        string `xml:"Code,attr"`
	CodeContext string `xml:"CodeContext,attr"`
}

// hertzSearchAvailabilityCore contains the rental core parameters such as pickup/return date and location.
type hertzSearchAvailabilityCore struct {
	Status        string             `xml:"Status,attr"`
	VehRentalCore hertzVehRentalCore `xml:"VehRentalCore"`
}

// hertzVehRentalCore holds pickup and return datetime/location details for the rental.
type hertzVehRentalCore struct {
	PickUpDateTime       string        `xml:"PickUpDateTime,attr"`
	ReturnDateTime       string        `xml:"ReturnDateTime,attr"`
	MultIslandRentalDays string        `xml:"MultIslandRentalDays,attr"`
	PickUpLocation       hertzLocation `xml:"PickUpLocation"`
	ReturnLocation       hertzLocation `xml:"ReturnLocation"`
}

// hertzLocation identifies a Hertz rental location by its location code.
type hertzLocation struct {
	LocationCode string `xml:"LocationCode,attr"`
}

// hertzSearchAvailabilityInfo carries supplementary request info such as tour details.
type hertzSearchAvailabilityInfo struct {
	TourInfo hertzSearchAvailabilityTourInfo `xml:"TourInfo"`
}

// hertzSearchAvailabilityTourInfo holds the tour number used for negotiated rates.
type hertzSearchAvailabilityTourInfo struct {
	TourNumber string `xml:"TourNumber,attr"`
}

// hertzCarAvailabilityXML is the parsed XML response containing available cars and location info.
type hertzCarAvailabilityXML struct {
	Errors []hertzResError  `xml:"Errors>Error"`
	Cars   []hertzCarXML    `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>VehAvails>VehAvail"`
	Info   hertzLocationXML `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>Info"`
}

// HasErrors reports whether the availability response contains any errors.
func (x hertzCarAvailabilityXML) HasErrors() bool {
	return len(x.Errors) > 0
}

// hertzLocationXML wraps the location details returned alongside availability results.
type hertzLocationXML struct {
	LocationDetails hertzLocationDetailsXML `xml:"LocationDetails"`
}

// hertzLocationDetailsXML holds additional info about a location, such as park type.
type hertzLocationDetailsXML struct {
	AdditionalInfo hertzAdditionalInfoXML `xml:"AdditionalInfo"`
}

// hertzAdditionalInfoXML contains supplementary location attributes like park location.
type hertzAdditionalInfoXML struct {
	ParkLocation hertzParkLocationXML `xml:"ParkLocation"`
}

// hertzParkLocationXML indicates the park/terminal location type (e.g. "Terminal", "Off Airport").
type hertzParkLocationXML struct {
	Location string `xml:"Location,attr"`
}

// hertzCarXML represents a single vehicle entry returned in the availability XML response.
type hertzCarXML struct {
	Core hertzVehAvailCoreXML `xml:"VehAvailCore"`
}

// hertzVehAvailCoreXML contains the reference, vehicle details, and rental rate for one availability record.
type hertzVehAvailCoreXML struct {
	Reference  hertzReferenceXML  `xml:"Reference"`
	Vehicle    hertzVehicleXML    `xml:"Vehicle"`
	RentalRate hertzRentalRateXML `xml:"RentalRate"`
}

// hertzReferenceXML holds the unique reference ID for a vehicle availability record.
type hertzReferenceXML struct {
	ID string `xml:"ID,attr"`
}

// hertzVehicleXML contains vehicle attributes as returned by the Hertz XML API.
type hertzVehicleXML struct {
	ImageURL         string               `xml:"PictureURL"`
	Acriss           string               `xml:"Code,attr"`
	HasAC            bool                 `xml:"AirConditionInd,attr"`
	TransmissionType string               `xml:"TransmissionType,attr"`
	Seats            int                  `xml:"PassengerQuantity,attr"`
	Doors            int                  `xml:"Doors,attr"`
	Bags             int                  `xml:"BaggageQuantity,attr"`
	VehType          hertzVehTypeXML      `xml:"VehType"`
	VehMakeModel     hertzVehMakeModelXML `xml:"VehMakeModel"`
}

// hertzVehTypeXML holds the vehicle category code (e.g. car, van, SUV).
type hertzVehTypeXML struct {
	CarType string `xml:"VehicleCategory,attr"`
}

// hertzVehMakeModelXML holds the make/model display name of the vehicle.
type hertzVehMakeModelXML struct {
	Model string `xml:"Name,attr"`
}

// hertzRentalRateXML contains the list of charges associated with a vehicle's rental rate.
type hertzRentalRateXML struct {
	Charges []hertzCharge `xml:"VehicleCharges>VehicleCharge"`
}

// toResponse converts the raw XML availability data into a hertzCarAvailabilityResponse.
func (x hertzCarAvailabilityXML) toResponse() hertzCarAvailabilityResponse {
	resp := hertzCarAvailabilityResponse{
		LocationType: x.Info.LocationDetails.AdditionalInfo.ParkLocation.Location,
		Cars:         make([]hertzCar, 0, len(x.Cars)),
	}
	for _, c := range x.Cars {
		resp.Cars = append(resp.Cars, hertzCar{
			ID:               c.Core.Reference.ID,
			Model:            c.Core.Vehicle.VehMakeModel.Model,
			ImageURL:         c.Core.Vehicle.ImageURL,
			CarType:          c.Core.Vehicle.VehType.CarType,
			Acriss:           c.Core.Vehicle.Acriss,
			HasAC:            c.Core.Vehicle.HasAC,
			TransmissionType: c.Core.Vehicle.TransmissionType,
			Seats:            c.Core.Vehicle.Seats,
			Doors:            c.Core.Vehicle.Doors,
			Bags:             c.Core.Vehicle.Bags,
			Charges:          c.Core.RentalRate.Charges,
		})
	}
	return resp
}

// hertzCarAvailabilityResponse holds the parsed list of available cars and the location type.
type hertzCarAvailabilityResponse struct {
	Cars         []hertzCar
	LocationType string
}

// hertzCar holds the parsed details for a single available vehicle.
type hertzCar struct {
	ID               string
	Model            string
	ImageURL         string
	CarType          string
	Acriss           string
	HasAC            bool
	TransmissionType string
	Seats            int
	Doors            int
	Bags             int
	Charges          []hertzCharge
}

// hertzCharge represents a single charge line item within a vehicle's rental rate.
type hertzCharge struct {
	Purpose       string  `xml:"Purpose,attr"`
	TaxInclusive  bool    `xml:"TaxInclusive,attr"`
	GuaranteedInd bool    `xml:"GuaranteedInd,attr"`
	Amount        float64 `xml:"Amount,attr"`
	CurrencyCode  string  `xml:"CurrencyCode,attr"`
}

// ----- Hertz Booking Types -----

// hertzBookingReq is the top-level OTA XML request envelope for creating a vehicle reservation.
type hertzBookingReq struct {
	XMLName      xml.Name         `xml:"http://www.opentravel.org/OTA/2003/05 OTA_VehResRQ"`
	XmlnsXsi     string           `xml:"xmlns:xsi,attr"`
	SchemaLoc    string           `xml:"xsi:schemaLocation,attr"`
	Version      string           `xml:"Version,attr"`
	SequenceNmbr string           `xml:"SequenceNmbr,attr"`
	POS          hertzPOS         `xml:"POS"`
	VehResRQCore hertzBookingCore `xml:"VehResRQCore"`
	VehResRQInfo hertzBookingInfo `xml:"VehResRQInfo"`
}

// hertzBookingCore holds the rental core, customer, and optional special equipment for a booking request.
type hertzBookingCore struct {
	Status            string                  `xml:"Status,attr"`
	VehRentalCore     hertzVehRentalCore      `xml:"VehRentalCore"`
	Customer          hertzCustomer           `xml:"Customer"`
	SpecialEquipPrefs *hertzSpecialEquipPrefs `xml:"SpecialEquipPrefs,omitempty"`
}

// hertzCustomer contains the primary customer details for the reservation.
type hertzCustomer struct {
	Primary hertzPrimaryCustomer `xml:"Primary"`
}

// hertzPrimaryCustomer holds the person name of the primary renter.
type hertzPrimaryCustomer struct {
	PersonName hertzPersonName `xml:"PersonName"`
}

// hertzPersonName represents a person's given (first) and family (last) name.
type hertzPersonName struct {
	GivenName string `xml:"GivenName"`
	Surname   string `xml:"Surname"`
}

// hertzSpecialEquipPrefs is a placeholder for optional special equipment preferences.
type hertzSpecialEquipPrefs struct {
	// TODO: define special equipment preference fields as needed
}

// hertzBookingInfo carries supplementary booking request info such as reference and confirmation preferences.
type hertzBookingInfo struct {
	Reference       hertzBookingReference `xml:"Reference"`
	WrittenConfInst hertzWrittenConfInst  `xml:"WrittenConfInst"`
}

// hertzBookingReference holds a reference type and ID used to link back to an availability record.
type hertzBookingReference struct {
	Type string `xml:"Type,attr"`
	ID   string `xml:"ID,attr"`
}

// hertzWrittenConfInst controls written confirmation preferences (e.g. language, whether to send).
type hertzWrittenConfInst struct {
	ConfirmInd string `xml:"ConfirmInd,attr"`
	LanguageID string `xml:"LanguageID,attr"`
}

// hertzBookingResXML is the parsed XML response for a booking request, covering both success and error cases.
type hertzBookingResXML struct {
	Errors   []hertzResError   `xml:"Errors>Error"`
	Warnings []hertzResWarning `xml:"Warnings>Warning"`
	ConfID   hertzConfID       `xml:"VehResRSCore>VehReservation>VehSegmentCore>ConfID"`
}

// hertzResError represents a single error element returned in a failed booking response.
type hertzResError struct {
	Type      string `xml:"Type,attr"`
	ShortText string `xml:"ShortText,attr"`
	Code      string `xml:"Code,attr"`
	RecordID  string `xml:"RecordID,attr"`
}

// hertzResWarning represents a single warning element returned in a booking response.
type hertzResWarning struct {
	Type      string `xml:"Type,attr"`
	ShortText string `xml:"ShortText,attr"`
	RecordID  string `xml:"RecordID,attr"`
}

// hertzConfID holds the confirmation type and ID returned upon a successful booking.
type hertzConfID struct {
	Type string `xml:"Type,attr"`
	ID   string `xml:"ID,attr"`
}
