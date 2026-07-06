package icount

import (
	"encoding/json"
	"fmt"
)

type CancelDocumentParams struct {
	DocType  string
	DocNum   string
	RefundCC bool
	Reason   string
}

// CancelDocument cancels an existing document in iCount.
func (i *Icount) CancelDocument(p CancelDocumentParams) (*ICountCancelDocumentResponse, error) {
	body, err := i.DoRequest(cancelDocEndpoint, ICountCancelDocumentRequest{
		DocType:  p.DocType,
		DocNum:   p.DocNum,
		RefundCC: p.RefundCC,
		Reason:   p.Reason,
	})
	if err != nil {
		return nil, fmt.Errorf("canceling document: %w", err)
	}

	var result ICountCancelDocumentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}
