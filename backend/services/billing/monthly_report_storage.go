package billing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"encore.app/internal/storage"
	"encore.dev/storage/objects"
)

// Monthly reports live under a deterministic key, so a re-run overwrites the file it replaces
// instead of piling up copies:
//
//	monthly-reports/organizations/{organizationID}/{YYYY-MM}.xlsx   // organic organizations
//	monthly-reports/offices/{officeID}/{YYYY-MM}.xlsx               // offices of inorganic ones
//
// The prefix down to the billing entity id is all a listing needs: list that prefix and the object
// names are the billing periods, already sorted ascending. Ids rather than names keep the key
// stable when an organization or office is renamed.
const (
	monthlyReportsPrefix = "monthly-reports"
	organizationsSegment = "organizations"
	officesSegment       = "offices"

	monthlyReportPeriodLayout = "2006-01"
	monthlyReportExtension    = ".xlsx"
	monthlyReportContentType  = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

// billingEntity owns a folder of reports: an organic organization, or an office of an inorganic
// one, which pays for itself.
type billingEntity struct {
	Type string
	ID   int64
}

// reportBillingEntity returns the entity a report is billed to.
func reportBillingEntity(r Report) (billingEntity, error) {
	if r.IsOrganic {
		return billingEntity{Type: organizationsSegment, ID: r.OrganizationID}, nil
	}

	if r.OfficeID == nil {
		return billingEntity{}, fmt.Errorf("inorganic report of organization %d has no office to bill", r.OrganizationID)
	}

	return billingEntity{Type: officesSegment, ID: *r.OfficeID}, nil
}

// newBillingEntity validates an entity that comes from outside this file — a request, or a key read
// back from the bucket — so a caller can never point us at an arbitrary object.
func newBillingEntity(entityType string, id int64) (billingEntity, error) {
	if entityType != organizationsSegment && entityType != officesSegment {
		return billingEntity{}, fmt.Errorf("unknown billing entity type %q", entityType)
	}

	if id <= 0 {
		return billingEntity{}, fmt.Errorf("invalid billing entity id %d", id)
	}

	return billingEntity{Type: entityType, ID: id}, nil
}

// monthlyReportPrefix returns the key prefix holding every report of one billing entity.
func monthlyReportPrefix(e billingEntity) string {
	return fmt.Sprintf("%s/%s/%d/", monthlyReportsPrefix, e.Type, e.ID)
}

func monthlyReportKey(r Report, period time.Time) (string, error) {
	entity, err := reportBillingEntity(r)
	if err != nil {
		return "", err
	}

	return monthlyReportPrefix(entity) + period.Format(monthlyReportPeriodLayout) + monthlyReportExtension, nil
}

// parseMonthlyReportKey reads a stored key back into the entity that owns it and the period it
// covers. A key that doesn't match the layout is rejected rather than guessed at, so an unrelated
// private file that ends up under the prefix never surfaces as a report.
func parseMonthlyReportKey(key string) (billingEntity, string, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 4 || parts[0] != monthlyReportsPrefix {
		return billingEntity{}, "", fmt.Errorf("unexpected monthly report key %q", key)
	}

	entityType, rawEntityID, filename := parts[1], parts[2], parts[3]

	entityID, err := strconv.ParseInt(rawEntityID, 10, 64)
	if err != nil {
		return billingEntity{}, "", fmt.Errorf("unexpected billing entity id in key %q", key)
	}

	entity, err := newBillingEntity(entityType, entityID)
	if err != nil {
		return billingEntity{}, "", fmt.Errorf("unexpected billing entity in key %q: %w", key, err)
	}

	if !strings.HasSuffix(filename, monthlyReportExtension) {
		return billingEntity{}, "", fmt.Errorf("unexpected monthly report filename in key %q", key)
	}

	period, err := parseMonthlyReportPeriod(strings.TrimSuffix(filename, monthlyReportExtension))
	if err != nil {
		return billingEntity{}, "", fmt.Errorf("unexpected monthly report filename in key %q: %w", key, err)
	}

	return entity, period, nil
}

// parseMonthlyReportPeriod normalizes a YYYY-MM billing period, rejecting anything else.
func parseMonthlyReportPeriod(period string) (string, error) {
	parsed, err := time.Parse(monthlyReportPeriodLayout, period)
	if err != nil {
		return "", fmt.Errorf("invalid billing period %q", period)
	}

	return parsed.Format(monthlyReportPeriodLayout), nil
}

// billingPeriod is the month a report bills for. The cron fires on the 1st over everything still
// open, so the period it closes is the month that just ended.
func billingPeriod(now time.Time) time.Time {
	return now.UTC().AddDate(0, -1, 0)
}

func storeMonthlyReport(ctx context.Context, r Report, excelReport []byte, period time.Time) error {
	key, err := monthlyReportKey(r, period)
	if err != nil {
		return err
	}

	writer := storage.PrivateFiles.Upload(ctx, key, objects.WithUploadAttrs(objects.UploadAttrs{
		ContentType: monthlyReportContentType,
	}))

	if _, err := io.Copy(writer, bytes.NewReader(excelReport)); err != nil {
		writer.Abort(err)
		return fmt.Errorf("failed to write monthly report %s: %w", key, err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to store monthly report %s: %w", key, err)
	}

	return nil
}
