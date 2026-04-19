package broker

import (
	"encoding/xml"
	"fmt"

	"encore.dev/rlog"
)

func (h Hertz) Cancel(bookingID, lastName, supplierCode string) error {
	xmlReq, err := h.buildCancelRequest(bookingID, lastName, supplierCode)
	if err != nil {
		return fmt.Errorf("building hertz cancel request: %w", err)
	}

	body, err := h.sendXMLRequest(xmlReq)
	if err != nil {
		return fmt.Errorf("sending hertz cancel request: %w", err)
	}

	var res hertzCancelResXML
	if err := xml.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("hertz Cancel unmarshal response: %w", err)
	}

	if err := h.parseErrors(res.Errors); err != nil {
		return fmt.Errorf("hertz cancellation errors: %w", err)
	}

	if warnings := h.parseWarnings(res.Warnings); warnings != "" {
		rlog.Warn("hertz cancellation warnings", warnings)
	}

	rlog.Info("cancellation successful", "booking_id", bookingID, "cancellation_id", fmt.Sprintf("UniqueID: %s", res.UniqueID.ID))

	return nil
}

// buildCancelRequest constructs the XML request string for a Hertz vehicle cancellation.
func (h Hertz) buildCancelRequest(bookingID, lastName, supplierCode string) (string, error) {
	req := hertzCancelReq{
		XMLName:   xml.Name{Space: hertzXMLNS, Local: "OTA_VehCancelRQ"},
		XmlnsXsi:  hertzXMLNSXSI,
		SchemaLoc: hertzCancelSchemaLocation,
		Version:   hertzXMLVersion,
		POS:       h.buildPOSReqItem(hertzBrand(supplierCode), false),
		VehCancelRQCore: hertzCancelCore{
			Type: "Book",
			UniqueID: hertzCancelUniqueID{
				Type: "14",
				ID:   bookingID,
			},
			PersonName: hertzCancelPersonName{
				Surname: lastName,
			},
		},
	}

	out, err := xml.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling hertz cancel request: %w", err)
	}

	return string(out), nil
}
