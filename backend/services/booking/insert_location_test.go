package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	locations_mocks "encore.app/services/booking/mocks"
	"go.uber.org/mock/gomock"
)

func TestInsertLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("validate the payload", func(t *testing.T) {
		cases := []struct {
			name            string
			params          InsertLocationParams
			isExpectedError bool
		}{
			{
				name: "valid payload",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
					City:        "Test City",
					Iata:        "TES",
				},
				isExpectedError: false,
			},
			{
				name: "missing broker",
				params: InsertLocationParams{
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "invalid broker",
				params: InsertLocationParams{
					Broker:      "invalid-broker",
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "missing name",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Country:     "Test Country",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "missing country",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "country code too short",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "U",
				},
				isExpectedError: true,
			},
			{
				name: "country code too long",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "USA",
				},
				isExpectedError: true,
			},
			{
				name: "iata code too short",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
					Iata:        "TE",
				},
				isExpectedError: true,
			},
			{
				name: "iata code too long",
				params: InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
					Iata:        "TEST",
				},
				isExpectedError: true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.params.Validate()
				if tc.isExpectedError && err == nil {
					t.Fatalf("expected validation error, got nil")
				}
				if !tc.isExpectedError && err != nil {
					t.Fatalf("expected no validation error, got %v", err)
				}
			})
		}
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := locations_mocks.NewMockQuerier(ctrl)

		q.EXPECT().UpsertLocationByIATA(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("database error")).Times(1)

		s := &Service{query: q}
		err := s.InsertLocation(ctx, InsertLocationParams{
			Broker:      broker.BrokerFlex,
			ID:          "flex-loc-err",
			Name:        "Error Location",
			Country:     "Error Country",
			CountryCode: "ER",
			Iata:        "ERR",
		})

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("inserts location successfully", func(t *testing.T) {
		params := InsertLocationParams{
			Broker:      broker.BrokerFlex,
			ID:          "flex-loc-success",
			Name:        "Success Location",
			Country:     "Success Country",
			CountryCode: "SU",
			City:        "Success City",
			Iata:        "SUC",
		}

		err := InsertLocation(ctx, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify using the database
		locID, err := query.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Iata:        params.Iata,
			Name:        params.Name,
			Country:     params.Country,
			CountryCode: params.CountryCode,
			City:        params.City,
		})
		if err != nil {
			t.Fatalf("failed to query inserted location: %v", err)
		}
		if locID == 0 {
			t.Fatalf("location was not inserted")
		}

		// Verify broker code
		brokerCode, err := query.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			LocationID: locID,
			Broker:     db.BrokerFlex,
		})
		if err != nil {
			t.Fatalf("failed to query inserted broker code: %v", err)
		}
		if brokerCode.BrokerLocationID != params.ID {
			t.Fatalf("expected broker location ID %s, got %s", params.ID, brokerCode.BrokerLocationID)
		}
	})
}
