package reservation

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/mocks"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

// testQuerier returns a real db.Querier backed by the Encore test database.
func testQuerier() *db.Queries {
	pool := sqldb.Driver[*pgxpool.Pool](reservationsDB)
	return db.New(pool)
}

func validCreateReservationParams() *CreateReservationRequest {
	return &CreateReservationRequest{
		UserID:              1,
		BrokerReservationID: "BRK-12345",
		Broker:              "flex",
		SupplierCode:        "SUP1",
		CarDetails: &broker.CarDetails{
			Model:        "Toyota Corolla",
			CarGroup:     "Economy",
			ImageURL:     "https://example.com/car.png",
			SupplierName: "Hertz",
			CarType:      "Sedan",
			Acriss:       "CDMR",
			HasAC:        true,
			IsAutoGear:   true,
			IsElectric:   false,
			Seats:        5,
			Bags:         2,
			Doors:        4,
		},
		PlanInclusions:      []string{"Unlimited Mileage", "Collision Damage Waiver"},
		CountryCode:         "US",
		CurrencyCode:        "USD",
		CurrencyRate:        3.65,
		PurchasePrice:       100.00,
		MarkupPercentage:    45.00,
		DiscountPercentage:  10,
		BrokerErpPrice:      15.00,
		BtErpPrice:          20,
		PickupDate:          "2026-04-01",
		ReturnDate:          "2026-04-05",
		RentalDays:          4,
		DriverTitle:         "Mr",
		DriverFirstName:     "John",
		DriverLastName:      "Doe",
		DriverAge:           30,
		PickupLocationName:  "Airport Terminal 1",
		DropoffLocationName: "City Center Office",
	}
}

func mockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func invalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// --- Tests ---

func TestCreateReservationValidation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(p *CreateReservationRequest)
		wantErr error
	}{
		{
			name:    "rejects missing broker reservation id",
			modify:  func(p *CreateReservationRequest) { p.BrokerReservationID = "" },
			wantErr: invalidValueErr("brokerReservationId"),
		},
		{
			name:    "rejects blank broker reservation id",
			modify:  func(p *CreateReservationRequest) { p.BrokerReservationID = "   " },
			wantErr: invalidValueErr("brokerReservationId"),
		},
		{
			name:    "rejects invalid broker",
			modify:  func(p *CreateReservationRequest) { p.Broker = "avis" },
			wantErr: invalidValueErr("broker"),
		},
		{
			name:    "rejects missing supplier code",
			modify:  func(p *CreateReservationRequest) { p.SupplierCode = "" },
			wantErr: invalidValueErr("supplierCode"),
		},
		{
			name:    "rejects missing car details",
			modify:  func(p *CreateReservationRequest) { p.CarDetails = nil },
			wantErr: invalidValueErr("carDetails"),
		},
		{
			name:    "rejects missing plan inclusions",
			modify:  func(p *CreateReservationRequest) { p.PlanInclusions = nil },
			wantErr: invalidValueErr("planInclusions"),
		},
		{
			name:    "rejects missing country code",
			modify:  func(p *CreateReservationRequest) { p.CountryCode = "" },
			wantErr: invalidValueErr("countryCode"),
		},
		{
			name:    "rejects missing currency code",
			modify:  func(p *CreateReservationRequest) { p.CurrencyCode = "" },
			wantErr: invalidValueErr("currencyCode"),
		},
		{
			name:    "rejects zero currency rate",
			modify:  func(p *CreateReservationRequest) { p.CurrencyRate = 0 },
			wantErr: invalidValueErr("currencyRate"),
		},
		{
			name:    "rejects negative currency rate",
			modify:  func(p *CreateReservationRequest) { p.CurrencyRate = -1 },
			wantErr: invalidValueErr("currencyRate"),
		},
		{
			name:    "rejects negative purchase price",
			modify:  func(p *CreateReservationRequest) { p.PurchasePrice = -1 },
			wantErr: invalidValueErr("purchasePrice"),
		},
		{
			name:    "rejects discount above 100",
			modify:  func(p *CreateReservationRequest) { p.DiscountPercentage = 101 },
			wantErr: invalidValueErr("discountPercentage"),
		},
		{
			name:    "rejects negative discount",
			modify:  func(p *CreateReservationRequest) { p.DiscountPercentage = -1 },
			wantErr: invalidValueErr("discountPercentage"),
		},
		{
			name:    "rejects invalid pickup date format",
			modify:  func(p *CreateReservationRequest) { p.PickupDate = "01-04-2026" },
			wantErr: invalidValueErr("pickupDate"),
		},
		{
			name:    "rejects invalid return date format",
			modify:  func(p *CreateReservationRequest) { p.ReturnDate = "2026/04/05" },
			wantErr: invalidValueErr("returnDate"),
		},
		{
			name:    "rejects zero rental days",
			modify:  func(p *CreateReservationRequest) { p.RentalDays = 0 },
			wantErr: invalidValueErr("rentalDays"),
		},
		{
			name:    "rejects missing driver title",
			modify:  func(p *CreateReservationRequest) { p.DriverTitle = "" },
			wantErr: invalidValueErr("driverTitle"),
		},
		{
			name:    "rejects invalid driver title",
			modify:  func(p *CreateReservationRequest) { p.DriverTitle = "Dr" },
			wantErr: invalidValueErr("driverTitle"),
		},
		{
			name:    "rejects missing driver first name",
			modify:  func(p *CreateReservationRequest) { p.DriverFirstName = "" },
			wantErr: invalidValueErr("driverFirstName"),
		},
		{
			name:    "rejects blank driver first name",
			modify:  func(p *CreateReservationRequest) { p.DriverFirstName = "   " },
			wantErr: invalidValueErr("driverFirstName"),
		},
		{
			name:    "rejects missing driver last name",
			modify:  func(p *CreateReservationRequest) { p.DriverLastName = "" },
			wantErr: invalidValueErr("driverLastName"),
		},
		{
			name:    "rejects driver age below 18",
			modify:  func(p *CreateReservationRequest) { p.DriverAge = 17 },
			wantErr: invalidValueErr("driverAge"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := validCreateReservationParams()
			tc.modify(p)
			api_errors.AssertApiError(t, tc.wantErr, p.Validate())
		})
	}

	t.Run("accepts valid params", func(t *testing.T) {
		if err := validCreateReservationParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestCreateReservation(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("creates reservation successfully", func(t *testing.T) {
		resp, err := s.CreateReservation(ctx, *validCreateReservationParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().InsertReservation(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := s.CreateReservation(ctx, *validCreateReservationParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
