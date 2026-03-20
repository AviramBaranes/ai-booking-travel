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

type hertzCarAvailabilityXML struct {
	Cars []hertzCarXML    `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>VehAvails>VehAvail"`
	Info hertzLocationXML `xml:"VehAvailRSCore>VehVendorAvails>VehVendorAvail>Info"`
}

type hertzLocationXML struct {
	LocationDetails hertzLocationDetailsXML `xml:"LocationDetails"`
}

type hertzLocationDetailsXML struct {
	AdditionalInfo hertzAdditionalInfoXML `xml:"AdditionalInfo"`
}

type hertzAdditionalInfoXML struct {
	ParkLocation hertzParkLocationXML `xml:"ParkLocation"`
}

type hertzParkLocationXML struct {
	Location string `xml:"Location,attr"`
}

type hertzCarXML struct {
	Core hertzVehAvailCoreXML `xml:"VehAvailCore"`
}

type hertzVehAvailCoreXML struct {
	Reference  hertzReferenceXML  `xml:"Reference"`
	Vehicle    hertzVehicleXML    `xml:"Vehicle"`
	RentalRate hertzRentalRateXML `xml:"RentalRate"`
}

type hertzReferenceXML struct {
	ID string `xml:"ID,attr"`
}

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

type hertzVehTypeXML struct {
	CarType string `xml:"VehicleCategory,attr"`
}

type hertzVehMakeModelXML struct {
	Model string `xml:"Name,attr"`
}

type hertzRentalRateXML struct {
	Charges []hertzCharge `xml:"VehicleCharges>VehicleCharge"`
}

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

type hertzCarAvailabilityResponse struct {
	Cars         []hertzCar
	LocationType string
}

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

type hertzCharge struct {
	Purpose       string  `xml:"Purpose,attr"`
	TaxInclusive  bool    `xml:"TaxInclusive,attr"`
	GuaranteedInd bool    `xml:"GuaranteedInd,attr"`
	Amount        float64 `xml:"Amount,attr"`
	CurrencyCode  string  `xml:"CurrencyCode,attr"`
}
