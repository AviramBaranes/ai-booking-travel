package billing

import (
	"context"

	"encore.app/services/accounts"
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
	ReturnDate              string
	RentalDays              int
	CountryCode             string
	Currency                string
	CarPrice                float64
	ERPPrice                float64
	TotalPrice              float64
}

type AgentInfo struct {
	AgentID    int32
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

	reports := generateReports(openReservations, billingContacts)
	sendReports(reports)

	return nil
}

func generateReports(openReservations *reservation.GetOpenReservationsResponse, billingContacts *accounts.GetBillingContactsResponse) []Report {
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

func generateTransactionGroups(openReservations *reservation.GetOpenReservationsResponse, contact accounts.BillingContact) []TransactionGroup {
	relevantAgents := make(map[int32]AgentInfo)
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

		reservation := toReportReservation(r, agentInfo)

		tgIndex := -1
		for i, tg := range tgs {
			if tg.Currency == r.CurrencyCode {
				tgIndex = i
				break
			}
		}

		if tgIndex == -1 {
			tgs = append(tgs, TransactionGroup{
				Currency:     r.CurrencyCode,
				Reservations: []Reservations{reservation},
				TotalAmount:  reservation.TotalPrice,
			})
		} else {
			tgs[tgIndex].Reservations = append(tgs[tgIndex].Reservations, reservation)
			tgs[tgIndex].TotalAmount += reservation.TotalPrice
		}
	}

	return nil
}

func toReportReservation(r reservation.OpenReservation, agentInfo AgentInfo) Reservations {
	var m float64 = 1
	if r.PaymentStatus == "refund_pending" {
		m = -1
	}

	return Reservations{
		ReservationID:           r.ID,
		OfficeName:              agentInfo.OfficeName,
		AgentName:               agentInfo.AgentName,
		DriverName:              r.DriverName,
		ReservationCreationDate: r.CreatedAt,
		ReservationBrokerID:     r.BrokerReservationID,
		VoucherDate:             r.VoucheredAt,
		VoucherNumber:           r.VoucherNumber,
		PickupDate:              r.PickupDate,
		ReturnDate:              r.DropoffDate,
		RentalDays:              r.RentalDays,
		CountryCode:             r.CountryCode,
		Currency:                r.CountryCode,
		CarPrice:                r.CarPrice * m,
		ERPPrice:                r.ERPPrice * m,
		TotalPrice:              r.TotalPrice * m,
	}
}

func sendReports(reports []Report) {
	for _, r := range reports {
		excelReport, err := generateExcelReport(r)
		if err != nil {
			rlog.Error("failed to generate excel report", "error", err, "contact_email", r.ContactEmail)
			continue
		}
		_ = excelReport
	}
}

var _ = cron.NewJob("monthly-billing", cron.JobConfig{
	Title:    "Send Monthly Billing",
	Schedule: "0 8 1 * *", // At 08:00 on day-of-month 1.
	Endpoint: GenerateMonthlyReport,
})
