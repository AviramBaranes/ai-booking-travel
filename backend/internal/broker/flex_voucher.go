package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

func (f *Flex) GenerateVoucher(d *VoucherData) (string, error) {
	form := url.Values{}
	form.Set("FCHReference", d.BookingReferenceID)
	form.Set("SupplierReference", "")
	form.Set("FullVoucher", "True")
	form.Set("Language", "")
	form.Set("AdditionalParameters", "")

	body, err := f.postForm("GetVoucher", form)
	if err != nil {
		return "", err
	}

	var resp flexGetVoucherResponse

	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("flex GetVoucher unmarshal response: %w", err)
	}

	if resp.ReturnCode != 0 {
		return "", fmt.Errorf("GetVoucher API returned error code %d with message: %s", resp.ReturnCode, resp.ErrorMessage)
	}

	return resp.HtmlVoucher, nil
}
