package broker

import "encoding/xml"

// OTA 2006A messages as Autohellas uses them. Only the elements and attributes we send or read are
// modelled; the guide and the parameter lists Autohellas sent describe the rest.

// hertzGRPOS identifies us on every message: our customer code number (CCN) as the requestor and
// the validation code (VC) as its company name.
type hertzGRPOS struct {
	Sources []hertzGRSource `xml:"Source"`
}

type hertzGRSource struct {
	ISOCountry  string             `xml:"ISOCountry,attr"`
	RequestorID hertzGRRequestorID `xml:"RequestorID"`
}

type hertzGRRequestorID struct {
	Type        string              `xml:"Type,attr"`
	ID          string              `xml:"ID,attr"`
	CompanyName *hertzGRCompanyName `xml:"CompanyName,omitempty"`
}

type hertzGRCompanyName struct {
	Code string `xml:"Code,attr"`
}

type hertzGRVehRentalCore struct {
	PickUpDateTime string                 `xml:"PickUpDateTime,attr"`
	ReturnDateTime string                 `xml:"ReturnDateTime,attr"`
	PickUpLocation hertzGRLocationCodeRef `xml:"PickUpLocation"`
	ReturnLocation hertzGRLocationCodeRef `xml:"ReturnLocation"`
}

type hertzGRLocationCodeRef struct {
	LocationCode string `xml:"LocationCode,attr"`
}

type hertzGRRateQualifier struct {
	RateQualifier string `xml:"RateQualifier,attr"`
}

// hertzGRVehAvailRateRQ quotes every car of one brand, the one the rate code selects.
type hertzGRVehAvailRateRQ struct {
	XMLName xml.Name              `xml:"OTA_VehAvailRateRQ"`
	Xmlns   string                `xml:"xmlns,attr"`
	Version string                `xml:"Version,attr"`
	POS     hertzGRPOS            `xml:"POS"`
	Core    hertzGRVehAvailRQCore `xml:"VehAvailRQCore"`
}

type hertzGRVehAvailRQCore struct {
	// Status "Available" asks for free-sell cars only, leaving out OnRequest ones.
	Status        string               `xml:"Status,attr"`
	RentalCore    hertzGRVehRentalCore `xml:"VehRentalCore"`
	DriverType    *hertzGRDriverType   `xml:"DriverType,omitempty"`
	RateQualifier hertzGRRateQualifier `xml:"RateQualifier"`
}

type hertzGRDriverType struct {
	Age int `xml:"Age,attr"`
}

type hertzGRVehAvailRateRS struct {
	Errors       []hertzGRError       `xml:"Errors>Error"`
	VendorAvails []hertzGRVendorAvail `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail"`
}

// hertzGRError is an OTA error. Autohellas puts the message in ShortText on some errors and in the
// element's text on others.
type hertzGRError struct {
	Type      string `xml:"Type,attr"`
	Code      string `xml:"Code,attr"`
	ShortText string `xml:"ShortText,attr"`
	Text      string `xml:",chardata"`
}

type hertzGRVendorAvail struct {
	Avails []hertzGRVehAvail `xml:"VehAvails>VehAvail"`
	Info   struct {
		ParkLocation struct {
			Location string `xml:"Location,attr"`
		} `xml:"LocationDetails>AdditionalInfo>ParkLocation"`
	} `xml:"Info"`
}

type hertzGRVehAvail struct {
	Core hertzGRVehAvailCore `xml:"VehAvailCore"`
	Info struct {
		PricedCoverages []hertzGRPricedCoverage `xml:"PricedCoverages>PricedCoverage"`
	} `xml:"VehAvailInfo"`
}

type hertzGRVehAvailCore struct {
	Status      string             `xml:"Status,attr"`
	Vehicle     hertzGRVehicle     `xml:"Vehicle"`
	RentalRate  hertzGRRentalRate  `xml:"RentalRate"`
	TotalCharge hertzGRTotalCharge `xml:"TotalCharge"`
	Fees        []hertzGRCharge    `xml:"Fees>Fee"`
}

type hertzGRVehicle struct {
	AirConditionInd   bool   `xml:"AirConditionInd,attr"`
	TransmissionType  string `xml:"TransmissionType,attr"`
	PassengerQuantity int    `xml:"PassengerQuantity,attr"`
	BaggageQuantity   int    `xml:"BaggageQuantity,attr"`
	// VendorCarType is Autohellas' car group, what a reservation asks for.
	VendorCarType string `xml:"VendorCarType,attr"`
	// Code is the ACRISS (SIPP) code.
	Code    string `xml:"Code,attr"`
	VehType struct {
		VehicleCategory string `xml:"VehicleCategory,attr"`
		DoorCount       int    `xml:"DoorCount,attr"`
	} `xml:"VehType"`
	MakeModel struct {
		Name string `xml:"Name,attr"`
	} `xml:"VehMakeModel"`
	PictureURL string `xml:"PictureURL"`
}

type hertzGRRentalRate struct {
	RateDistance struct {
		Unlimited    bool   `xml:"Unlimited,attr"`
		Quantity     int    `xml:"Quantity,attr"`
		DistUnitName string `xml:"DistUnitName,attr"`
	} `xml:"RateDistance"`
	VehicleCharges []hertzGRCharge `xml:"VehicleCharges>VehicleCharge"`
}

// hertzGRCharge is a charge, fee or coverage price. IncludedInRate marks what the rate code already
// pays for; IncludedInEstTotalInd what the grand total adds up.
type hertzGRCharge struct {
	Amount                float64 `xml:"Amount,attr"`
	CurrencyCode          string  `xml:"CurrencyCode,attr"`
	Description           string  `xml:"Description,attr"`
	Purpose               string  `xml:"Purpose,attr"`
	IncludedInRate        bool    `xml:"IncludedInRate,attr"`
	IncludedInEstTotalInd bool    `xml:"IncludedInEstTotalInd,attr"`
}

type hertzGRTotalCharge struct {
	RateTotalAmount      float64 `xml:"RateTotalAmount,attr"`
	EstimatedTotalAmount float64 `xml:"EstimatedTotalAmount,attr"`
	CurrencyCode         string  `xml:"CurrencyCode,attr"`
}

type hertzGRPricedCoverage struct {
	Coverage struct {
		Code string `xml:"Code,attr"`
	} `xml:"Coverage"`
	Charge hertzGRCharge `xml:"Charge"`
}

// hertzGRVehResRQ books one car: a car group of the brand the rate code selects.
type hertzGRVehResRQ struct {
	XMLName xml.Name             `xml:"OTA_VehResRQ"`
	Xmlns   string               `xml:"xmlns,attr"`
	Version string               `xml:"Version,attr"`
	POS     hertzGRPOS           `xml:"POS"`
	Core    hertzGRVehResRQCore  `xml:"VehResRQCore"`
	Info    *hertzGRVehResRQInfo `xml:"VehResRQInfo,omitempty"`
}

type hertzGRVehResRQCore struct {
	Status        string               `xml:"Status,attr"`
	RentalCore    hertzGRVehRentalCore `xml:"VehRentalCore"`
	Customer      hertzGRCustomer      `xml:"Customer"`
	VehPref       hertzGRVehPref       `xml:"VehPref"`
	DriverType    *hertzGRDriverType   `xml:"DriverType,omitempty"`
	RateQualifier hertzGRRateQualifier `xml:"RateQualifier"`
}

type hertzGRCustomer struct {
	PersonName hertzGRPersonName `xml:"Primary>PersonName"`
}

type hertzGRPersonName struct {
	GivenName string `xml:"GivenName,omitempty"`
	Surname   string `xml:"Surname"`
}

// hertzGRVehPref asks for a car group by Autohellas' own group code, the option Autohellas
// recommends.
type hertzGRVehPref struct {
	VendorCarType string `xml:"VendorCarType,attr"`
}

type hertzGRVehResRQInfo struct {
	ArrivalDetails *hertzGRArrivalDetails `xml:"ArrivalDetails,omitempty"`
}

// hertzGRArrivalDetails carries the flight number; transportation code 14 is air.
type hertzGRArrivalDetails struct {
	TransportationCode string `xml:"TransportationCode,attr"`
	Number             string `xml:"Number,attr"`
}

type hertzGRVehResRS struct {
	Errors []hertzGRError `xml:"Errors>Error"`
	Core   struct {
		// ReservationStatus is "Committed" for a confirmed booking and "On Hold" for one waiting
		// for Autohellas' reservations desk.
		ReservationStatus string `xml:"ReservationStatus,attr"`
		Segment           struct {
			ConfID      hertzGRUniqueID    `xml:"ConfID"`
			TotalCharge hertzGRTotalCharge `xml:"TotalCharge"`
			Fees        []hertzGRCharge    `xml:"Fees>Fee"`
		} `xml:"VehReservation>VehSegmentCore"`
	} `xml:"VehResRSCore"`
}

// hertzGRUniqueID identifies a reservation: type 14 is Autohellas' reservation number, type 16 our
// own reference.
type hertzGRUniqueID struct {
	Type string `xml:"Type,attr"`
	ID   string `xml:"ID,attr"`
}

type hertzGRVehCancelRQ struct {
	XMLName xml.Name `xml:"OTA_VehCancelRQ"`
	Xmlns   string   `xml:"xmlns,attr"`
	Version string   `xml:"Version,attr"`
	POS     hertzGRPOS
	Core    struct {
		CancelType string            `xml:"CancelType,attr"`
		UniqueID   hertzGRUniqueID   `xml:"UniqueID"`
		PersonName hertzGRPersonName `xml:"PersonName"`
	} `xml:"VehCancelRQCore"`
}

type hertzGRVehCancelRS struct {
	Errors []hertzGRError `xml:"Errors>Error"`
	Core   struct {
		CancelStatus string `xml:"CancelStatus,attr"`
	} `xml:"VehCancelRSCore"`
}
