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

var secrets struct {
	translationToken string
}

// GetPendingTranslationsResponse is the response for GetPendingTranslations endpoint
type GetPendingTranslationsResponse struct {
	Translations []db.BrokerTranslation `json:"translations"`
}

// GetPendingTranslationsRequest is the request for GetPendingTranslations endpoint
type GetPendingTranslationsRequest struct {
	Token string `header:"X-Translation-Token" encore:"sensitive"`
}

// GetPendingTranslations returns the list of pending translations for brokers. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/pending method=GET
func (s *Service) GetPendingTranslations(ctx context.Context, p *GetPendingTranslationsRequest) (*GetPendingTranslationsResponse, error) {
	if p.Token != secrets.translationToken {
		rlog.Warn("invalid translation token", "provided_token", p.Token)
		return nil, api_errors.ErrNotFound
	}

	ts, err := s.query.ListPendingTranslations(ctx)
	if err != nil {
		rlog.Error("failed to get pending translations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetPendingTranslationsResponse{
		Translations: ts,
	}, nil
}

// TranslateTranslationRequest is the request for TranslateTranslation endpoint
type TranslateTranslationRequest struct {
	Token      string `header:"X-Translation-Token" encore:"sensitive"`
	ID         int32  `json:"id"`
	TargetText string `json:"targetText"`
	Confidence int32  `json:"confidence"`
}

// TranslateTranslation translates a pending translation. It requires a valid translation token in the header.
// encore: api public path=/booking/translations/translate method=PATCH
func (s *Service) TranslateTranslation(ctx context.Context, p *TranslateTranslationRequest) error {
	if p.Token != secrets.translationToken {
		rlog.Warn("invalid translation token", "provided_token", p.Token)
		return api_errors.ErrNotFound
	}

	err := s.query.TranslatePendingTranslation(ctx, db.TranslatePendingTranslationParams{
		ID:              p.ID,
		TargetText:      &p.TargetText,
		ConfidenceScore: &p.Confidence,
	})

	if err != nil {
		rlog.Error("failed to translate translation", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}

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

type UpdateBrokerTranslationRequest struct {
	TargetText string `json:"target_text" validate:"required,notblank"`
}

func (p UpdateBrokerTranslationRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdateBrokerTranslation updates a broker translation target text by ID.
//
//encore:api auth method=PUT path=/broker-translations/:id tag:admin
func (s *Service) UpdateBrokerTranslation(ctx context.Context, id int32, params UpdateBrokerTranslationRequest) error {
	err := s.query.UpdateBrokerTranslation(ctx, db.UpdateBrokerTranslationParams{
		ID:         id,
		TargetText: &params.TargetText,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to update broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}

// VerifyBrokerTranslation marks a broker translation as verified by ID.
//
//encore:api auth method=PATCH path=/broker-translations/:id/verify tag:admin
func (s *Service) VerifyBrokerTranslation(ctx context.Context, id int32) error {
	err := s.query.VerifyBrokerTranslation(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to verify broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}

// DeleteBrokerTranslation deletes a broker translation by ID.
//
//encore:api auth method=DELETE path=/broker-translations/:id tag:admin
func (s *Service) DeleteBrokerTranslation(ctx context.Context, id int32) error {
	err := s.query.DeleteBrokerTranslation(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
