package broker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Hertz is a struct that implements the Broker interface for the Hertz car rental service
type Hertz struct {
	usErpDayCharge int
	caErpDayCharge int
	httpClient     *http.Client
	r              io.Reader
}

// NewHertz creates a new instance of the Hertz broker with a default reader that is not initialized.
func NewHertz(usErpDayCharge, caErpDayCharge int) Hertz {
	return Hertz{
		usErpDayCharge: usErpDayCharge,
		caErpDayCharge: caErpDayCharge,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
	}
}

// NewHertzWithReader creates a new instance of the Hertz broker with the provided reader, which is used to read the CSV data for locations.
func NewHertzWithReader(r io.Reader) Hertz {
	return Hertz{r: r}
}

// Name returns the name of the broker
func (h Hertz) Name() Name {
	return BrokerHertz
}

// sendXMLRequest is a helper function that sends an XML request to the Hertz API.
func (h Hertz) sendXMLRequest(requestBody string) ([]byte, error) {

	req, err := http.NewRequest("POST", hertzBaseURL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return nil, fmt.Errorf("hertz create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hertz do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("hertz read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("hertz unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
