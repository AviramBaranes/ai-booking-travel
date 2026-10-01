package broker

import "encoding/xml"

// avanceQuoteXML is a price-quote response. Every amount is an integer in minor units, so 9993 is
// 99.93. Failures arrive as a 200 carrying an <errors> block rather than as an HTTP error.
type avanceQuoteXML struct {
	XMLName    xml.Name            `xml:"pricequote"`
	ID         string              `xml:"id,attr"`
	Currency   string              `xml:"currency,attr"`
	FuelPolicy string              `xml:"fuelpolicy,attr"`
	OWPrepaid  bool                `xml:"owprepaid,attr"`
	Errors     []avanceErrorXML    `xml:"errors>error"`
	Categories []avanceCategoryXML `xml:"rates>category"`
}

type avanceErrorXML struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}

// avanceCategoryXML is one vehicle group within a quote.
//
// prepaidamount is what we collect at booking and what Avance invoices us; totalrate adds the
// charges the customer settles at the station (young driver, one-way, out of hours). excess and
// the options' exc and dep are quoted excluding VAT.
type avanceCategoryXML struct {
	Availability  string            `xml:"availability,attr"`
	Category      string            `xml:"cat,attr"`
	Code          string            `xml:"code,attr"`
	Acriss        string            `xml:"acriss,attr"`
	ImageURL      string            `xml:"imageurl,attr"`
	Model         string            `xml:"model,attr"`
	Pax           int               `xml:"pax,attr"`
	Doors         int               `xml:"doors,attr"`
	Suitcases     int               `xml:"suitcases,attr"`
	PrepaidAmount int               `xml:"prepaidamount,attr"`
	OneWayCharge  int               `xml:"onewaycharge,attr"`
	Excess        int               `xml:"excess,attr"`
	Unlimited     bool              `xml:"unlimited,attr"`
	IncludedKm    int               `xml:"IncKlm,attr"`
	Options       []avanceOptionXML `xml:"options>option"`
	Taxes         avanceTaxesXML    `xml:"taxes"`
}

// avanceOptionXML is an insurance waiver, extra or surcharge priced for the whole rental.
type avanceOptionXML struct {
	Code       string `xml:"code,attr"`
	Rate       int    `xml:"rate,attr"`
	FirstFree  bool   `xml:"firstfree,attr"`
	Inclusive  bool   `xml:"inclusive,attr"`
	Mandatory  bool   `xml:"mandatory,attr"`
	Prepaid    bool   `xml:"prepaid,attr"`
	ChargeType string `xml:"chargetype,attr"`
	// Exc and Dep are pointers because an absent attribute means something different from zero:
	// on some FDW0 groups the FDW waiver carries a deposit but no excess attribute at all.
	Exc *int `xml:"exc,attr"`
	Dep *int `xml:"dep,attr"`
}

type avanceTaxesXML struct {
	Tax1Rate int `xml:"tax1rate,attr"` // basis points: 2400 is 24%, always VAT in Avance's quotes
}

// avanceStationsXML is the stations response.
type avanceStationsXML struct {
	Stations []avanceStationXML `xml:"station"`
}

type avanceStationXML struct {
	Code string `xml:"code,attr"`
	Info struct {
		Type       string              `xml:"StationType,attr"`
		Address    string              `xml:"Address,attr"`
		ZipCode    string              `xml:"ZipCode,attr"`
		City       string              `xml:"City,attr"`
		Phone      string              `xml:"Phone,attr"`
		PickupInfo string              `xml:"PickupInstructions>PickupInfo"`
		Hours      []avanceWorkHourXML `xml:"OperationHours>WorkHour"`
	} `xml:"StationInformation"`
}

type avanceWorkHourXML struct {
	Day      string `xml:"Day,attr"`
	Closed   bool   `xml:"Closed,attr"`
	OpensAt  string `xml:"OpensAt,attr"`
	ClosesAt string `xml:"ClosesAt,attr"`
}

// avanceOptionsXML is the options response: the display name of every option code.
type avanceOptionsXML struct {
	Options []struct {
		Code  string `xml:"code,attr"`
		Name  string `xml:"name,attr"`
		Quant bool   `xml:"quant,attr"`
	} `xml:"option"`
}

// avanceReservationResXML is the new-res and cancel-res response. status is OK or an ERR code, and
// res-status is the reservation's own status: RES, REQ or CNC.
type avanceReservationResXML struct {
	Reservation struct {
		IRN       string `xml:"irn,attr"`
		Status    string `xml:"status,attr"`
		ResStatus string `xml:"res-status,attr"`
	} `xml:"reservation"`
	Errors []avanceErrorXML `xml:"errors>error"`
}
