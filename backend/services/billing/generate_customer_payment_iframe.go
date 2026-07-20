package billing

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/icount"
	"encore.app/internal/lang"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/accounts/handlers/user"
	"encore.app/services/billing/db"
	"encore.app/services/booking"
	"encore.dev"
	"encore.dev/beta/auth"
	"encore.dev/rlog"
)

type GenerateCustomerPaymentIframeParams struct {
	Phone           string               `json:"phone" validate:"required,israeli_phone"`
	FirstName       string               `json:"firstName" validate:"required,notblank,alphaunicode"`
	LastName        string               `json:"lastName" validate:"required,notblank,alphaunicode"`
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
	URL                 string `json:"url"`
	PendingPaymentToken string `json:"pendingPaymentToken"`
}

// GenerateCustomerPaymentIframe creates a pending payment record and returns an iCount iframe URL for the customer to complete payment.
// Pricing and booking details are derived server-side from the snapshot.
// If the request is authenticated as a customer, the user ID is stored on the pending record.
// Requests authenticated under any other role are rejected.
// Unauthenticated requests proceed without a user ID (guest checkout).
//
// encore:api public path=/billing/generate-customer-payment-iframe method=POST
func (s *Service) GenerateCustomerPaymentIframe(ctx context.Context, p GenerateCustomerPaymentIframeParams) (*GenerateCustomerPaymentIframeResponse, error) {
	userID, err := getUserID(ctx, p)
	if err != nil {
		return nil, err
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

	planILSPrice, err := s.getILSPrice(ctx, planSummary.TotalPrice, planSummary.CurrencyCode)
	if err != nil {
		return nil, err
	}

	row, err := s.query.CreatePendingPayment(ctx, db.CreatePendingPaymentParams{
		UserID:          userID,
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
	// baseURL = "https://6f76-46-116-122-75.ngrok-free.app"
	langCode := lang.FromContext(ctx, "he")

	ic := icount.NewIcount()
	resp, err := ic.GenerateIframe(icount.GenerateIframeParams{
		ClientName:   fmt.Sprintf("%s %s", p.FirstName, p.LastName),
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Email:        p.Email,
		Phone:        p.Phone,
		PaypageID:    cfg.Icount.CustomersPaypageID(),
		Sum:          planILSPrice,
		CurrencyCode: "ILS", //customers always pay in ILS, regardless of the plan's currency
		Description:  buildOrderDescription(planSummary.CarModel, planSummary.PickupDate, planSummary.PickupTime, planSummary.DropoffDate, planSummary.DropoffTime, planSummary.PickupLocationName, planSummary.DropoffLocationName, langCode),
		OrderID:      row.ID,
		PageLang:     langCode,
		SuccessURL:   fmt.Sprintf("%s%s", baseURL, SuccessPaymentGatewayPath),
		IpnURL:       fmt.Sprintf("%s%s?id=%d", baseURL, CustomerPaymentIPNGatewayPath, row.ID),
	})
	if err != nil {
		rlog.Error("failed to generate paypage iframe for customer", "error", err, "pendingPaymentID", row.ID)
		return nil, api_errors.ErrInternalError
	}

	return &GenerateCustomerPaymentIframeResponse{
		URL:                 resp.SaleURL,
		PendingPaymentToken: dbadapters.UuidToString(row.Token),
	}, nil
}

// getILSPrice calculates the price of the order in ILS based on the plan's price and currency code.
func (s *Service) getILSPrice(ctx context.Context, planPrice float64, currencyCode string) (float64, error) {
	rate, err := s.ratesCache.GetCurrencyRate(ctx, currencyCode)
	rlog.Info("currency rates", "rate", rate)
	if err != nil {
		rlog.Error("currency code not found in rates", "currency_code", currencyCode)
		return 0, api_errors.ErrInternalError
	}

	return planPrice * rate, nil
}

// getUserID retrieves the user ID from the auth context if the request is authenticated as a customer.
// If the request is unauthenticated, it attempts to get or create a customer using the provided first name, last name, email, and phone.
func getUserID(ctx context.Context, p GenerateCustomerPaymentIframeParams) (int64, error) {
	if auth.Data() != nil {
		authData := auth.Data().(*accounts.AuthData)
		if authData.Role != accounts.UserRoleCustomer {
			return 0, api_errors.ErrUnauthorized
		}

		return authData.UserID, nil
	}

	resp, err := accounts.GetOrCreateCustomer(ctx, user.GetOrCreateCustomerParams{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Phone:     p.Phone,
	})
	if err != nil {
		rlog.Error("failed to get or create customer", "error", err, "phone", p.Phone, "email", p.Email)
		return 0, err
	}

	return resp.UserID, nil
}
