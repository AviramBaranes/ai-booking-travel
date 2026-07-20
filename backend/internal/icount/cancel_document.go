package icount

import (
	"encoding/json"
	"fmt"

	"encore.dev/rlog"
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

	if !result.Status {
		rlog.Error("iCount cancel document failed", "reason", result.Reason, "docType", p.DocType, "docNum", p.DocNum)
		return nil, fmt.Errorf("iCount cancel document failed: %s", result.Reason)
	}

	return &result, nil
}
