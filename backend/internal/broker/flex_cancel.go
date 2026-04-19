package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

func (f Flex) Cancel(bookingID, _, _1 string) error {
	form := url.Values{}
	form.Set("FCHReference", bookingID)
	form.Set("Language", "UK")
	form.Set("IgnorePayment", "")
	form.Set("AdditionalParameters", "")

	body, err := f.postForm("CancelBooking", form)
	if err != nil {
		return err
	}

	var resp flexCancelResponse

	if err := xml.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("flex CancelBooking unmarshal response: %w", err)
	}

	if resp.ReturnCode != 0 {
		return fmt.Errorf("CancelBooking API returned error code %d with message: %s", resp.ReturnCode, resp.ErrorMessage)
	}

	return nil
}
