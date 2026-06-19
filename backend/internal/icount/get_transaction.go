package icount

import (
	"encoding/json"
	"fmt"
)

func (i *Icount) GetTransactions(confCode string) (*ICountGetTransactionResponse, error) {
	body, err := i.DoRequest(getTransactionEndpoint, ICountGetTransactionRequest{
		ConfirmationCode: confCode,
	})

	if err != nil {
		return nil, fmt.Errorf("getting transactions: %w", err)
	}

	var result ICountGetTransactionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response for confirmation code:%s, json: %s, err: %w", confCode, string(body), err)
	}

	return &result, nil
}
