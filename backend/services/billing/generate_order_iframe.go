package billing

import (
	"context"
	"fmt"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/icount"
	"encore.app/internal/lang"
	"encore.app/services/accounts"
	"encore.app/services/accounts/handlers/lookup"
	"encore.app/services/reservation"
	"encore.dev"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"golang.org/x/sync/errgroup"
)

type GenerateOrderIframeParams struct {
	OrderID int64 `json:"orderId"`
	IsILS   bool  `json:"isIls"`
}

type GenerateOrderIframeResponse struct {
	URL string `json:"url"`
}

var (
	ErrOrderNotUnpaid = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "invalid payment status on reservation, only unpaid reservations can be paid", api_errors.ErrorDetails{
		Code: api_errors.CodeOrderNotUnpaid,
	})

	ErrOrderNotBooked = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "invalid reservation status, only booked reservations can be paid", api_errors.ErrorDetails{
		Code: api_errors.CodeOrderNotBooked,
	})
)

// GenerateOrderIframe generates an iframe for the order page. It returns the URL of the iframe and an error if any occurs.
// It should be used for agents that can't voucher an order because of a finished obligo.
//
// encore:api auth path=/billing/generate-order-iframe method=POST tag:agent
func (s *Service) GenerateOrderIframe(ctx context.Context, p GenerateOrderIframeParams) (*GenerateOrderIframeResponse, error) {
	authData := auth.Data().(*accounts.AuthData)

	g, groupCtx := errgroup.WithContext(ctx)

	var clientName string
	var order *reservation.GetReservationResponse
	// var rates map[string]float64

	g.Go(func() error {
		cn, err := getClientName(groupCtx, authData.UserID)
		if err != nil {
			return err
		}
		clientName = cn
		return nil
	})

	g.Go(func() error {
		o, err := getOrder(groupCtx, p.OrderID)
		if err != nil {
			return err
		}
		order = o
		return nil
	})

	if err := g.Wait(); err != nil {
		rlog.Error("failed to get data for generating order iframe", "error", err)
		return nil, err
	}

	sum, currency, err := s.getPrice(ctx, order, p.IsILS)
	if err != nil {
		rlog.Error("failed to get price", "error", err)
		return nil, err
	}

	baseURL := encore.Meta().APIBaseURL.String()
	// baseURL = "https://2d9b-31-154-63-122.ngrok-free.app"

	langCode := lang.FromContext(ctx, "he")
	ic := icount.NewIcount()
	resp, err := ic.GenerateIframe(icount.GenerateIframeParams{
		ClientName:   clientName,
		PaypageID:    cfg.Icount.AgentsPaypageID(),
		Sum:          sum,
		CurrencyCode: currency,
		Description:  buildOrderDescription(order.CarDetails.Model, order.PickupDate, order.PickupTime, order.DropoffDate, order.DropoffTime, order.PickupLocationName, order.DropoffLocationName, langCode),
		OrderID:      p.OrderID,
		PageLang:     langCode,
		SuccessURL:   fmt.Sprintf("%s%s", baseURL, SuccessPaymentGatewayPath),
		IpnURL:       getIpnURL(baseURL, p.OrderID, authData),
	})

	if err != nil {
		rlog.Error("failed to generate paypage iframe", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GenerateOrderIframeResponse{
		URL: resp.SaleURL,
	}, nil
}

// getClientName uses the accounts service to get the name of the user with the given auth data. It returns the name of the user and an error if any occurs.
func getClientName(ctx context.Context, userID int64) (string, error) {
	resp, err := accounts.GetAccountsLookup(ctx, lookup.GetAccountsLookupParams{
		UserIDs: []int64{userID},
	})
	if err != nil {
		rlog.Error("failed to get user name", "error", err)
		return "", api_errors.ErrInternalError
	}

	if len(resp.Users) == 0 {
		rlog.Error("no user found for user id", "user_id", userID)
		return "", api_errors.ErrInternalError
	}

	return resp.Users[0].Name, nil
}

// getOrder retrieves the order with the given order ID from the reservation service.
func getOrder(ctx context.Context, orderID int64) (*reservation.GetReservationResponse, error) {
	resp, err := reservation.GetReservation(ctx, orderID)
	if err != nil {
		rlog.Error("failed to get order", "error", err, "order_id", orderID)
		return nil, api_errors.ErrInternalError
	}

	if resp.PaymentStatus != reservation.PaymentStatusUnpaid {
		rlog.Error("order already paid", "order_id", orderID)
		return nil, ErrOrderNotUnpaid
	}

	if resp.ReservationStatus != reservation.ReservationStatusBooked {
		rlog.Error("order not booked", "order_id", orderID, "reservation_status", resp.ReservationStatus)
		return nil, ErrOrderNotBooked
	}

	return resp, nil
}

// getPrice calculates the price of the order. If isILS is false, it returns the total price of the order and its currency code. If isILS is true, it converts the total price to ILS using the exchange rate from iCount and returns the converted price and "ILS" as the currency code.
func (s *Service) getPrice(ctx context.Context, order *reservation.GetReservationResponse, isILS bool) (float64, string, error) {
	if !isILS {
		return order.TotalPriceFloat, order.CurrencyCode, nil
	}

	rate, err := s.ratesCache.GetCurrencyRate(ctx, order.CurrencyCode)
	rlog.Info("currency rates", "rate", rate)
	if err != nil {
		rlog.Error("currency code not found in rates", "currency_code", order.CurrencyCode)
		return 0, "", api_errors.ErrInternalError
	}

	return order.TotalPriceFloat * rate, "ILS", nil
}

// getIpnURL constructs the IPN URL for the order payment notification. It takes the base URL of the API, the order ID, and the auth data of the user making the request. It returns the constructed IPN URL.
func getIpnURL(baseURL string, orderID int64, authData *accounts.AuthData) string {
	billingEntityParam := "&office_id=" + strconv.FormatInt(authData.OrganizationContext.OfficeID, 10)
	if authData.OrganizationContext.IsOrganic {
		billingEntityParam = "&organization_id=" + strconv.FormatInt(authData.OrganizationContext.OrganizationID, 10)
	}
	return fmt.Sprintf("%s%s?id=%d%s", baseURL, OrderPaymentIPNGatewayPath, orderID, billingEntityParam)
}

// buildOrderDescription generates a localised car rental description for iCount payment pages.
func buildOrderDescription(model, pickupDate, pickupTime, dropoffDate, dropoffTime, pickupLocation, dropoffLocation, langCode string) string {
	sameLocation := pickupLocation == dropoffLocation

	switch langCode {
	case "he":
		if sameLocation {
			return fmt.Sprintf(
				"השכרת רכב %s או דומה.\nאיסוף והחזרה ב%s.\nאיסוף: %s בשעה %s\nהחזרה: %s בשעה %s",
				model, pickupLocation, pickupDate, pickupTime, dropoffDate, dropoffTime,
			)
		}
		return fmt.Sprintf(
			"השכרת רכב %s או דומה.\nאיסוף מ%s בתאריך %s בשעה %s.\nהחזרה ב%s בתאריך %s בשעה %s.",
			model, pickupLocation, pickupDate, pickupTime, dropoffLocation, dropoffDate, dropoffTime,
		)
	default:
		if sameLocation {
			return fmt.Sprintf(
				"Car rental for %s or similar.\nPickup and return at %s.\nPickup: %s at %s\nReturn: %s at %s",
				model, pickupLocation, pickupDate, pickupTime, dropoffDate, dropoffTime,
			)
		}
		return fmt.Sprintf(
			"Car rental for %s or similar.\nPickup at %s on %s at %s.\nReturn at %s on %s at %s.",
			model, pickupLocation, pickupDate, pickupTime, dropoffLocation, dropoffDate, dropoffTime,
		)
	}
}
