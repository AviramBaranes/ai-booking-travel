package billing

import (
	"context"
	"encoding/base64"
	"math"
	"time"

	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.app/services/notifications"
	"encore.app/services/reservation"
	"encore.dev/cron"
	"encore.dev/rlog"
)

type Report struct {
	OrganizationName  string
	IsOrganic         bool
	ContactName       string
	ContactEmail      string
	TransactionGroups []TransactionGroup
}

type TransactionGroup struct {
	TotalAmount  float64
	Currency     string
	Reservations []Reservations
}

type Reservations struct {
	ReservationID           int64
	OfficeName              string
	AgentName               string
	DriverName              string
	ReservationCreationDate string
	ReservationBrokerID     string
	VoucherDate             string
	VoucherNumber           string
	PickupDate              string
	DropoffDate             string
	RentalDays              int
	CountryCode             string
	Currency                string
	CarPrice                float64
	ERPPrice                float64
	TotalPrice              float64
	// PenaltyType is empty for a rental row, and holds the fee type for a cancellation or
	// no-show row. It drives both the label on the reservation id and the blank price columns.
	PenaltyType string
}

type AgentInfo struct {
	AgentID    int64
	AgentName  string
	OfficeName string
}

// encore:api private
func GenerateMonthlyReport(ctx context.Context) error {
	openReservations, err := reservation.GetOpenReservations(ctx)
	if err != nil {
		rlog.Error("failed to get open reservations for monthly report", "error", err)
		return err
	}

	agentsSet := make(map[int64]struct{})
	for _, r := range openReservations.Reservations {
		agentsSet[r.AgentID] = struct{}{}
	}
	for _, p := range openReservations.Penalties {
		agentsSet[p.AgentID] = struct{}{}
	}

	agentsIDs := make([]int64, 0, len(agentsSet))
	for id := range agentsSet {
		agentsIDs = append(agentsIDs, id)
	}

	billingContacts, err := accounts.GetBillingContacts(ctx, contact.GetBillingContactsParams{
		AgentsIDs: agentsIDs,
	})
	if err != nil {
		rlog.Error("failed to get billing contacts for monthly report", "error", err)
		return err
	}

	reports := generateReports(openReservations, billingContacts)
	sendReports(ctx, reports)

	return nil
}

func generateReports(openReservations *reservation.GetOpenReservationsResponse, billingContacts *contact.GetBillingContactsResponse) []Report {
	reports := make([]Report, 0, len(billingContacts.Contacts))
	for _, c := range billingContacts.Contacts {
		report := Report{
			OrganizationName:  c.OrganizationName,
			IsOrganic:         c.IsOrganic,
			ContactName:       c.ContactName,
			ContactEmail:      c.ContactEmail,
			TransactionGroups: generateTransactionGroups(openReservations, c),
		}

		reports = append(reports, report)
	}

	return reports
}

func generateTransactionGroups(openReservations *reservation.GetOpenReservationsResponse, contact contact.BillingContact) []TransactionGroup {
	relevantAgents := make(map[int64]AgentInfo)
	for _, office := range contact.Offices {
		for _, agent := range office.Agents {
			relevantAgents[agent.ID] = AgentInfo{
				AgentID:    agent.ID,
				AgentName:  agent.Name,
				OfficeName: office.Name,
			}
		}
	}

	tgs := make([]TransactionGroup, 0)

	for _, r := range openReservations.Reservations {
		agentInfo, isAgentRelevant := relevantAgents[r.AgentID]
		if !isAgentRelevant {
			continue
		}

		tgs = addToTransactionGroups(tgs, r.CurrencyCode, toReportReservation(r, agentInfo))
	}

	for _, p := range openReservations.Penalties {
		agentInfo, isAgentRelevant := relevantAgents[p.AgentID]
		if !isAgentRelevant {
			continue
		}

		tgs = addToTransactionGroups(tgs, p.CurrencyCode, toReportPenalty(p, agentInfo))
	}

	return tgs
}

// addToTransactionGroups files a report row under its currency group, creating the group on first
// use, and adds its price to the group total.
func addToTransactionGroups(tgs []TransactionGroup, currencyCode string, row Reservations) []TransactionGroup {
	tgIndex := -1
	for i, tg := range tgs {
		if tg.Currency == currencyCode {
			tgIndex = i
			break
		}
	}

	if tgIndex == -1 {
		return append(tgs, TransactionGroup{
			Currency:     currencyCode,
			Reservations: []Reservations{row},
			TotalAmount:  row.TotalPrice,
		})
	}

	tgs[tgIndex].Reservations = append(tgs[tgIndex].Reservations, row)
	tgs[tgIndex].TotalAmount += row.TotalPrice

	return tgs
}

func toReportReservation(r reservation.OpenReservation, agentInfo AgentInfo) Reservations {
	m := 1.0
	if r.PaymentStatus == "refund_pending" {
		m = -1
	}

	return Reservations{
		ReservationID:           r.ID,
		OfficeName:              agentInfo.OfficeName,
		AgentName:               agentInfo.AgentName,
		DriverName:              r.DriverName,
		ReservationCreationDate: formatDate(r.CreatedAt),
		ReservationBrokerID:     r.BrokerReservationID,
		VoucherDate:             formatDate(r.VoucheredAt),
		VoucherNumber:           r.VoucherNumber,
		PickupDate:              formatDate(r.PickupDate),
		DropoffDate:             formatDate(r.DropoffDate),
		RentalDays:              r.RentalDays,
		CountryCode:             r.CountryCode,
		Currency:                r.CurrencyCode,
		CarPrice:                roundPrice(r.CarPrice * m),
		ERPPrice:                roundPrice(r.ERPPrice * m),
		TotalPrice:              roundPrice(r.TotalPrice * m),
	}
}

// toReportPenalty maps a fee to a report row. Every column other than the amount and the fee label
// repeats the reservation the fee was charged on, so the row reads like the booking it belongs to.
// The car and coverage prices are left unset: a fee is a single amount, not a priced rental.
func toReportPenalty(p reservation.OpenPenalty, agentInfo AgentInfo) Reservations {
	return Reservations{
		ReservationID:           p.ReservationID,
		OfficeName:              agentInfo.OfficeName,
		AgentName:               agentInfo.AgentName,
		DriverName:              p.DriverName,
		ReservationCreationDate: formatDate(p.CreatedAt),
		ReservationBrokerID:     p.BrokerReservationID,
		VoucherDate:             formatDate(p.VoucheredAt),
		VoucherNumber:           p.VoucherNumber,
		PickupDate:              formatDate(p.PickupDate),
		DropoffDate:             formatDate(p.DropoffDate),
		RentalDays:              p.RentalDays,
		CountryCode:             p.CountryCode,
		Currency:                p.CurrencyCode,
		TotalPrice:              roundPrice(p.Amount),
		PenaltyType:             p.Type,
	}
}

func sendReports(ctx context.Context, reports []Report) {
	for _, r := range reports {
		excelReport, err := generateExcelReport(r)
		if err != nil {
			rlog.Error("failed to generate excel report", "error", err, "contact_email", r.ContactEmail)
			continue
		}

		base64ExcelReport := base64.StdEncoding.EncodeToString(excelReport)

		if err = notifications.SendMonthlyReport(ctx, notifications.SendMonthlyReportParams{
			ContactName:  r.ContactName,
			ContactEmail: r.ContactEmail,
			ExcelBase64:  base64ExcelReport,
		}); err != nil {
			rlog.Error("failed to send monthly report", "error", err, "contact_email", r.ContactEmail)
		}
	}
}

var _ = cron.NewJob("monthly-billing", cron.JobConfig{
	Title:    "Send Monthly Billing",
	Schedule: "0 7 1 * *", // At 08:00 on day-of-month 1.
	Endpoint: GenerateMonthlyReport,
})

func formatDate(date string) string {
	if date == "" {
		return ""
	}

	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t.Format("02/01/2006")
		}
	}

	return date
}

func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}
