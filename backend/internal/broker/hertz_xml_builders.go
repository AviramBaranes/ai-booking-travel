package broker

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

const (
	hertzSearchAvailabilityXMLNS          = "http://www.opentravel.org/OTA/2003/05"
	hertzSearchAvailabilityXMLNSXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	hertzSearchAvailabilitySchemaLocation = "http://www.opentravel.org/OTA/2003/05 OTA_VehAvailRateRQ.xsd"
	hertzSearchAvailabilityVersion        = "1.008"

	hertzSearchAvailabilityISOCountry          = "IL"
	hertzSearchAvailabilityVendorNumberType    = "4"  //Type="4" - Indicates a unique Vendor Number (VN) will be used.
	hertzSearchAvailabilityBrandType           = "8"  //Type “8” – Identifies car rental Brand.
	hertzSearchAvailabilityConsumerProductCode = "CP" //This indicates that a consumer product (CP) code will be used.
)

func (h Hertz) buildSearchAvailabilityRequest(p SearchAvailabilityParams, brandID hertzBrand, planCode string, dayCount int) (string, error) {
	req := hertzSearchAvailabilityReq{
		XMLName:        xml.Name{Space: hertzSearchAvailabilityXMLNS, Local: "OTA_VehAvailRateRQ"},
		XmlnsXsi:       hertzSearchAvailabilityXMLNSXSI,
		SchemaLocation: hertzSearchAvailabilitySchemaLocation,
		Version:        hertzSearchAvailabilityVersion,
		POS:            h.buildPOSReqItem(brandID),
		VehAvailRQCore: hertzSearchAvailabilityCore{
			Status:        "All",
			VehRentalCore: h.buildRentalCoreReqItem(p, dayCount),
		},
		VehAvailRQInfo: hertzSearchAvailabilityInfo{
			TourInfo: hertzSearchAvailabilityTourInfo{
				TourNumber: planCode,
			},
		},
	}

	out, err := xml.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling xml hertz search availability request %w", err)
	}

	return string(out), nil
}

// formatDateTime formats the given date and time strings into the format required by the Hertz API, which is "YYYY-MM-DDTHH:MM:SS"
func (h Hertz) formatDateTime(date string, time string) string {
	return fmt.Sprintf("%sT%s:00", date, time)
}

// buildPOSReqItem creates a POS request item for the given brand ID and plan code, which is used in the search availability request to identify the brand and plan being requested.
func (h Hertz) buildPOSReqItem(brandID hertzBrand) hertzPOS {
	return hertzPOS{
		Source: []hertzSource{
			{
				ISOCountry:    hertzSearchAvailabilityISOCountry,
				AgentDutyCode: secrets.hertzAgentDutyCode,
				RequestorID: hertzRequestorID{
					Type: hertzSearchAvailabilityVendorNumberType,
					ID:   secrets.hertzVendorNumber,
					CompanyName: &hertzCompanyName{
						Code:        hertzSearchAvailabilityConsumerProductCode,
						CodeContext: secrets.hertzCodeContext,
					},
				},
			},
			{
				RequestorID: hertzRequestorID{
					Type: hertzSearchAvailabilityBrandType,
					ID:   string(brandID),
				},
			},
		},
	}
}

// buildRentalCoreReqItem creates a VehRentalCore request item
func (h Hertz) buildRentalCoreReqItem(p SearchAvailabilityParams, dayCount int) hertzVehRentalCore {
	return hertzVehRentalCore{
		PickUpDateTime:       h.formatDateTime(p.PickupDate, p.PickupTime),
		ReturnDateTime:       h.formatDateTime(p.DropoffDate, p.DropoffTime),
		MultIslandRentalDays: strconv.Itoa(dayCount),
		PickUpLocation: hertzLocation{
			LocationCode: p.PickupLocation,
		},
		ReturnLocation: hertzLocation{
			LocationCode: p.DropoffLocation,
		},
	}
}
