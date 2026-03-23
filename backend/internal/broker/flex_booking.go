package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type flexBookingResponse struct {
	ReturnCode    int    `xml:"ReturnCode"`
	BookingNumber string `xml:"BookingNumber"`
	ErrorMessage  string `xml:"ErrorMessage"`
}

// Book books a rental using the provided booking parameters and returns a BookingResponse or an error if the booking fails.
func (f Flex) Book(p BookingParams) (BookingResponse, error) {
	form := url.Values{}
	form.Set("RateQualifier", p.RateQualifier)
	form.Set("ProductID", p.PlanID)
	form.Set("SupplierCode", p.SupplierCode)
	form.Set("CarType", p.Acriss)

	form.Set("DriverTitle", p.DriverTitle)
	form.Set("DriverInitial", p.DriverFirstName)
	form.Set("DriverLastName", p.DriverLastName)
	form.Set("DriverAge", p.DriverAge)

	form.Set("PickupLocationID", p.PickupLocation)
	form.Set("DropoffLocationID", p.DropoffLocation)
	form.Set("PickupDate", formatDate(p.PickupDate))
	form.Set("DropoffDate", formatDate(p.DropoffDate))
	form.Set("PickUpTime", p.PickupTime)
	form.Set("DropoffTime", p.DropoffTime)

	form.Set("FlightNumber", p.FlightNumber)
	form.Set("ERPRequired", getERPRequiredString(p.IncludeERP))
	form.Set("Extras", getExtrasString(p.SelectedAddOns))

	form.Set("NettCost", "")
	form.Set("AgentRef", "")
	form.Set("AdditionalParameters", "")

	body, err := f.postForm("PlaceBooking", form)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("failed to place booking: %w", err)
	}

	var resp flexBookingResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return BookingResponse{}, fmt.Errorf("flex PlaceBooking unmarshal response: %w", err)
	}

	if resp.ReturnCode != 0 {
		return BookingResponse{}, fmt.Errorf("booking failed with error: %s", resp.ErrorMessage)
	}

	return BookingResponse{
		ConfirmationNumber: resp.BookingNumber,
	}, nil
}

// getExtrasString converts a slice of SelectAddOn into a comma-separated string format required for the booking request.
func getExtrasString(addOns []SelectAddOn) string {
	extras := ""
	for _, addOn := range addOns {
		extras += fmt.Sprintf("%d:%d,", addOn.ID, addOn.Quantity)
	}
	if len(extras) > 0 {
		extras = extras[:len(extras)-1]
	}
	return extras
}

// getERPRequiredString converts true to "True" and false to "False" for the ERPRequired parameter in the booking request.
func getERPRequiredString(includeERP bool) string {
	if includeERP {
		return "True"
	}
	return "False"
}
