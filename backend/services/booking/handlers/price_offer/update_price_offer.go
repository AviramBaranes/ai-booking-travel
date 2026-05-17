package price_offer

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// UpdatePriceOfferParams holds the mutable fields for a price offer update.
type UpdatePriceOfferParams struct {
	Status              *string `json:"status" encore:"optional" validate:"omitempty,oneof=open booked declined"`
	Name                *string `json:"name" encore:"optional" validate:"omitempty,notblank"`
	OfferedCurrencyCode *string `json:"offeredCurrencyCode" encore:"optional" validate:"omitempty,len=3,uppercase_only"`
	OfferedPrice        *int32  `json:"offeredPrice" encore:"optional" validate:"omitempty,gt=0"`
}

func (p UpdatePriceOfferParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdatePriceOffer updates a price offer's mutable fields for the authenticated agent.
func (s *PriceOfferService) UpdatePriceOffer(ctx context.Context, id int64, p UpdatePriceOfferParams) error {
	authData := auth.GetAuthData()

	var status db.NullOfferStatus
	if p.Status != nil {
		status = nullOfferStatusFromString(*p.Status)
	}

	err := s.query.UpdatePriceOffer(ctx, db.UpdatePriceOfferParams{
		ID:                  id,
		AgentID:             authData.UserID,
		Status:              status,
		Name:                p.Name,
		OfferedCurrencyCode: p.OfferedCurrencyCode,
		OfferedPrice:        p.OfferedPrice,
	})
	if err != nil {
		if isNotFound(err) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to update price offer", "id", id, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
