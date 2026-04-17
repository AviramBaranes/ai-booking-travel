package reservation

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"go.uber.org/mock/gomock"
)

func TestGetReservation(t *testing.T) {
	ctx := context.Background()
	t.Run("return 404 for non-existent reservation", func(t *testing.T) {
		_, err := GetReservation(ctx, 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})
	t.Run("return sends an internal error on db failure", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().GetReservationByID(gomock.Any(), gomock.Any()).Return(db.GetReservationByIDRow{}, errors.New("db error"))

		_, err := s.GetReservation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
	t.Run("return 404 for reservation of different user", func(t *testing.T) {
		res, err := CreateReservation(ctx, *validCreateReservationParams())
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}

		ctx := authContext(999)
		_, err = GetReservation(ctx, res.ID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns reservation with the right calculations", func(t *testing.T) {
		tests := []struct {
			name   string
			modify func(p *CreateReservationRequest)
			want   GetReservationResponse
		}{
			{
				name: "only car",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 100
					p.MarkupPercentage = 50
					p.BrokerErpPrice = 0
					p.BtErpPrice = 0
					p.DiscountPercentage = 0
				},
				want: GetReservationResponse{
					CarFullPrice:   150,
					ErpPrice:       0,
					TotalPrice:     150,
					DiscountAmount: 0,
				},
			},
			{
				name: "car with bt erp only",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 100
					p.MarkupPercentage = 50
					p.BrokerErpPrice = 0
					p.DiscountPercentage = 0
					p.BtErpPrice = 15
				},
				want: GetReservationResponse{
					CarFullPrice:   150,
					ErpPrice:       15,
					TotalPrice:     165,
					DiscountAmount: 0,
				},
			},
			{
				name: "car with both bt erp and broker erp",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 100
					p.MarkupPercentage = 50
					p.DiscountPercentage = 0
					p.BtErpPrice = 15
					p.BrokerErpPrice = 10
				},
				want: GetReservationResponse{
					CarFullPrice:   150,
					ErpPrice:       30, // markup on the broker erp
					TotalPrice:     180,
					DiscountAmount: 0,
				},
			},
			{
				name: "only car with discount",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 100
					p.MarkupPercentage = 50
					p.DiscountPercentage = 10
					p.BrokerErpPrice = 0
					p.BtErpPrice = 0
				},
				want: GetReservationResponse{
					CarFullPrice:   150,
					ErpPrice:       0,
					TotalPrice:     135,
					DiscountAmount: 15,
				},
			},
			{
				name: "car with bt erp and discount",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 102.5 // so round is up (make sure we not doing floor)
					p.MarkupPercentage = 50
					p.DiscountPercentage = 10
					p.BtErpPrice = 15
					p.BrokerErpPrice = 0
				},
				want: GetReservationResponse{
					CarFullPrice:   154,
					ErpPrice:       15,
					TotalPrice:     153, // rounding (153.75 - 10% discount) + 15 erp (no discount on bt erp)
					DiscountAmount: 16,
				},
			},
			{
				name: "car with both bt erp and broker erp and discount",
				modify: func(p *CreateReservationRequest) {
					p.PurchasePrice = 101.5 // so round is down
					p.MarkupPercentage = 50
					p.DiscountPercentage = 10
					p.BtErpPrice = 15
					p.BrokerErpPrice = 10
				},
				want: GetReservationResponse{
					CarFullPrice:   152,
					ErpPrice:       30,
					TotalPrice:     166,
					DiscountAmount: 16,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				params := validCreateReservationParams()
				tt.modify(params)

				res, err := CreateReservation(ctx, *params)
				if err != nil {
					t.Fatalf("failed to create reservation: %v", err)
				}

				ctx := authContext(params.UserID)

				got, err := GetReservation(ctx, res.ID)
				if err != nil {
					t.Fatalf("failed to get reservation: %v", err)
				}

				if got.CarFullPrice != tt.want.CarFullPrice {
					t.Errorf("unexpected car full price: got %d, want %d", got.CarFullPrice, tt.want.CarFullPrice)
				}
				if got.ErpPrice != tt.want.ErpPrice {
					t.Errorf("unexpected erp price: got %d, want %d", got.ErpPrice, tt.want.ErpPrice)
				}
				if got.TotalPrice != tt.want.TotalPrice {
					t.Errorf("unexpected total price: got %d, want %d", got.TotalPrice, tt.want.TotalPrice)
				}
				if got.DiscountAmount != tt.want.DiscountAmount {
					t.Errorf("unexpected discount amount: got %d, want %d", got.DiscountAmount, tt.want.DiscountAmount)
				}
				if got.PickupTime != params.PickupTime {
					t.Errorf("unexpected pickup time: got %s, want %s", got.PickupTime, params.PickupTime)
				}
				if got.DropoffTime != params.DropoffTime {
					t.Errorf("unexpected dropoff time: got %s, want %s", got.DropoffTime, params.DropoffTime)
				}
			})
		}
	})
}
