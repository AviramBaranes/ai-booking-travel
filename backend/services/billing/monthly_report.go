package billing

import (
	"context"

	"encore.app/services/accounts"
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

	agentsSet := make(map[int32]struct{})
	for _, r := range openReservations.Reservations {
		agentsSet[r.AgentID] = struct{}{}
	}

	agentsIDs := make([]int32, 0, len(agentsSet))
	for id := range agentsSet {
		agentsIDs = append(agentsIDs, id)
	}

	billingContacts, err := accounts.GetBillingContacts(ctx, &accounts.GetBillingContactsRequest{
		AgentsIDs: agentsIDs,
	})
	if err != nil {
		rlog.Error("failed to get billing contacts for monthly report", "error", err)
		return err
	}

	_ = billingContacts

	return nil
}

var _ = cron.NewJob("monthly-billing", cron.JobConfig{
	Title:    "Send Monthly Billing",
	Schedule: "0 8 1 * *", // At 08:00 on day-of-month 1.
	Endpoint: GenerateMonthlyReport,
})
