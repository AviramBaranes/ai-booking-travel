package price_offer

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// ListPriceOffersRequest holds the query parameters for listing price offers.
type ListPriceOffersRequest struct {
	Name   string `query:"name" encore:"optional"`
	Status string `query:"status" encore:"optional"`
	Page   int32  `query:"page" validate:"required,gte=1"`
}

func (p ListPriceOffersRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// PriceOfferSummary is a compact price offer representation for list views.
type PriceOfferSummary struct {
	ID                  int64  `json:"id"`
	Status              string `json:"status"`
	Name                string `json:"name"`
	PickupLocationName  string `json:"pickupLocationName"`
	DropoffLocationName string `json:"dropoffLocationName"`
	PickupDate          string `json:"pickupDate"`
	DropoffDate         string `json:"dropoffDate"`
	PickupTime          string `json:"pickupTime"`
	DropoffTime         string `json:"dropoffTime"`
	CurrencyCode        string `json:"currencyCode"`
	TotalPrice          int32  `json:"totalPrice"`
	OfferedCurrencyCode string `json:"offeredCurrencyCode"`
	OfferedPrice        int32  `json:"offeredPrice"`
	CreatedAt           string `json:"createdAt"`
}

// ListPriceOffersResponse wraps a page of price offer summaries with a total count.
type ListPriceOffersResponse struct {
	PriceOffers []PriceOfferSummary `json:"priceOffers"`
	Total       int64               `json:"total"`
}

const listPriceOffersLimit int32 = 8

// ListPriceOffers returns a paginated list of the authenticated agent's price offers.
func (s *PriceOfferService) ListPriceOffers(ctx context.Context, p ListPriceOffersRequest) (*ListPriceOffersResponse, error) {
	rows, err := s.listPriceOffersByAgent(ctx, p)
	if err != nil {
		return nil, err
	}

	total, err := s.countPriceOffersByAgent(ctx, p)
	if err != nil {
		return nil, err
	}

	return &ListPriceOffersResponse{
		PriceOffers: mapRowsToPriceOfferSummaries(rows),
		Total:       total,
	}, nil
}

func (s *PriceOfferService) listPriceOffersByAgent(ctx context.Context, p ListPriceOffersRequest) ([]db.ListPriceOffersByAgentRow, error) {
	authData := auth.GetAuthData()
	offset := (p.Page - 1) * listPriceOffersLimit

	rows, err := s.query.ListPriceOffersByAgent(ctx, db.ListPriceOffersByAgentParams{
		AgentID:    authData.UserID,
		Status:     nullOfferStatusFromString(p.Status),
		NameSearch: nilIfEmpty(p.Name),
		PageSize:   listPriceOffersLimit,
		PageOffset: offset,
	})
	if err != nil {
		if isNotFound(err) {
			return []db.ListPriceOffersByAgentRow{}, nil
		}
		rlog.Error("failed to list price offers", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return rows, nil
}

func (s *PriceOfferService) countPriceOffersByAgent(ctx context.Context, p ListPriceOffersRequest) (int64, error) {
	authData := auth.GetAuthData()

	count, err := s.query.CountPriceOffersByAgent(ctx, db.CountPriceOffersByAgentParams{
		AgentID:    authData.UserID,
		Status:     nullOfferStatusFromString(p.Status),
		NameSearch: nilIfEmpty(p.Name),
	})
	if err != nil {
		if isNotFound(err) {
			return 0, nil
		}
		rlog.Error("failed to count price offers", "error", err)
		return 0, api_errors.ErrInternalError
	}

	return count, nil
}

func mapRowsToPriceOfferSummaries(rows []db.ListPriceOffersByAgentRow) []PriceOfferSummary {
	summaries := make([]PriceOfferSummary, len(rows))
	for i, r := range rows {
		summaries[i] = PriceOfferSummary{
			ID:                  r.ID,
			Status:              string(r.Status),
			Name:                r.Name,
			PickupLocationName:  r.PickupLocation,
			DropoffLocationName: r.DropoffLocation,
			PickupDate:          dbadapters.DateToString(r.PickupDate),
			DropoffDate:         dbadapters.DateToString(r.DropoffDate),
			PickupTime:          r.PickupTime,
			DropoffTime:         r.DropoffTime,
			CurrencyCode:        r.CurrencyCode,
			TotalPrice:          r.TotalPrice,
			OfferedCurrencyCode: r.OfferedCurrencyCode,
			OfferedPrice:        r.OfferedPrice,
			CreatedAt:           dbadapters.TimestamptzToString(r.CreatedAt),
		}
	}
	return summaries
}
