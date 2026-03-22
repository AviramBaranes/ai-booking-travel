package broker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Flex is a struct that implements the Broker interface for the Flex car rental service
type Flex struct {
	httpClient     *http.Client
	countriesCache []flexCountry
	erpDayCharge   int
}

const defaultTimeout = 5 * time.Minute

// NewFlexWithErpCfg creates a new instance of the Flex broker with a default HTTP client and timeout.
func NewFlexWithErpCfg(erpDayCharge int) Flex {
	return Flex{
		erpDayCharge: erpDayCharge,
		httpClient:   &http.Client{Timeout: defaultTimeout},
	}
}

// NewFlex creates a new instance of the Flex broker with a default HTTP client and timeout.
func NewFlex() Flex {
	return Flex{
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

// Name returns the name of the broker
func (f Flex) Name() Name {
	return BrokerFlex
}

// postForm sends a POST request to the specified endpoint of the Flex broker with the given form data.
func (f Flex) postForm(ep string, form url.Values) ([]byte, error) {
	form.Set("AgentCode", secrets.flexAgentCode)
	form.Set("Password", secrets.flexPassword)

	url := flexBaseURL + "/" + ep
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("flex %s create request: %w", ep, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json,application/xml,application/xml;q=0.9,*/*;q=0.8")

	res, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flex %s do request: %w", ep, err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("flex %s read response: %w", ep, err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("flex %s unexpected status %d: %s", ep, res.StatusCode, string(b))
	}

	return b, nil
}
