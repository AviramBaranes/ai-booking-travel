package reports

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type BillingEntity string

const (
	BillingEntityBusiness BillingEntity = "organization"
	BillingEntityOffice   BillingEntity = "office"
)

type BusinessesBalancesReportRow struct {
	BillingEntityType        BillingEntity `json:"billingEntityType"`
	BillingEntityID          int64         `json:"billingEntityId"`
	BillingEntityName        string        `json:"billingEntityName"`
	TotalOpenBalanceInEuro   float64       `json:"totalOpenBalanceInEuro"`
	TotalOpenBalanceInDollar float64       `json:"totalOpenBalanceInDollar"`
	TotalInOtherCurrency     float64       `json:"totalInOtherCurrency"`
}

type BusinessesBalancesReportResponse struct {
	Businesses []BusinessesBalancesReportRow `json:"businesses"`
	Total      int64                         `json:"total"`
}

func (s *ReportsService) GetBusinessesBalancesReport(ctx context.Context) (*BusinessesBalancesReportResponse, error) {
	rows, err := s.query.ListBusinessesBalancesReport(ctx)
	if err != nil {
		rlog.Error("failed to list businesses balances report", "error", err)
		return nil, api_errors.ErrInternalError
	}

	billingEntitiesSet := buildBillingEntitiesSet(rows)
	accountsLookup, err := accounts.GetAccountsLookup(ctx, accounts.GetAccountsLookupParams{
		OrganizationIDs: idsFromSet(billingEntitiesSet.organizationIDs),
		OfficeIDs:       idsFromSet(billingEntitiesSet.officeIDs),
	})
	if err != nil {
		rlog.Error("failed to get accounts lookup for businesses balances report", "error", err)
		return nil, err
	}

	businesses := buildBusinessesBalancesReportRows(rows, accountsLookup)

	return &BusinessesBalancesReportResponse{
		Businesses: businesses,
		Total:      int64(len(businesses)),
	}, nil
}

func buildBillingEntitiesSet(rows []db.ListBusinessesBalancesReportRow) accountsSet {
	organizationIDs := make(map[int64]struct{})
	officeIDs := make(map[int64]struct{})

	for _, row := range rows {
		switch BillingEntity(row.BillingEntityType) {
		case BillingEntityBusiness:
			organizationIDs[row.BillingEntityID] = struct{}{}
		case BillingEntityOffice:
			officeIDs[row.BillingEntityID] = struct{}{}
		}
	}

	return accountsSet{
		organizationIDs: organizationIDs,
		officeIDs:       officeIDs,
		userIDs:         map[int64]struct{}{},
	}
}

func buildBusinessesBalancesReportRows(rows []db.ListBusinessesBalancesReportRow, accountsLookup *accounts.GetAccountsLookupResponse) []BusinessesBalancesReportRow {
	organizationNames := namesByID(accountsLookup.Organizations)
	officeNames := namesByID(accountsLookup.Offices)

	businesses := make([]BusinessesBalancesReportRow, 0, len(rows))
	for _, row := range rows {
		billingEntityType := BillingEntity(row.BillingEntityType)

		businesses = append(businesses, BusinessesBalancesReportRow{
			BillingEntityType:        billingEntityType,
			BillingEntityID:          row.BillingEntityID,
			BillingEntityName:        billingEntityName(billingEntityType, row.BillingEntityID, organizationNames, officeNames),
			TotalOpenBalanceInEuro:   row.TotalEur,
			TotalOpenBalanceInDollar: row.TotalUsd,
			TotalInOtherCurrency:     row.TotalOtherConverted,
		})
	}

	return businesses
}

func billingEntityName(entityType BillingEntity, entityID int64, organizationNames map[int64]string, officeNames map[int64]string) string {
	switch entityType {
	case BillingEntityBusiness:
		return organizationNames[entityID]
	case BillingEntityOffice:
		return officeNames[entityID]
	default:
		return ""
	}
}
