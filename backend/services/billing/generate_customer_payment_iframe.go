package billing

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/icount"
	"encore.app/internal/lang"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/billing/db"
	"encore.app/services/booking"
	"encore.dev"
	"encore.dev/beta/auth"
	"encore.dev/rlog"
)

const CustomerPaymentIPNGatewayPath = "/customer-payment-ipn-gateway"

type GenerateCustomerPaymentIframeParams struct {
	Phone           string               `json:"phone" validate:"required,israeli_phone"`
	FirstName       string               `json:"firstName" validate:"required,notblank"`
	LastName        string               `json:"lastName" validate:"required,notblank"`
	Email           string               `json:"email" validate:"required,email"`
	SnapshotID      int64                `json:"snapshotId" validate:"required"`
	RateQualifier   string               `json:"rateQualifier" validate:"required"`
	SupplierCode    string               `json:"supplierCode" validate:"required"`
	PlanID          string               `json:"planId" validate:"required"`
	IncludeERP      bool                 `json:"includeERP"`
	SelectedAddOns  []broker.SelectAddOn `json:"selectedAddOns"`
	DriverTitle     string               `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Mrs' 'Ms' 'Miss' 'Dr'"`
	DriverFirstName string               `json:"driverFirstName" validate:"required,uppercase_only"`
	DriverLastName  string               `json:"driverLastName" validate:"required,uppercase_only"`
	FlightNumber    *string              `json:"flightNumber" encore:"optional" validate:"omitempty"`
}

func (p GenerateCustomerPaymentIframeParams) Validate() error {
	return validation.ValidateStruct(p)
}

type GenerateCustomerPaymentIframeResponse struct {
	URL              string `json:"url"`
	PendingPaymentID int64  `json:"pendingPaymentId"`
}

// GenerateCustomerPaymentIframe creates a pending payment record and returns an iCount iframe URL for the customer to complete payment.
// Pricing and booking details are derived server-side from the snapshot.
// If the request is authenticated as a customer, the user ID is stored on the pending record.
// Requests authenticated under any other role are rejected.
// Unauthenticated requests proceed without a user ID (guest checkout).
//
// encore:api public path=/billing/generate-customer-payment-iframe method=POST
func (s *Service) GenerateCustomerPaymentIframe(ctx context.Context, p *GenerateCustomerPaymentIframeParams) (*GenerateCustomerPaymentIframeResponse, error) {
	var userID *int64
	if auth.Data() != nil {
		authData := auth.Data().(*accounts.AuthData)
		if authData.Role != accounts.UserRoleCustomer {
			return nil, api_errors.ErrUnauthorized
		}
		uid := authData.UserID
		userID = &uid
	}

	planSummary, err := booking.GetSnapshotPlanSummary(ctx, &booking.GetSnapshotPlanSummaryParams{
		SnapshotID:    p.SnapshotID,
		RateQualifier: p.RateQualifier,
		SupplierCode:  p.SupplierCode,
		PlanID:        p.PlanID,
		IncludeERP:    p.IncludeERP,
	})
	if err != nil {
		rlog.Error("failed to get snapshot plan summary", "snapshotID", p.SnapshotID, "error", err)
		return nil, err
	}

	selectedAddOnsJSON, err := json.Marshal(p.SelectedAddOns)
	if err != nil {
		rlog.Error("failed to marshal selected add-ons", "error", err)
		return nil, api_errors.ErrInternalError
	}

	pendingID, err := s.query.CreatePendingPayment(ctx, db.CreatePendingPaymentParams{
		UserID:          userID,
		Phone:           p.Phone,
		UserFirstName:   p.FirstName,
		UserLastName:    p.LastName,
		UserEmail:       p.Email,
		SnapshotID:      p.SnapshotID,
		RateQualifier:   p.RateQualifier,
		SupplierCode:    p.SupplierCode,
		PlanID:          p.PlanID,
		IncludeErp:      p.IncludeERP,
		SelectedAddons:  selectedAddOnsJSON,
		DriverTitle:     p.DriverTitle,
		DriverFirstName: p.DriverFirstName,
		DriverLastName:  p.DriverLastName,
		FlightNumber:    p.FlightNumber,
	})
	if err != nil {
		rlog.Error("failed to create pending customer payment", "error", err)
		return nil, api_errors.ErrInternalError
	}

	baseURL := encore.Meta().APIBaseURL.String()
	langCode := lang.FromContext(ctx, "he")

	ic := icount.NewIcount()
	resp, err := ic.GenerateIframe(icount.GenerateIframeParams{
		ClientName:   fmt.Sprintf("%s %s", p.FirstName, p.LastName),
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Email:        p.Email,
		Phone:        p.Phone,
		PaypageID:    cfg.Icount.CustomersPaypageID(),
		Sum:          planSummary.TotalPrice,
		CurrencyCode: planSummary.CurrencyCode,
		Description:  buildOrderDescription(planSummary.CarModel, planSummary.PickupDate, planSummary.PickupTime, planSummary.DropoffDate, planSummary.DropoffTime, planSummary.PickupLocationName, planSummary.DropoffLocationName, langCode),
		OrderID:      pendingID,
		PageLang:     langCode,
		SuccessURL:   fmt.Sprintf("%s%s", baseURL, SuccessPaymentGatewayPath),
		IpnURL:       fmt.Sprintf("%s%s?pending_id=%d", baseURL, CustomerPaymentIPNGatewayPath, pendingID),
	})
	if err != nil {
		rlog.Error("failed to generate paypage iframe for customer", "error", err, "pendingPaymentID", pendingID)
		return nil, api_errors.ErrInternalError
	}

	return &GenerateCustomerPaymentIframeResponse{
		URL:              resp.SaleURL,
		PendingPaymentID: pendingID,
	}, nil
}
