package booking

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
	locations_mocks "encore.app/services/booking/mocks"
	"go.uber.org/mock/gomock"
)

// mockBroker implements broker.Broker for testing insertLocations without HTTP.
type mockBroker struct {
	name  broker.Name
	pages []broker.LocationPage
	errs  map[string]error // cursor -> error
	calls int
}

func (m *mockBroker) Name() broker.Name { return m.name }

func (m *mockBroker) GetLocationsPage(cursor string) (broker.LocationPage, error) {
	if err, ok := m.errs[cursor]; ok {
		page := broker.LocationPage{}
		if m.calls < len(m.pages) {
			page = m.pages[m.calls]
		}
		m.calls++
		return page, err
	}
	if m.calls >= len(m.pages) {
		return broker.LocationPage{}, nil
	}
	page := m.pages[m.calls]
	m.calls++
	return page, nil
}

func TestInsertLocations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()

	t.Run("inserts single page of locations with IATA", func(t *testing.T) {
		loc := broker.Location{
			ID: "loc-iata-1", Name: "Airport", Country: "US Country",
			CountryCode: "US", City: "New York", Iata: "JFK",
		}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify location was inserted
		locID, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "US Country", CountryCode: "US", City: "New York",
			Name: "Airport", Iata: "JFK",
		})
		if err != nil {
			t.Fatalf("failed to query location: %v", err)
		}
		if locID == 0 {
			t.Fatal("location was not inserted")
		}

		// Verify broker code
		brokerCode, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerFlex, BrokerLocationID: "loc-iata-1", LocationID: locID,
		})
		if err != nil {
			t.Fatalf("failed to query broker code: %v", err)
		}
		if brokerCode.BrokerLocationID != "loc-iata-1" {
			t.Fatalf("expected broker location ID %q, got %q", "loc-iata-1", brokerCode.BrokerLocationID)
		}
	})

	t.Run("inserts location without IATA uses UpsertByCountryCodeName", func(t *testing.T) {
		loc := broker.Location{
			ID: "loc-no-iata-1", Name: "Downtown Office", Country: "France",
			CountryCode: "FR", City: "Paris", Iata: "",
		}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify via country code name upsert (returns existing ID)
		locID, err := q.UpsertLocationByCountryCodeName(ctx, db.UpsertLocationByCountryCodeNameParams{
			Country: "France", CountryCode: "FR", City: "Paris", Name: "Downtown Office",
		})
		if err != nil {
			t.Fatalf("failed to query location: %v", err)
		}
		if locID == 0 {
			t.Fatal("location was not inserted")
		}

		brokerCode, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerFlex, BrokerLocationID: "loc-no-iata-1", LocationID: locID,
		})
		if err != nil {
			t.Fatalf("failed to query broker code: %v", err)
		}
		if brokerCode.BrokerLocationID != "loc-no-iata-1" {
			t.Fatalf("expected broker location ID %q, got %q", "loc-no-iata-1", brokerCode.BrokerLocationID)
		}
	})

	t.Run("handles multiple pages", func(t *testing.T) {
		loc1 := broker.Location{ID: "mp-loc-1", Name: "Location A", Country: "Country", CountryCode: "US", City: "CityA", Iata: "AAA"}
		loc2 := broker.Location{ID: "mp-loc-2", Name: "Location B", Country: "Country", CountryCode: "US", City: "CityB", Iata: "BBB"}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc1}, NextPage: "page2"},
				{Locations: []broker.Location{loc2}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify both locations exist
		id1, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "Country", CountryCode: "US", City: "CityA", Name: "Location A", Iata: "AAA",
		})
		if err != nil || id1 == 0 {
			t.Fatalf("location A not found: %v", err)
		}

		id2, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "Country", CountryCode: "US", City: "CityB", Name: "Location B", Iata: "BBB",
		})
		if err != nil || id2 == 0 {
			t.Fatalf("location B not found: %v", err)
		}
	})

	t.Run("skips broker code insert when it already exists", func(t *testing.T) {
		loc := broker.Location{ID: "dup-loc-1", Name: "Dup Location", Country: "Country", CountryCode: "DE", City: "Berlin", Iata: "DUP"}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		// Insert once
		ls := location.NewLocationService(q)
		err := ls.InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("first insert failed: %v", err)
		}

		// Insert again — should not error (broker code already exists)
		b.calls = 0
		err = ls.InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("second insert failed: %v", err)
		}
	})

	t.Run("returns error when first page fails with no next page", func(t *testing.T) {
		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: nil, NextPage: ""},
			},
			errs: map[string]error{"": errors.New("broker down")},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("skips page on error and continues to next", func(t *testing.T) {
		loc := broker.Location{ID: "skip-loc-1", Name: "Skip B", Country: "Country", CountryCode: "US", City: "CityY", Iata: "SKP"}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: nil, NextPage: "page2"},
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
			errs: map[string]error{"": errors.New("temporary error")},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify the second page's location was inserted
		locID, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "Country", CountryCode: "US", City: "CityY", Name: "Skip B", Iata: "SKP",
		})
		if err != nil || locID == 0 {
			t.Fatalf("skipped location not found: %v", err)
		}
	})

	t.Run("returns error when upsert fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := locations_mocks.NewMockQuerier(ctrl)
		loc := broker.Location{ID: "err-loc-1", Name: "A", Country: "C", CountryCode: "US", City: "X", Iata: "ERR"}

		q.EXPECT().UpsertLocationByIATA(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when InsertLocationBrokerCode fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := locations_mocks.NewMockQuerier(ctrl)
		loc := broker.Location{ID: "err-loc-2", Name: "A", Country: "C", CountryCode: "US", City: "X", Iata: "EBC"}

		q.EXPECT().UpsertLocationByIATA(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		q.EXPECT().GetLocationBrokerCode(gomock.Any(), gomock.Any()).Return(db.LocationBrokerCode{}, db.ErrNoRows)
		q.EXPECT().InsertLocationBrokerCode(gomock.Any(), gomock.Any()).Return(db.LocationBrokerCode{}, errors.New("insert failed"))

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when GetLocationBrokerCode returns unexpected error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := locations_mocks.NewMockQuerier(ctrl)
		loc := broker.Location{ID: "err-loc-3", Name: "A", Country: "C", CountryCode: "US", City: "X", Iata: "EGT"}

		q.EXPECT().UpsertLocationByIATA(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		q.EXPECT().GetLocationBrokerCode(gomock.Any(), gomock.Any()).Return(db.LocationBrokerCode{}, errors.New("unexpected db error"))

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("handles empty page gracefully", func(t *testing.T) {
		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("stops when broker returns same cursor to avoid infinite loop", func(t *testing.T) {
		loc := broker.Location{ID: "loop-loc-1", Name: "Loop", Country: "C", CountryCode: "US", City: "X", Iata: "LOP"}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: "same-cursor"},
				{Locations: []broker.Location{loc}, NextPage: "same-cursor"},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err == nil {
			t.Fatal("expected error for same cursor, got nil")
		}
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("normalizes IATA with whitespace and lowercase", func(t *testing.T) {
		loc := broker.Location{ID: "norm-loc-1", Name: "Normalized", Country: "C", CountryCode: "US", City: "X", Iata: " jfk "}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify IATA was normalized to "JFK"
		locID, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "C", CountryCode: "US", City: "X", Name: "Normalized", Iata: "JFK",
		})
		if err != nil || locID == 0 {
			t.Fatalf("normalized location not found: %v", err)
		}
	})

	t.Run("invalid IATA length falls back to UpsertByCountryCodeName", func(t *testing.T) {
		loc := broker.Location{ID: "short-iata-1", Name: "Short IATA", Country: "C", CountryCode: "US", City: "X", Iata: "AB"}

		b := &mockBroker{
			name: broker.BrokerFlex,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify location was inserted via country code name fallback
		locID, err := q.UpsertLocationByCountryCodeName(ctx, db.UpsertLocationByCountryCodeNameParams{
			Country: "C", CountryCode: "US", City: "X", Name: "Short IATA",
		})
		if err != nil || locID == 0 {
			t.Fatalf("location with short IATA not found: %v", err)
		}
	})

	t.Run("works with hertz broker", func(t *testing.T) {
		loc := broker.Location{ID: "hertz-1", Name: "Hertz Office", Country: "Germany", CountryCode: "DE", City: "Berlin", Iata: "BER"}

		b := &mockBroker{
			name: broker.BrokerHertz,
			pages: []broker.LocationPage{
				{Locations: []broker.Location{loc}, NextPage: ""},
			},
		}

		err := location.NewLocationService(q).InsertLocations(ctx, b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		locID, err := q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
			Country: "Germany", CountryCode: "DE", City: "Berlin", Name: "Hertz Office", Iata: "BER",
		})
		if err != nil || locID == 0 {
			t.Fatalf("hertz location not found: %v", err)
		}

		brokerCode, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerHertz, BrokerLocationID: "hertz-1", LocationID: locID,
		})
		if err != nil {
			t.Fatalf("hertz broker code not found: %v", err)
		}
		if brokerCode.Broker != db.BrokerHertz {
			t.Fatalf("expected broker %q, got %q", db.BrokerHertz, brokerCode.Broker)
		}
	})
}

func TestExtractFile(t *testing.T) {
	t.Run("returns error for non-multipart content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", nil)
		req.Header.Set("Content-Type", "application/json")

		_, err := location.ExtractFile(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, location.ErrInvalidContentType, err)
	})

	t.Run("returns error when no file field in form", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, err := location.ExtractFile(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, location.ErrGetFileFromForm, err)
	})

	t.Run("extracts file successfully", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "locations.csv")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("test file content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		file, err := location.ExtractFile(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer file.Close()

		buf := make([]byte, 100)
		n, _ := file.Read(buf)
		if string(buf[:n]) != "test file content" {
			t.Fatalf("expected file content 'test file content', got %q", string(buf[:n]))
		}
	})
}

func TestInsertLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("validate the payload", func(t *testing.T) {
		cases := []struct {
			name            string
			params          location.InsertLocationParams
			isExpectedError bool
		}{
			{
				name: "valid payload",
				params: location.InsertLocationParams{
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
				params: location.InsertLocationParams{
					ID:          "flex-loc-1",
					Name:        "Test Location",
					Country:     "Test Country",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "invalid broker",
				params: location.InsertLocationParams{
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
				params: location.InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Country:     "Test Country",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "missing country",
				params: location.InsertLocationParams{
					Broker:      broker.BrokerFlex,
					ID:          "flex-loc-1",
					Name:        "Test Location",
					CountryCode: "US",
				},
				isExpectedError: true,
			},
			{
				name: "country code too short",
				params: location.InsertLocationParams{
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
				params: location.InsertLocationParams{
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
				params: location.InsertLocationParams{
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
				params: location.InsertLocationParams{
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
		err := s.InsertLocation(ctx, location.InsertLocationParams{
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
		params := location.InsertLocationParams{
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
