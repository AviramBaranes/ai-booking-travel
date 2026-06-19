package icount

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var secrets struct {
	IcountPassword string
}

type Icount struct {
	httpClient *http.Client
	cid        string
	user       string
}

// NewIcount initializes and returns an Icount client with the provided credentials.
// The password is read from the shared IcountPassword secret.
func NewIcount() *Icount {
	return &Icount{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type icountAPIEndPoint string

const (
	createDocEndpoint       icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/doc/create"
	fetchCurrenciesEndpoint icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/currency/get_rates"
	generatePaypageEndpoint icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/paypage/generate_sale"
	getTransactionEndpoint  icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/cc/transactions"
)

func (i *Icount) DoRequest(endpoint icountAPIEndPoint, requestBody any) ([]byte, error) {
	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("creating json: %w", err)
	}

	req, err := http.NewRequest("POST", string(endpoint), bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secrets.IcountPassword))

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
