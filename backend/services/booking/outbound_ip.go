package booking

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var errOutboundIPLookupFailed = api_errors.NewErrorWithDetail(
	errs.Internal,
	"Failed to resolve the outbound IP",
	api_errors.InternalErrorDetails,
)

// OutboundIPResponse reports the public IP our outbound traffic appears to come from.
type OutboundIPResponse struct {
	IP string `json:"ip"`
}

// TEMPORARY debug endpoint: resolves the egress IP of the running instance so it
// can be handed to suppliers for allow-listing. Remove once that is done.
//
//encore:api private method=GET path=/debug/outbound-ip
func (s *Service) OutboundIP(ctx context.Context) (*OutboundIPResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		rlog.Error("failed to build outbound ip request", "error", err)
		return nil, errOutboundIPLookupFailed
	}

	resp, err := client.Do(req)
	if err != nil {
		rlog.Error("failed to resolve outbound ip", "error", err)
		return nil, errOutboundIPLookupFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rlog.Error("failed to read outbound ip response", "error", err)
		return nil, errOutboundIPLookupFailed
	}

	ip := strings.TrimSpace(string(body))
	rlog.Info("Outbound IP", "ip", ip)

	return &OutboundIPResponse{IP: ip}, nil
}
