package reservation

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type nameEntry struct {
	Title     string
	FirstName string
	LastName  string
}

var fakeNames = []nameEntry{
	{"Mr", "James", "Smith"},
	{"Ms", "Sarah", "Johnson"},
	{"Mr", "Michael", "Williams"},
	{"Ms", "Emily", "Brown"},
	{"Mr", "David", "Jones"},
	{"Ms", "Laura", "Garcia"},
	{"Mr", "Robert", "Miller"},
	{"Ms", "Anna", "Davis"},
	{"Mr", "Daniel", "Martinez"},
	{"Ms", "Sophie", "Anderson"},
	{"Mr", "Thomas", "Wilson"},
	{"Ms", "Rachel", "Taylor"},
	{"Mr", "Andrew", "Thomas"},
	{"Ms", "Olivia", "Moore"},
	{"Mr", "Peter", "Jackson"},
	{"Ms", "Emma", "White"},
	{"Mr", "George", "Harris"},
	{"Ms", "Natalie", "Martin"},
	{"Mr", "Steven", "Thompson"},
	{"Ms", "Jessica", "Robinson"},
	{"Mr", "Mark", "Clark"},
	{"Ms", "Hannah", "Lewis"},
	{"Mr", "Paul", "Walker"},
	{"Ms", "Victoria", "Hall"},
	{"Mr", "Kevin", "Allen"},
	{"Ms", "Charlotte", "Young"},
	{"Mr", "Brian", "King"},
	{"Ms", "Rebecca", "Wright"},
	{"Mr", "Jason", "Scott"},
	{"Ms", "Diana", "Green"},
}

type SeedReservationsResponse struct {
	CreatedIDs []int64 `json:"createdIds"`
}

// SeedReservations creates 30 fake reservations based on reservations 1-4.
//
//encore:api private
func (s *Service) SeedReservations(ctx context.Context) (*SeedReservationsResponse, error) {
	sourceIDs := []int64{1, 4}

	// Fetch source reservations
	var sources []db.GetReservationByIDRow
	for _, id := range sourceIDs {
		row, err := s.query.GetReservationByID(ctx, id)
		if err != nil {
			rlog.Error("failed to get source reservation", "id", id, "error", err)
			return nil, api_errors.ErrInternalError
		}
		sources = append(sources, row)
	}

	// All 30 share the user_id from reservation 1
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var createdIDs []int64
	for i := 0; i < 30; i++ {
		src := sources[i%len(sources)]
		name := fakeNames[i%len(fakeNames)]

		// Generate random pickup date between 2026-04-01 and 2026-05-25
		startRange := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
		daysRange := 55 // ~April 1 to May 25
		pickupDate := startRange.AddDate(0, 0, rng.Intn(daysRange))

		// Rental days between 2 and 14
		rentalDays := 2 + rng.Intn(13)
		returnDate := pickupDate.AddDate(0, 0, rentalDays)

		// Ensure return date doesn't exceed end of June
		maxReturn := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
		if returnDate.After(maxReturn) {
			returnDate = maxReturn
			rentalDays = int(returnDate.Sub(pickupDate).Hours() / 24)
			if rentalDays < 1 {
				rentalDays = 1
				returnDate = pickupDate.AddDate(0, 0, 1)
			}
		}

		brokerResID := src.BrokerReservationID + "FAKE"

		id, err := s.query.InsertReservation(ctx, db.InsertReservationParams{
			UserID:              src.UserID,
			BrokerReservationID: brokerResID,
			Broker:              src.Broker,
			SupplierCode:        src.SupplierCode,
			CarDetails:          src.CarDetails,
			PlanInclusions:      src.PlanInclusions,
			CountryCode:         src.CountryCode,
			CurrencyCode:        src.CurrencyCode,
			CurrencyRate:        src.CurrencyRate,
			PurchasePrice:       src.PurchasePrice,
			MarkupPercentage:    src.MarkupPercentage,
			DiscountPercentage:  src.DiscountPercentage,
			BrokerErpPrice:      src.BrokerErpPrice,
			BtErpPrice:          src.BtErpPrice,
			VatPercentage:       src.VatPercentage,
			TotalPrice:          src.TotalPrice,
			PickupDate:          db.DateFromString(pickupDate.Format("2006-01-02")),
			ReturnDate:          db.DateFromString(returnDate.Format("2006-01-02")),
			PickupTime:          src.PickupTime,
			DropoffTime:         src.DropoffTime,
			RentalDays:          int32(rentalDays),
			DriverTitle:         name.Title,
			DriverFirstName:     name.FirstName,
			DriverLastName:      name.LastName,
			DriverAge:           src.DriverAge,
			PickupLocationName:  src.PickupLocationName,
			DropoffLocationName: src.DropoffLocationName,
		})
		if err != nil {
			rlog.Error("failed to insert seed reservation", "index", i, "error", err)
			return nil, fmt.Errorf("failed to insert reservation %d: %w", i, err)
		}

		createdIDs = append(createdIDs, id)
	}

	rlog.Info("seeded reservations", "count", len(createdIDs))
	return &SeedReservationsResponse{CreatedIDs: createdIDs}, nil
}
