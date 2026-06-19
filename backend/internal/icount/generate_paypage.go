package icount

import (
	"encoding/json"
	"fmt"
)

type GenerateIframeParams struct {
	PaypageID    int
	CurrencyCode string
	Sum          float64
	ClientName   string
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	Description  string
	OrderID      int64
	PageLang     string
	SuccessURL   string
	IpnURL       string
}

// GenerateIframe generates a paypage iframe using the iCount API.
func (i *Icount) GenerateIframe(p GenerateIframeParams) (*ICountGeneratePaypageSaleResponse, error) {
	icountReq := i.createRequest(p)

	body, err := i.DoRequest(generatePaypageEndpoint, icountReq)
	if err != nil {
		return nil, fmt.Errorf("generating paypage iframe: %w", err)
	}

	var result ICountGeneratePaypageSaleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}

func (i *Icount) createRequest(p GenerateIframeParams) ICountGeneratePaypageSaleRequest {
	return ICountGeneratePaypageSaleRequest{
		ClientName:   p.ClientName,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Email:        p.Email,
		Phone:        p.Phone,
		PaypageID:    p.PaypageID,
		Sum:          p.Sum,
		CurrencyCode: p.CurrencyCode,
		Description:  p.Description,
		IsIframe:     true,
		XOrderID:     p.OrderID,
		SuccessURL:   p.SuccessURL,
		IpnURL:       p.IpnURL,
	}
}
