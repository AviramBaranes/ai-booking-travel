package booking

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
)

// seedHertzGRStation inserts a location and a Hertz GR station serving the given brands there.
func seedHertzGRStation(t *testing.T, q *db.Queries, loc db.InsertLocationParams, code string, brands []string) (db.Location, db.LocationBrokerCode) {
	t.Helper()
	ctx := context.Background()

	l, err := q.InsertLocation(ctx, loc)
	if err != nil {
		t.Fatalf("failed to seed location: %v", err)
	}
	lbc, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
		LocationID: l.ID, Broker: db.BrokerHertzGr, BrokerLocationID: code, Brands: brands,
	})
	if err != nil {
		t.Fatalf("failed to seed station: %v", err)
	}
	t.Cleanup(func() {
		_, _ = q.DeleteLocationBrokerCode(ctx, lbc.ID)
		_ = q.DeleteLocationByID(ctx, l.ID)
	})

	return l, lbc
}

// TestInsertHertzGRLocations covers what the parser tests cannot reach: that Postgres accepts the
// 'hertz_gr' value and the brands column, and that one location holds a Hertz GR station per brand.
func TestInsertHertzGRLocations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	ls := location.NewLocationService(q)
	suffix := fmt.Sprintf("%d-%d", time.Now().UnixNano(), nextLocSeq())

	hertzDesk := "hgr-hertz-" + suffix
	thriftyDesk := "hgr-thrifty-" + suffix
	clashingDesk := "hgr-clash-" + suffix

	airport, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{
			Country: "Greece", CountryCode: "GR", Name: "Hertz GR Test Airport " + suffix,
			City: strPtr("Mykonos"), Iata: strPtr("ZHG"),
		},
		db.BrokerFlex, "flex-hgr-"+suffix,
	)

	importStations := func(t *testing.T, locs ...broker.Location) {
		t.Helper()
		b := &mockBroker{name: broker.BrokerHertzGR, pages: []broker.LocationPage{{Locations: locs}}}
		if err := ls.InsertLocations(ctx, b, ""); err != nil {
			t.Fatalf("InsertLocations: %v", err)
		}
	}
	station := func(code string, brands ...string) broker.Location {
		return broker.Location{
			ID: code, Name: "Hertz GR Test Airport", Country: "Greece", CountryCode: "GR",
			City: "Mykonos", Iata: "ZHG", Brands: brands,
		}
	}
	codeOf := func(t *testing.T, code string) db.LocationBrokerCode {
		t.Helper()
		c, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerHertzGr, BrokerLocationID: code, LocationID: -1,
		})
		if err != nil {
			t.Fatalf("GetLocationBrokerCode(%s): %v", code, err)
		}
		return c
	}

	importStations(t,
		station(hertzDesk, "ZE"),
		station(thriftyDesk, "ZR", "ZT"),
		station(clashingDesk, "ZT"),
	)

	t.Run("each brand's desk attaches to the existing airport", func(t *testing.T) {
		for code, want := range map[string][]string{hertzDesk: {"ZE"}, thriftyDesk: {"ZR", "ZT"}} {
			c := codeOf(t, code)
			if c.LocationID != airport.ID {
				t.Errorf("%s LocationID = %d, want the existing airport %d", code, c.LocationID, airport.ID)
			}
			if !slices.Equal(c.Brands, want) {
				t.Errorf("%s Brands = %v, want %v", code, c.Brands, want)
			}
		}
	})

	t.Run("a desk serving a brand the location already has is skipped", func(t *testing.T) {
		_, err := q.GetLocationIDByAnyBrokerCode(ctx, db.GetLocationIDByAnyBrokerCodeParams{
			Broker: db.BrokerHertzGr, BrokerLocationID: clashingDesk,
		})
		if !errors.Is(err, db.ErrNoRows) {
			t.Errorf("err = %v, want the clashing desk left out", err)
		}
	})

	t.Run("a re-import refreshes a known desk's brands", func(t *testing.T) {
		importStations(t, station(hertzDesk, "FF", "ZE"))
		if got := codeOf(t, hertzDesk).Brands; !slices.Equal(got, []string{"FF", "ZE"}) {
			t.Errorf("Brands = %v, want [FF ZE]", got)
		}
	})

	t.Run("a re-import keeps brands that would clash with another desk", func(t *testing.T) {
		importStations(t, station(hertzDesk, "ZE", "ZR"))
		if got := codeOf(t, hertzDesk).Brands; !slices.Equal(got, []string{"FF", "ZE"}) {
			t.Errorf("Brands = %v, want the previous [FF ZE]", got)
		}
	})

	t.Run("a broker without brands still holds one station per location", func(t *testing.T) {
		_, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
			LocationID: airport.ID, Broker: db.BrokerFlex, BrokerLocationID: "flex-hgr-second-" + suffix,
		})
		if !db.IsUniqueViolation(err) {
			t.Errorf("err = %v, want a unique violation", err)
		}
	})
}

func TestInsertLocationBrands(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()

	valid := location.InsertLocationParams{
		Broker: broker.BrokerHertzGR, ID: "hgr-manual", Name: "Hertz GR Manual",
		Country: "Greece", CountryCode: "GR", Brands: []string{"ZE", "ZT"},
	}

	t.Run("validates brands against the broker", func(t *testing.T) {
		cases := []struct {
			name   string
			broker broker.Name
			brands []string
			valid  bool
		}{
			{"hertz gr with brands", broker.BrokerHertzGR, []string{"ZE", "ZT"}, true},
			{"hertz gr without brands", broker.BrokerHertzGR, nil, false},
			{"hertz gr with an unknown brand", broker.BrokerHertzGR, []string{"ZE", "XX"}, false},
			{"hertz gr with a repeated brand", broker.BrokerHertzGR, []string{"ZE", "ZE"}, false},
			{"flex without brands", broker.BrokerFlex, nil, true},
			{"flex with brands", broker.BrokerFlex, []string{"ZE"}, false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				p := valid
				p.Broker = tc.broker
				p.Brands = tc.brands
				err := p.Validate()
				if tc.valid {
					if err != nil {
						t.Errorf("expected no error, got %v", err)
					}
					return
				}
				api_errors.AssertApiError(t, invalidValueErr("brands"), err)
			})
		}
	})

	t.Run("stores the brands of a manually added station", func(t *testing.T) {
		p := valid
		p.ID = fmt.Sprintf("hgr-manual-%d-%d", time.Now().UnixNano(), nextLocSeq())
		p.Name = "Hertz GR Manual " + p.ID
		if err := InsertLocation(ctx, p); err != nil {
			t.Fatalf("InsertLocation: %v", err)
		}

		c, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerHertzGr, BrokerLocationID: p.ID, LocationID: -1,
		})
		if err != nil {
			t.Fatalf("GetLocationBrokerCode: %v", err)
		}
		t.Cleanup(func() {
			_, _ = q.DeleteLocationBrokerCode(ctx, c.ID)
			_ = q.DeleteLocationByID(ctx, c.LocationID)
		})

		if !slices.Equal(c.Brands, []string{"ZE", "ZT"}) {
			t.Errorf("Brands = %v, want [ZE ZT]", c.Brands)
		}
	})
}
