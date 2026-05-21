package reservation

import (
	"context"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportParams struct {
	Page            int32  `query:"page" validate:"required,gte=1"`
	PickupDateFrom  string `query:"pickupDateFrom,omitempty" encore:"optional"`
	PickupDateTo    string `query:"pickupDateTo,omitempty" encore:"optional"`
	CreatedDateFrom string `query:"createdDateFrom,omitempty" encore:"optional"`
	CreatedDateTo   string `query:"createdDateTo,omitempty" encore:"optional"`
	VoucheredAtFrom string `query:"voucheredAtFrom,omitempty" encore:"optional"`
	VoucheredAtTo   string `query:"voucheredAtTo,omitempty" encore:"optional"`
	Status          string `query:"status,omitempty" encore:"optional"`
	Broker          string `query:"broker,omitempty" encore:"optional"`
	Supplier        string `query:"supplier,omitempty" encore:"optional"`
	OrganizationID  int64  `query:"organizationId,omitempty" encore:"optional"`
	OfficeID        int64  `query:"officeId,omitempty" encore:"optional"`
	AgentID         int64  `query:"agentId,omitempty" encore:"optional"`
	IsBusiness      bool   `query:"isBusiness,omitempty" encore:"optional"`
}

const reportPageSize int64 = 50

func nilIfZero(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func nullBrokerFromString(s string) db.NullBroker {
	if s == "" {
		return db.NullBroker{}
	}
	return db.NullBroker{Broker: db.Broker(s), Valid: true}
}

func timestamptzFromString(s string) pgtype.Timestamptz {
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
	}
	return dbadapters.DBTime(t)
}

func (s *Service) getReports(ctx context.Context, p ReportParams, isBusiness bool) ([]db.Reservation, int64, error) {
	offset := int64(p.Page-1) * reportPageSize

	queryParams := db.ListReservationsReportParams{
		PickupDateFrom:  dbadapters.DateFromString(p.PickupDateFrom),
		PickupDateTo:    dbadapters.DateFromString(p.PickupDateTo),
		CreatedDateFrom: timestamptzFromString(p.CreatedDateFrom),
		CreatedDateTo:   timestamptzFromString(p.CreatedDateTo),
		VoucheredAtFrom: timestamptzFromString(p.VoucheredAtFrom),
		VoucheredAtTo:   timestamptzFromString(p.VoucheredAtTo),
		Status:          nullStatusFromString(p.Status),
		Broker:          nullBrokerFromString(p.Broker),
		SupplierCode:    nilIfEmpty(p.Supplier),
		OrganizationID:  nilIfZero(p.OrganizationID),
		OfficeID:        nilIfZero(p.OfficeID),
		AgentID:         nilIfZero(p.AgentID),
		IsBusiness:      isBusiness,
		PageSize:        reportPageSize,
		PageOffset:      offset,
	}

	reservations, err := s.query.ListReservationsReport(ctx, queryParams)
	if err != nil {
		rlog.Error("failed to list reservations report", "error", err)
		return nil, 0, api_errors.ErrInternalError
	}

	countParams := db.CountReservationsReportParams{
		PickupDateFrom:  queryParams.PickupDateFrom,
		PickupDateTo:    queryParams.PickupDateTo,
		CreatedDateFrom: queryParams.CreatedDateFrom,
		CreatedDateTo:   queryParams.CreatedDateTo,
		VoucheredAtFrom: queryParams.VoucheredAtFrom,
		VoucheredAtTo:   queryParams.VoucheredAtTo,
		Status:          queryParams.Status,
		Broker:          queryParams.Broker,
		SupplierCode:    queryParams.SupplierCode,
		OrganizationID:  queryParams.OrganizationID,
		OfficeID:        queryParams.OfficeID,
		AgentID:         queryParams.AgentID,
		IsBusiness:      isBusiness,
	}

	total, err := s.query.CountReservationsReport(ctx, countParams)
	if err != nil {
		rlog.Error("failed to count reservations report", "error", err)
		return nil, 0, api_errors.ErrInternalError
	}

	return reservations, total, nil
}
