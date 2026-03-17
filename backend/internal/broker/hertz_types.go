package broker

import (
	"encoding/xml"
)

type hertzSearchAvailabilityReq struct {
	XMLName        xml.Name                    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_VehAvailRateRQ"`
	XmlnsXsi       string                      `xml:"xmlns:xsi,attr"`
	SchemaLocation string                      `xml:"xsi:schemaLocation,attr"`
	Version        string                      `xml:"Version,attr"`
	POS            hertzPOS                    `xml:"POS"`
	VehAvailRQCore hertzSearchAvailabilityCore `xml:"VehAvailRQCore"`
	VehAvailRQInfo hertzSearchAvailabilityInfo `xml:"VehAvailRQInfo"`
}

type hertzPOS struct {
	Source []hertzSource `xml:"Source"`
}

type hertzSource struct {
	ISOCountry    string           `xml:"ISOCountry,attr,omitempty"`
	AgentDutyCode string           `xml:"AgentDutyCode,attr,omitempty"`
	RequestorID   hertzRequestorID `xml:"RequestorID"`
}

type hertzRequestorID struct {
	Type        string            `xml:"Type,attr"`
	ID          string            `xml:"ID,attr"`
	CompanyName *hertzCompanyName `xml:"CompanyName,omitempty"`
}

type hertzCompanyName struct {
	Code        string `xml:"Code,attr"`
	CodeContext string `xml:"CodeContext,attr"`
}

type hertzSearchAvailabilityCore struct {
	Status        string             `xml:"Status,attr"`
	VehRentalCore hertzVehRentalCore `xml:"VehRentalCore"`
}

type hertzVehRentalCore struct {
	PickUpDateTime       string        `xml:"PickUpDateTime,attr"`
	ReturnDateTime       string        `xml:"ReturnDateTime,attr"`
	MultIslandRentalDays string        `xml:"MultIslandRentalDays,attr"`
	PickUpLocation       hertzLocation `xml:"PickUpLocation"`
	ReturnLocation       hertzLocation `xml:"ReturnLocation"`
}

type hertzLocation struct {
	LocationCode string `xml:"LocationCode,attr"`
}

type hertzSearchAvailabilityInfo struct {
	TourInfo hertzSearchAvailabilityTourInfo `xml:"TourInfo"`
}

type hertzSearchAvailabilityTourInfo struct {
	TourNumber string `xml:"TourNumber,attr"`
}

type hertzCarAvailabilityResponse struct {
	Cars         []hertzCar        `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>VehAvails>VehAvail"`
	LocationInfo hertzLocationInfo `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>Info"`
}

type hertzLocationInfo struct {
	LocationDetails struct {
		AdditionalInfo struct {
			ParkLocation struct {
				LocationType string `xml:"Location,attr"`
			} `xml:"ParkLocation"`
		} `xml:"AdditionalInfo"`
	} `xml:"LocationDetails"`
}

type hertzCar struct {
	VehAvailCore hertzVehAvailCore `xml:"VehAvailCore"`
}

type hertzVehAvailCore struct {
	Reference struct {
		ID string `xml:"ID,attr"`
	} `xml:"Reference"`
	Vehicle struct {
		ImageURL         string `xml:"PictureURL"`
		Acriss           string `xml:"Code,attr"`
		HasAC            bool   `xml:"AirConditionInd,attr"`
		TransmissionType string `xml:"TransmissionType,attr"`
		Seats            int    `xml:"PassengerQuantity,attr"`
		Doors            int    `xml:"Doors,attr"`
		Bags             int    `xml:"BaggageQuantity,attr"`
		VehType          struct {
			CarType string `xml:"VehicleCategory,attr"`
		} `xml:"VehType"`
		VehMakeModel struct {
			Model string `xml:"Name,attr"`
		} `xml:"VehMakeModel"`
	} `xml:"Vehicle"`
	RentalRate struct {
		Charges []hertzCharge `xml:"VehicleCharges>VehicleCharge"`
	} `xml:"RentalRate"`
}

type hertzCharge struct {
	Purpose       string  `xml:"Purpose,attr"`
	TaxInclusive  bool    `xml:"TaxInclusive,attr"`
	GuaranteedInd bool    `xml:"GuaranteedInd,attr"`
	Amount        float64 `xml:"Amount,attr"`
	CurrencyCode  string  `xml:"CurrencyCode,attr"`
}
