package broker

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	hertzBookingSchemaLocation = "http://www.opentravel.org/OTA/2003/05 OTA_VehResRQ.xsd"
	hertzBookingVersion        = "1.008"
)

// Book creates a vehicle reservation with the Hertz API and returns the confirmation number.
func (h Hertz) Book(p BookingParams) (BookingResponse, error) {
	xmlReq, err := h.buildBookingRequest(p)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("building hertz booking request: %w", err)
	}

	body, err := h.sendXMLRequest(xmlReq)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("sending hertz booking request: %w", err)
	}

	var res hertzBookingResXML
	if err := xml.Unmarshal(body, &res); err != nil {
		return BookingResponse{}, fmt.Errorf("hertz Book unmarshal response: %w", err)
	}

	if len(res.Errors) > 0 {
		messages := make([]string, 0, len(res.Errors))
		for _, e := range res.Errors {
			messages = append(messages, e.ShortText)
		}
		return BookingResponse{}, fmt.Errorf("hertz booking failed: %s", strings.Join(messages, "; "))
	}

	return BookingResponse{
		ConfirmationNumber: res.ConfID.ID,
	}, nil
}

// buildBookingRequest constructs the XML request string for a Hertz vehicle reservation.
func (h Hertz) buildBookingRequest(p BookingParams) (string, error) {
	req := hertzBookingReq{
		XmlnsXsi:     hertzSearchAvailabilityXMLNSXSI,
		SchemaLoc:    hertzBookingSchemaLocation,
		Version:      hertzBookingVersion,
		SequenceNmbr: "1",
		POS:          h.buildPOSReqItem(hertzBrand(p.SupplierCode)),
		VehResRQCore: hertzBookingCore{
			Status: "All",
			VehRentalCore: hertzVehRentalCore{
				PickUpDateTime: h.formatDateTime(p.PickupDate, p.PickupTime),
				ReturnDateTime: h.formatDateTime(p.DropoffDate, p.DropoffTime),
				PickUpLocation: hertzLocation{LocationCode: p.PickupLocation},
				ReturnLocation: hertzLocation{LocationCode: p.DropoffLocation},
			},
			Customer: hertzCustomer{
				Primary: hertzPrimaryCustomer{
					PersonName: hertzPersonName{
						GivenName: p.DriverFirstName,
						Surname:   p.DriverLastName,
					},
				},
			},
		},
		VehResRQInfo: hertzBookingInfo{
			Reference: hertzBookingReference{
				Type: "16",
				ID:   p.RateQualifier,
			},
			WrittenConfInst: hertzWrittenConfInst{
				ConfirmInd: "false",
				LanguageID: "ENUS",
			},
		},
	}

	out, err := xml.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling hertz booking request: %w", err)
	}

	return string(out), nil
}
