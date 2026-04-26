package billing

import (
	"context"

	"encore.app/services/reservation"
	"encore.dev/cron"
	"encore.dev/rlog"
)

// encore:api private
func GenerateMonthlyReport(ctx context.Context) error {
	openReservations, err := reservation.GetOpenReservations(ctx)
	if err != nil {
		rlog.Error("failed to get open reservations for monthly report", "error", err)
		return err
	}

	_ = openReservations // TODO: Implement report generation logic using the open reservations data.
	return nil
}

var _ = cron.NewJob("monthly-billing", cron.JobConfig{
	Title:    "Send Monthly Billing",
	Schedule: "0 8 1 * *", // At 08:00 on day-of-month 1.
	Endpoint: GenerateMonthlyReport,
})
