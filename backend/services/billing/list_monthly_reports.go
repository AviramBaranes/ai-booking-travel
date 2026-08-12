package billing

import (
	"context"
	"fmt"
	"sort"

	"encore.app/internal/api_errors"
	"encore.app/internal/storage"
	"encore.app/services/accounts"
	"encore.dev/rlog"
	"encore.dev/storage/objects"
)

type ListMonthlyReportsResponse struct {
	Folders []MonthlyReportFolder `json:"folders"`
}

// MonthlyReportFolder is every report archived for one billing entity, presented as a folder named
// after the entity.
type MonthlyReportFolder struct {
	EntityType string              `json:"entityType"` // "organizations" or "offices"
	EntityID   int64               `json:"entityId"`
	Name       string              `json:"name"`
	Files      []MonthlyReportFile `json:"files"`
}

type MonthlyReportFile struct {
	Period string `json:"period"` // YYYY-MM
	Size   int64  `json:"size"`
}

// ListMonthlyReports lists the archived monthly reports, grouped into a folder per billing entity.
// Only entities that have at least one report get a folder.
//
//encore:api auth method=GET path=/monthly-reports tag:admin
func (s *Service) ListMonthlyReports(ctx context.Context) (*ListMonthlyReportsResponse, error) {
	filesByEntity, err := listArchivedMonthlyReports(ctx)
	if err != nil {
		rlog.Error("failed to list archived monthly reports", "error", err)
		return nil, api_errors.ErrInternalError
	}

	names, err := billingEntityNames(ctx)
	if err != nil {
		return nil, err
	}

	folders := make([]MonthlyReportFolder, 0, len(filesByEntity))
	for entity, files := range filesByEntity {
		// Newest period first, the way a folder of monthly files is usually read.
		sort.Slice(files, func(i, j int) bool { return files[i].Period > files[j].Period })

		name, ok := names[entity]
		if !ok {
			// The entity was deleted or turned organic since the report was archived; keep the
			// folder reachable under its id rather than hiding files an admin may still need.
			name = fmt.Sprintf("#%d", entity.ID)
		}

		folders = append(folders, MonthlyReportFolder{
			EntityType: entity.Type,
			EntityID:   entity.ID,
			Name:       name,
			Files:      files,
		})
	}

	sort.Slice(folders, func(i, j int) bool { return folders[i].Name < folders[j].Name })

	return &ListMonthlyReportsResponse{Folders: folders}, nil
}

// listArchivedMonthlyReports groups every object under the monthly reports prefix by the billing
// entity that owns it. Keys that don't match the report layout are skipped: the bucket also holds
// other private files.
func listArchivedMonthlyReports(ctx context.Context) (map[billingEntity][]MonthlyReportFile, error) {
	filesByEntity := make(map[billingEntity][]MonthlyReportFile)

	query := &objects.Query{Prefix: monthlyReportsPrefix + "/"}
	for entry, err := range storage.PrivateFiles.List(ctx, query) {
		if err != nil {
			return nil, err
		}

		entity, period, err := parseMonthlyReportKey(entry.Name)
		if err != nil {
			rlog.Warn("skipping unrecognized object under the monthly reports prefix", "error", err)
			continue
		}

		filesByEntity[entity] = append(filesByEntity[entity], MonthlyReportFile{
			Period: period,
			Size:   entry.Size,
		})
	}

	return filesByEntity, nil
}

// billingEntityNames resolves the display name of every billing entity: accounts owns the names,
// and its two billing-entity listings are exactly the organizations and offices reports are filed
// under.
func billingEntityNames(ctx context.Context) (map[billingEntity]string, error) {
	organizations, err := accounts.ListOrganicOrganizations(ctx)
	if err != nil {
		rlog.Error("failed to list organic organizations for monthly reports", "error", err)
		return nil, api_errors.ErrInternalError
	}

	offices, err := accounts.ListInorganicOffices(ctx)
	if err != nil {
		rlog.Error("failed to list inorganic offices for monthly reports", "error", err)
		return nil, api_errors.ErrInternalError
	}

	names := make(map[billingEntity]string, len(organizations.Organizations)+len(offices.Offices))
	for _, o := range organizations.Organizations {
		names[billingEntity{Type: organizationsSegment, ID: o.ID}] = o.Name
	}
	for _, o := range offices.Offices {
		names[billingEntity{Type: officesSegment, ID: o.ID}] = o.Name
	}

	return names, nil
}
