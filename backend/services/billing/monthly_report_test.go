package billing

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"sort"
	"strconv"
	"sync"
	"testing"

	"encore.app/services/accounts"
	"encore.app/services/notifications"
	"encore.app/services/reservation"
	"encore.dev/et"
	"github.com/xuri/excelize/v2"
)

func TestGenerateMonthlyReport(t *testing.T) {
	ctx := context.Background()

	t.Run("open reservations error", func(t *testing.T) {
		et.MockEndpoint(reservation.GetOpenReservations, func(_ context.Context) (*reservation.GetOpenReservationsResponse, error) {
			return nil, errors.New("reservations down")
		})
		et.MockEndpoint(accounts.GetBillingContacts, func(_ context.Context, _ *accounts.GetBillingContactsRequest) (*accounts.GetBillingContactsResponse, error) {
			t.Fatal("GetBillingContacts should not be called")
			return nil, nil
		})
		et.MockEndpoint(notifications.SendMonthlyReport, func(_ context.Context, _ notifications.SendMonthlyReportRequest) error {
			t.Fatal("SendMonthlyReport should not be called")
			return nil
		})

		if err := GenerateMonthlyReport(ctx); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("billing contacts error", func(t *testing.T) {
		et.MockEndpoint(reservation.GetOpenReservations, func(_ context.Context) (*reservation.GetOpenReservationsResponse, error) {
			return &reservation.GetOpenReservationsResponse{
				Reservations: []reservation.OpenReservation{{ID: 1, AgentID: 1, CurrencyCode: "USD", TotalPrice: 10}},
			}, nil
		})
		et.MockEndpoint(accounts.GetBillingContacts, func(_ context.Context, _ *accounts.GetBillingContactsRequest) (*accounts.GetBillingContactsResponse, error) {
			return nil, errors.New("contacts down")
		})
		et.MockEndpoint(notifications.SendMonthlyReport, func(_ context.Context, _ notifications.SendMonthlyReportRequest) error {
			t.Fatal("SendMonthlyReport should not be called")
			return nil
		})

		if err := GenerateMonthlyReport(ctx); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("passes distinct agent IDs to billing contacts", func(t *testing.T) {
		et.MockEndpoint(reservation.GetOpenReservations, func(_ context.Context) (*reservation.GetOpenReservationsResponse, error) {
			return &reservation.GetOpenReservationsResponse{
				Reservations: []reservation.OpenReservation{
					{ID: 1, AgentID: 10, CurrencyCode: "USD"},
					{ID: 2, AgentID: 10, CurrencyCode: "USD"}, // duplicate agent
					{ID: 3, AgentID: 20, CurrencyCode: "EUR"},
					{ID: 4, AgentID: 30, CurrencyCode: "ILS"},
				},
			}, nil
		})

		var gotIDs []int32
		et.MockEndpoint(accounts.GetBillingContacts, func(_ context.Context, p *accounts.GetBillingContactsRequest) (*accounts.GetBillingContactsResponse, error) {
			gotIDs = append([]int32(nil), p.AgentsIDs...)
			return &accounts.GetBillingContactsResponse{Contacts: nil}, nil
		})
		et.MockEndpoint(notifications.SendMonthlyReport, func(_ context.Context, _ notifications.SendMonthlyReportRequest) error {
			t.Fatal("SendMonthlyReport should not be called when there are no contacts")
			return nil
		})

		if err := GenerateMonthlyReport(ctx); err != nil {
			t.Fatalf("GenerateMonthlyReport: %v", err)
		}

		sort.Slice(gotIDs, func(i, j int) bool { return gotIDs[i] < gotIDs[j] })
		want := []int32{10, 20, 30}
		if len(gotIDs) != len(want) {
			t.Fatalf("agentsIDs = %v, want %v", gotIDs, want)
		}
		for i := range want {
			if gotIDs[i] != want[i] {
				t.Errorf("agentsIDs[%d] = %d, want %d", i, gotIDs[i], want[i])
			}
		}
	})

	t.Run("successful report sending", func(t *testing.T) {
		reservations := []reservation.OpenReservation{
			// OrgX / Office A
			{ID: 1, AgentID: 1, CurrencyCode: "USD", CarPrice: 80, ERPPrice: 20, TotalPrice: 100, DriverName: "Driver1"},
			{ID: 2, AgentID: 2, CurrencyCode: "USD", CarPrice: 40, ERPPrice: 10, TotalPrice: 50, DriverName: "Driver2"},
			// OrgX / Office B
			{ID: 3, AgentID: 3, CurrencyCode: "EUR", CarPrice: 150, ERPPrice: 50, TotalPrice: 200, DriverName: "Driver3"},
			// OrgY / Office C
			{ID: 4, AgentID: 4, CurrencyCode: "ILS", CarPrice: 250, ERPPrice: 50, TotalPrice: 300, DriverName: "Driver4"},
			// OrgY / Office D
			{ID: 5, AgentID: 5, CurrencyCode: "USD", CarPrice: 120, ERPPrice: 30, TotalPrice: 150, DriverName: "Driver5"},
			// Agent not assigned to any contact - must be skipped.
			{ID: 6, AgentID: 99, CurrencyCode: "USD", TotalPrice: 999, DriverName: "Ignored"},
		}
		et.MockEndpoint(reservation.GetOpenReservations, func(_ context.Context) (*reservation.GetOpenReservationsResponse, error) {
			return &reservation.GetOpenReservationsResponse{Reservations: reservations}, nil
		})

		contacts := []accounts.BillingContact{
			{
				ContactName: "Alice", ContactEmail: "alice@example.com",
				OrganizationID: 1, OrganizationName: "OrgX", IsOrganic: true,
				Offices: []accounts.Office{
					{ID: 1, Name: "Office A", Agents: []accounts.Agent{{ID: 1, Name: "Agent1"}, {ID: 2, Name: "Agent2"}}},
					{ID: 2, Name: "Office B", Agents: []accounts.Agent{{ID: 3, Name: "Agent3"}}},
				},
			},
			{
				ContactName: "Bob", ContactEmail: "bob@example.com",
				OrganizationID: 2, OrganizationName: "OrgY", IsOrganic: true,
				Offices: []accounts.Office{
					{ID: 3, Name: "Office C", Agents: []accounts.Agent{{ID: 4, Name: "Agent4"}}},
					{ID: 4, Name: "Office D", Agents: []accounts.Agent{{ID: 5, Name: "Agent5"}}},
				},
			},
		}
		et.MockEndpoint(accounts.GetBillingContacts, func(_ context.Context, _ *accounts.GetBillingContactsRequest) (*accounts.GetBillingContactsResponse, error) {
			return &accounts.GetBillingContactsResponse{Contacts: contacts}, nil
		})

		var mu sync.Mutex
		sent := map[string]notifications.SendMonthlyReportRequest{}
		et.MockEndpoint(notifications.SendMonthlyReport, func(_ context.Context, p notifications.SendMonthlyReportRequest) error {
			mu.Lock()
			defer mu.Unlock()
			sent[p.ContactEmail] = p
			return nil
		})

		if err := GenerateMonthlyReport(ctx); err != nil {
			t.Fatalf("GenerateMonthlyReport: %v", err)
		}

		if len(sent) != 2 {
			t.Fatalf("expected 2 sent reports, got %d", len(sent))
		}

		expected := map[string]struct {
			contactName string
			sheetName   string
			totals      map[string]float64
		}{
			"alice@example.com": {
				contactName: "Alice",
				sheetName:   "OrgX",
				totals:      map[string]float64{"USD": 150, "EUR": 200},
			},
			"bob@example.com": {
				contactName: "Bob",
				sheetName:   "OrgY",
				totals:      map[string]float64{"ILS": 300, "USD": 150},
			},
		}

		for email, exp := range expected {
			got, ok := sent[email]
			if !ok {
				t.Errorf("no report sent to %s", email)
				continue
			}
			if got.ContactName != exp.contactName {
				t.Errorf("%s: contact name = %q, want %q", email, got.ContactName, exp.contactName)
			}

			raw, err := base64.StdEncoding.DecodeString(got.ExcelBase64)
			if err != nil {
				t.Errorf("%s: base64 decode: %v", email, err)
				continue
			}
			f, err := excelize.OpenReader(bytes.NewReader(raw))
			if err != nil {
				t.Errorf("%s: open xlsx: %v", email, err)
				continue
			}

			sheets := f.GetSheetList()
			if len(sheets) != 1 || sheets[0] != exp.sheetName {
				t.Errorf("%s: sheets = %v, want [%q]", email, sheets, exp.sheetName)
				_ = f.Close()
				continue
			}

			rows, err := f.GetRows(exp.sheetName)
			if err != nil {
				t.Errorf("%s: GetRows: %v", email, err)
				_ = f.Close()
				continue
			}
			if len(rows) == 0 || len(rows[0]) != monthlyReportColCount {
				t.Errorf("%s: header row missing or wrong width: %v", email, rows)
				_ = f.Close()
				continue
			}

			// Total rows only populate the last two columns (currency + amount).
			totals := map[string]float64{}
			for _, row := range rows[1:] {
				if len(row) < monthlyReportColCount {
					continue
				}
				if row[colOfficeName] != "" {
					continue
				}
				currency := row[monthlyReportColCount-2]
				amountStr := row[monthlyReportColCount-1]
				if currency == "" || amountStr == "" {
					continue
				}
				amount, parseErr := strconv.ParseFloat(amountStr, 64)
				if parseErr != nil {
					t.Errorf("%s: parse total %q: %v", email, amountStr, parseErr)
					continue
				}
				totals[currency] += amount
			}

			if len(totals) != len(exp.totals) {
				t.Errorf("%s: totals = %v, want %v", email, totals, exp.totals)
			}
			for c, v := range exp.totals {
				if totals[c] != v {
					t.Errorf("%s: total[%s] = %v, want %v", email, c, totals[c], v)
				}
			}

			_ = f.Close()
		}
	})
}
