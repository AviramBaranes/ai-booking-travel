package icount

import (
	"encoding/json"
	"fmt"
)

func (i *Icount) FetchCurrencies() (*GetCurrenciesRatesResponse, error) {
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

	addCheckRateToCurrencies(&result)

	return &result, nil
}

const CHECK_CURRENCY_RATE_PERCENT = 1

func addCheckRateToCurrencies(r *GetCurrenciesRatesResponse) {
	for currency, rate := range r.Rates {
		if rate > 0 {
			r.Rates[currency] = rate * (1 + CHECK_CURRENCY_RATE_PERCENT/100.0)
		}
	}
}
