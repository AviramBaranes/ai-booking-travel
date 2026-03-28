package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

const TranslationsLimit = 15

type ListBrokerTranslationsRequest struct {
	Page    int    `query:"page" validate:"required,min=1"`
	Search  string `query:"search"`
	Status  string `query:"status"`
	SortDir string `query:"sortDir" validate:"omitempty,oneof=asc desc"`
}

func (r ListBrokerTranslationsRequest) Validate() error {
	if r.Status != "" &&
		r.Status != string(db.BrokerTranslationStatusPending) &&
		r.Status != string(db.BrokerTranslationStatusTranslated) &&
		r.Status != string(db.BrokerTranslationStatusVerified) {
		return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code: api_errors.CodeInvalidValue, Field: "status",
		})
	}
	return validation.ValidateStruct(r)
}

type BrokerTranslationRow struct {
	ID              int32   `json:"id"`
	SourceText      string  `json:"source_text"`
	TargetText      *string `json:"target_text"`
	Status          string  `json:"status"`
	ConfidenceScore *int32  `json:"confidence_score"`
}

type ListBrokerTranslationsResponse struct {
	Translations []BrokerTranslationRow `json:"translations"`
	Total        int64                  `json:"total"`
}

//encore:api auth method=GET path=/broker-translations tag:admin
func (s *Service) ListBrokerTranslations(ctx context.Context, p ListBrokerTranslationsRequest) (*ListBrokerTranslationsResponse, error) {
	search := nilIfEmpty(p.Search)
	status := db.NullBrokerTranslationStatusFromString(p.Status)
	sortDir := defaultSortDir(p.SortDir)

	total, err := s.query.CountAllTranslations(ctx, db.CountAllTranslationsParams{
		Search: search,
		Status: status,
	})
	if err != nil {
		rlog.Error("failed to count broker translations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rows, err := s.query.ListAllTranslations(ctx, db.ListAllTranslationsParams{
		Search:      search,
		Status:      status,
		SortDir:     sortDir,
		QueryOffset: int32((p.Page - 1) * TranslationsLimit),
		QueryLimit:  TranslationsLimit,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListBrokerTranslationsResponse{Translations: []BrokerTranslationRow{}}, nil
		}
		rlog.Error("failed to list broker translations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &ListBrokerTranslationsResponse{
		Translations: toBrokerTranslationRows(rows),
		Total:        total,
	}, nil
}

func toBrokerTranslationRows(rows []db.BrokerTranslation) []BrokerTranslationRow {
	result := make([]BrokerTranslationRow, len(rows))
	for i, row := range rows {
		result[i] = BrokerTranslationRow{
			ID:              row.ID,
			SourceText:      row.SourceText,
			TargetText:      row.TargetText,
			Status:          string(row.Status),
			ConfidenceScore: row.ConfidenceScore,
		}
	}
	return result
}

func defaultSortDir(dir string) string {
	if dir == "" {
		return "asc"
	}
	return dir
}
