package reports

import (
	"context"
	"strings"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportParams struct {
	Page                int32  `query:"page" validate:"required,gte=1"`
	PageSize            int64  `query:"pageSize" validate:"required,gte=1"`
	BrokerReservationID string `query:"brokerReservationId,omitempty" encore:"optional"`
	PickupDateFrom      string `query:"pickupDateFrom,omitempty" encore:"optional"`
	PickupDateTo        string `query:"pickupDateTo,omitempty" encore:"optional"`
	CreatedDateFrom     string `query:"createdDateFrom,omitempty" encore:"optional"`
	CreatedDateTo       string `query:"createdDateTo,omitempty" encore:"optional"`
	VoucheredAtFrom     string `query:"voucheredAtFrom,omitempty" encore:"optional"`
	VoucheredAtTo       string `query:"voucheredAtTo,omitempty" encore:"optional"`
	Status              string `query:"status,omitempty" encore:"optional"`
	Broker              string `query:"broker,omitempty" encore:"optional"`
	Supplier            string `query:"supplier,omitempty" encore:"optional"`
	OrganizationID      int64  `query:"organizationId,omitempty" encore:"optional"`
	OfficeID            int64  `query:"officeId,omitempty" encore:"optional"`
	AgentID             int64  `query:"agentId,omitempty" encore:"optional"`
	UserType            string `query:"userType,required" validate:"required,oneof=agent customer both"`
	IsExport            bool   `query:"isExport,omitempty" encore:"optional"`
	SkipCanceled        bool   `query:"skipCanceled,omitempty" encore:"optional"`
}

func nilIfZero(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func timestamptzFromString(s string, endOfDay bool) pgtype.Timestamptz {
	if s == "" {
		return pgtype.Timestamptz{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// try date-only
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return pgtype.Timestamptz{}
		}
		if endOfDay {
			t = t.Add(24*time.Hour - time.Nanosecond)
		}
	}
	return dbadapters.DBTime(t)
}

type getReportResult struct {
	Reservations       []db.Reservation
	Count              int64
	TotalSales         float64
	TotalCarCost       float64
	TotalBrokerErpCost float64
}

func (s *ReportsService) getReports(ctx context.Context, p ReportParams) (*getReportResult, error) {
	offset := int64(p.Page-1) * p.PageSize
	if p.IsExport {
		offset = 0
		p.PageSize = 1000000
	}

	var isBusiness *bool
	if p.UserType == "agent" {
		ib := true
		isBusiness = &ib
	} else if p.UserType == "customer" {
		ib := false
		isBusiness = &ib
	}

	queryParams := db.ListReservationsReportParams{
		BrokerReservationID: &p.BrokerReservationID,
		PickupDateFrom:      dbadapters.DateFromString(p.PickupDateFrom),
		PickupDateTo:        dbadapters.DateFromString(p.PickupDateTo),
		CreatedDateFrom:     timestamptzFromString(p.CreatedDateFrom, false),
		CreatedDateTo:       timestamptzFromString(p.CreatedDateTo, true),
		VoucheredAtFrom:     timestamptzFromString(p.VoucheredAtFrom, false),
		VoucheredAtTo:       timestamptzFromString(p.VoucheredAtTo, true),
		Status:              dbadapters.StringToNullStatus[db.ReservationStatus](p.Status),
		Broker:              dbadapters.StringToNullStatus[db.Broker](p.Broker),
		SupplierCodes:       splitSupplierCodes(p.Supplier),
		OrganizationID:      nilIfZero(p.OrganizationID),
		OfficeID:            nilIfZero(p.OfficeID),
		AgentID:             nilIfZero(p.AgentID),
		IsBusiness:          isBusiness,
		PageSize:            p.PageSize,
		PageOffset:          offset,
		SkipCanceled:        p.SkipCanceled,
	}

	reservations, err := s.query.ListReservationsReport(ctx, queryParams)
	if err != nil {
		rlog.Error("failed to list reservations report", "error", err)
		return nil, api_errors.ErrInternalError
	}

	countParams := db.CountReservationsReportParams{
		BrokerReservationID: queryParams.BrokerReservationID,
		PickupDateFrom:      queryParams.PickupDateFrom,
		PickupDateTo:        queryParams.PickupDateTo,
		CreatedDateFrom:     queryParams.CreatedDateFrom,
		CreatedDateTo:       queryParams.CreatedDateTo,
		VoucheredAtFrom:     queryParams.VoucheredAtFrom,
		VoucheredAtTo:       queryParams.VoucheredAtTo,
		Status:              queryParams.Status,
		Broker:              queryParams.Broker,
		SupplierCodes:       queryParams.SupplierCodes,
		OrganizationID:      queryParams.OrganizationID,
		OfficeID:            queryParams.OfficeID,
		AgentID:             queryParams.AgentID,
		IsBusiness:          isBusiness,
		SkipCanceled:        p.SkipCanceled,
	}

	total, err := s.query.CountReservationsReport(ctx, countParams)
	if err != nil {
		rlog.Error("failed to count reservations report", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &getReportResult{
		Reservations:       reservations,
		Count:              total.Count,
		TotalSales:         total.TotalSales,
		TotalCarCost:       total.TotalCarCost,
		TotalBrokerErpCost: total.TotalBrokerErpCost,
	}, nil
}

func splitSupplierCodes(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
