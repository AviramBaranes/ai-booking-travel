package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"go.uber.org/mock/gomock"
)

func TestDeleteLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes broker code and keeps location when other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityTLV := "Tel Aviv"
		iataTLV := "TLV"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", City: &cityTLV, Name: "Ben Gurion Airport", Iata: &iataTLV,
			},
			db.BrokerFlex, "flex-tlv-del-keep",
		)

		// Add a second broker code to the same location
		_, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
			LocationID:       loc.ID,
			Broker:           db.BrokerHertz,
			BrokerLocationID: "hertz-tlv-del-keep",
		})
		if err != nil {
			t.Fatalf("failed to insert second broker code: %v", err)
		}

		// Delete the first broker code
		err = s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should still exist
		_, err = q.GetLocationById(ctx, loc.ID)
		if err != nil {
			t.Fatalf("location should still exist, got error: %v", err)
		}
	})

	t.Run("deletes broker code and location when no other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityParis := "Paris"
		iataCDG := "CDG"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "France", CountryCode: "FR", City: &cityParis, Name: "Charles de Gaulle", Iata: &iataCDG,
			},
			db.BrokerFlex, "flex-cdg-del-orphan",
		)

		err := s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should be deleted
		_, err = q.GetLocationById(ctx, loc.ID)
		if !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("expected location to be deleted, got err: %v", err)
		}
	})

	t.Run("returns not found for non-existent broker code", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		err := s.DeleteLocation(ctx, 999999)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when delete broker code fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(0), errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(10), nil)
		q.EXPECT().CountLocationBrokerCodesByLocationID(gomock.Any(), int64(10)).
			Return(int64(0), errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when delete location fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(10), nil)
		q.EXPECT().CountLocationBrokerCodesByLocationID(gomock.Any(), int64(10)).
			Return(int64(0), nil)
		q.EXPECT().DeleteLocationByID(gomock.Any(), int64(10)).
			Return(errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
