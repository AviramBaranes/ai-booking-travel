package booking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"encore.app/services/booking/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestGetPriceOfferByIdPicksTheBrandsStation covers a Hertz GR offer at a location holding a
// station per brand: the offer must resolve to its own brand's desk.
func TestGetPriceOfferByIdPicksTheBrandsStation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	pool := sqldb.Driver[*pgxpool.Pool](bookingsDB)
	suffix := fmt.Sprintf("%d-%d", time.Now().UnixNano(), nextLocSeq())

	loc, hertzDesk := seedHertzGRStation(t, q,
		db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Offer Brands Airport " + suffix},
		"hgr-offer-ze-"+suffix, []string{"ZE"},
	)
	thriftyDesk, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
		LocationID: loc.ID, Broker: db.BrokerHertzGr, BrokerLocationID: "hgr-offer-zr-zt-" + suffix, Brands: []string{"ZR", "ZT"},
	})
	if err != nil {
		t.Fatalf("seed Thrifty desk: %v", err)
	}
	t.Cleanup(func() { _, _ = q.DeleteLocationBrokerCode(ctx, thriftyDesk.ID) })

	seedOffer := func(t *testing.T, supplierCode string) int64 {
		t.Helper()
		var id int64
		err := pool.QueryRow(ctx, `
			INSERT INTO price_offers (
				agent_id, name, pickup_location_id, dropoff_location_id, pickup_date, dropoff_date,
				pickup_time, dropoff_time, driver_age, rental_days, plan_id, broker, rate_qualifier,
				supplier_code, car_details, currency_code, currency_rate, purchase_price,
				markup_percentage, broker_erp_price, bt_erp_price, total_price, offered_currency_code,
				offered_price
			) VALUES (
				1, 'Brand Offer', $1, $1, '2030-01-01', '2030-01-05', '10:00', '10:00', '30', 4,
				'1', 'hertz_gr', 'ITRT296:A', $2, '{}', 'EUR', 1, 100, 10, 0, 0, 110, 'EUR', 110
			) RETURNING id`, loc.ID, supplierCode).Scan(&id)
		if err != nil {
			t.Fatalf("seed offer: %v", err)
		}
		t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM price_offers WHERE id = $1`, id) })
		return id
	}

	for brand, want := range map[string]string{"ZT": thriftyDesk.BrokerLocationID, "ZE": hertzDesk.BrokerLocationID} {
		t.Run(brand, func(t *testing.T) {
			offer, err := q.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{ID: seedOffer(t, brand), AgentID: 1})
			if err != nil {
				t.Fatalf("GetPriceOfferById: %v", err)
			}
			if offer.PickupBrokerLocationID != want || offer.DropoffBrokerLocationID != want {
				t.Errorf("stations = %s / %s, want %s", offer.PickupBrokerLocationID, offer.DropoffBrokerLocationID, want)
			}
		})
	}
}
