package billing

import (
	"context"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/storage"
	"encore.dev/rlog"
	"encore.dev/storage/objects"
)

// monthlyReportURLTTL is how long a download link stays valid — long enough to click, short enough
// that a link that leaks is worthless.
const monthlyReportURLTTL = 5 * time.Minute

type GetMonthlyReportURLResponse struct {
	URL string `json:"url"`
}

// GetMonthlyReportURL returns a short-lived link to one archived monthly report. The key is rebuilt
// from the validated path params rather than taken from the caller, since the bucket also holds
// other private files.
//
//encore:api auth method=GET path=/monthly-reports/:entityType/:entityID/:period/download-url tag:admin
func (s *Service) GetMonthlyReportURL(ctx context.Context, entityType string, entityID int64, period string) (*GetMonthlyReportURLResponse, error) {
	entity, err := newBillingEntity(entityType, entityID)
	if err != nil {
		return nil, api_errors.NewValidationError(err.Error())
	}

	period, err = parseMonthlyReportPeriod(period)
	if err != nil {
		return nil, api_errors.NewValidationError(err.Error())
	}

	key := monthlyReportPrefix(entity) + period + monthlyReportExtension

	exists, err := storage.PrivateFiles.Exists(ctx, key)
	if err != nil {
		rlog.Error("failed to check monthly report existence", "error", err, "key", key)
		return nil, api_errors.ErrInternalError
	}
	if !exists {
		return nil, api_errors.ErrNotFound
	}

	signed, err := storage.PrivateFiles.SignedDownloadURL(ctx, key, objects.WithTTL(monthlyReportURLTTL))
	if err != nil {
		rlog.Error("failed to sign monthly report download url", "error", err, "key", key)
		return nil, api_errors.ErrInternalError
	}

	return &GetMonthlyReportURLResponse{URL: signed.URL}, nil
}
