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

type icount struct {
	httpClient *http.Client
	cid        string
	user       string
}

// NewIcount initializes and returns an icount client with the provided credentials.
// The password is read from the shared IcountPassword secret.
func NewIcount(cid, user string) icount {
	return icount{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cid:        cid,
		user:       user,
	}
}

type icountAPIEndPoint string

const (
	createDocEndpoint       icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/document/create"
	fetchCurrenciesEndpoint icountAPIEndPoint = "https://api.icount.co.il/api/v3.php/currency/get_rates"
)

func (i icount) DoRequest(endpoint icountAPIEndPoint, requestBody any) ([]byte, error) {
	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("creating json: %w", err)
	}

	req, err := http.NewRequest("POST", string(endpoint), bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
