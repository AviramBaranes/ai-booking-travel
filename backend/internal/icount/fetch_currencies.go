package icount

import (
	"encoding/json"
	"fmt"
)

func (i *icount) FetchCurrencies() (*GetCurrenciesRatesResponse, error) {
	reqBody := GetCurrenciesRatesRequest{
		CID:  i.cid,
		User: i.user,
		Pass: secrets.IcountPassword,
	}

	body, err := i.DoRequest(fetchCurrenciesEndpoint, reqBody)
	if err != nil {
		return nil, err
	}

	var result GetCurrenciesRatesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("icount error: %s", result.Reason)
	}

	return &result, nil

}
