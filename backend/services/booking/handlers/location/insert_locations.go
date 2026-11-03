package location

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	fileupload "encore.app/internal/file_upload"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// InsertFlexLocationsByCountryCode fetches all Flex locations in a specific country and upserts them into the database.
func (s *LocationService) InsertFlexLocationsByCountryCode(ctx context.Context, countryCode string) error {
	flex := broker.NewFlex()
	page, err := flex.GetLocationsPage(countryCode)
	if err != nil {
		rlog.Error("failed to insert Flex locations", "error", err)
		return api_errors.ErrInternalError
	}

	if len(page.Locations) > 0 {
		err = insertBatch(ctx, s.query, page.Locations, flex.Name())
		if err != nil {
			rlog.Error("failed to insert Flex locations batch", "error", err)
			return api_errors.ErrInternalError
		}
	}
	return nil
}

type InsertFlexLocationsParams struct {
	CountryCode string `json:"countryCode"`
}

// InsertFlexLocations fetches all Flex locations from the broker and upserts them into the database.
func (s *LocationService) InsertFlexLocations(ctx context.Context, p InsertFlexLocationsParams) error {
	flex := broker.NewFlex()
	err := s.InsertLocations(ctx, flex, p.CountryCode)
	if err != nil {
		rlog.Error("failed to insert Flex locations", "error", err)
		return api_errors.ErrInternalError
	}
	return nil
}

// InsertHertzLocations reads a CSV file from the multipart request and upserts Hertz locations.
func (s *LocationService) InsertHertzLocations(w http.ResponseWriter, req *http.Request) {
	file, err := fileupload.ExtractFile(req)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}
	defer file.Close()

	hertz := broker.NewHertzWithReader(file)

	ctx := req.Context()
	err = s.InsertLocations(ctx, hertz, "")
	if err != nil {
		rlog.Error("failed to insert locations", "error", err)
		http.Error(w, "failed to insert locations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// InsertAvanceLocations reads an xlsx workbook from the multipart request and upserts Avance locations.
func (s *LocationService) InsertAvanceLocations(w http.ResponseWriter, req *http.Request) {
	file, err := fileupload.ExtractFile(req)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}
	defer file.Close()

	avance := broker.NewAvanceWithReader(file)

	ctx := req.Context()
	err = s.InsertLocations(ctx, avance, "")
	if err != nil {
		rlog.Error("failed to insert locations", "broker", avance.Name(), "error", err)
		http.Error(w, "failed to insert locations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// InsertHertzGRLocations fetches the Autohellas station feed and upserts Hertz GR locations.
func (s *LocationService) InsertHertzGRLocations(ctx context.Context) error {
	hertzGR := broker.NewHertzGR(broker.HertzGRConfig{})
	err := s.InsertLocations(ctx, hertzGR, "")
	if err != nil {
		rlog.Error("failed to insert locations", "broker", hertzGR.Name(), "error", err)
		return api_errors.ErrInternalError
	}
	return nil
}

// InsertLocations fetches locations from the given broker and inserts them into the database.
func (s *LocationService) InsertLocations(ctx context.Context, b broker.LocationSearcher, cursor string) error {
	skippedCursors := make([]string, 0)
	for {
		page, err := b.GetLocationsPage(cursor)
		if err != nil {
			if cursor == "" && page.NextPage == "" {
				rlog.Error("failed to get first page of locations", "broker", b.Name(), "error", err)
				return api_errors.ErrInternalError
			}

			rlog.Error("failed to get locations page, skipping to next page", "broker", b.Name(), "cursor", cursor, "error", err)
			skippedCursors = append(skippedCursors, cursor)
			cursor = page.NextPage
			if cursor == "" {
				rlog.Info("no more pages to skip to after error", "broker", b.Name())
				break
			}
			continue
		}

		if len(page.Locations) > 0 {
			err = insertBatch(ctx, s.query, page.Locations, b.Name())
			if err != nil {
				rlog.Error("failed to insert locations batch", "broker", b.Name(), "cursor", cursor, "error", err)
				return api_errors.ErrInternalError
			}
		}

		if page.NextPage == "" {
			rlog.Info("no more pages to fetch", "broker", b.Name())
			break
		}
		next := page.NextPage
		if next == cursor {
			rlog.Error("broker returned same next page cursor, stopping to avoid infinite loop",
				"broker", b.Name(), "cursor", cursor)
			return api_errors.ErrInternalError
		}
		cursor = next
	}

	if len(skippedCursors) > 0 {
		rlog.Warn("skipped some pages of locations due to errors", "broker", b.Name(), "skipped_cursors", skippedCursors)
	}

	return nil
}

func normalizeIata(iata string) string {
	s := strings.TrimSpace(strings.ToUpper(iata))
	if len(s) != 3 {
		return ""
	}
	return s
}

func insertBatch(ctx context.Context, q db.Querier, locs []broker.Location, brokerName broker.Name) error {
	dbBroker, err := toDbBroker(brokerName)
	if err != nil {
		return err
	}

	for _, loc := range locs {
		// A known code keeps the location it points to, so a re-import cannot undo an admin merge
		// by recreating the station under its supplier name.
		knownLocationID, err := q.GetLocationIDByAnyBrokerCode(ctx, db.GetLocationIDByAnyBrokerCodeParams{
			Broker:           dbBroker,
			BrokerLocationID: loc.ID,
		})
		if err == nil {
			if len(loc.Brands) > 0 {
				if err := refreshBrands(ctx, q, knownLocationID, dbBroker, loc); err != nil {
					return err
				}
			}
			continue
		}
		if !errors.Is(err, db.ErrNoRows) {
			return fmt.Errorf("get location by broker code: %w", err)
		}

		locationID, err := resolveLocation(ctx, q, loc)
		if err != nil {
			return fmt.Errorf("failed to insert for supplier %s, locationId %s: %w", brokerName, loc.ID, err)
		}

		if len(loc.Brands) > 0 {
			if err := insertBrandedStation(ctx, q, locationID, dbBroker, loc); err != nil {
				return err
			}
			continue
		}

		_, err = q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker:           dbBroker,
			BrokerLocationID: loc.ID,
			LocationID:       locationID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				if _, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
					LocationID:       locationID,
					Broker:           dbBroker,
					BrokerLocationID: loc.ID,
				}); err != nil {
					return fmt.Errorf("insert location broker code: %w", err)
				}
			} else {
				return fmt.Errorf("get location broker code: %w", err)
			}
		}
	}

	return nil
}

// insertBrandedStation attaches a station of a broker with brands to its location. The location
// may already hold another station of the broker, as long as it serves other brands.
func insertBrandedStation(ctx context.Context, q db.Querier, locationID int64, dbBroker db.Broker, loc broker.Location) error {
	codes, err := q.GetAllLocationBrokerCodesByLocationIDs(ctx, []int64{locationID})
	if err != nil {
		return fmt.Errorf("get location broker codes: %w", err)
	}
	if other := stationConflict(codes, dbBroker, loc.ID, loc.Brands); other != "" {
		rlog.Warn("skipping station whose brands another station at the location already serves",
			"broker", dbBroker, "code", loc.ID, "brands", loc.Brands, "location_id", locationID, "other_code", other)
		return nil
	}

	if _, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
		LocationID:       locationID,
		Broker:           dbBroker,
		BrokerLocationID: loc.ID,
		Brands:           loc.Brands,
	}); err != nil {
		return fmt.Errorf("insert location broker code: %w", err)
	}

	return nil
}

// refreshBrands updates the brands of a station we already have, since a supplier can add or drop
// a brand at a desk. The station stays where it is.
func refreshBrands(ctx context.Context, q db.Querier, locationID int64, dbBroker db.Broker, loc broker.Location) error {
	codes, err := q.GetAllLocationBrokerCodesByLocationIDs(ctx, []int64{locationID})
	if err != nil {
		return fmt.Errorf("get location broker codes: %w", err)
	}
	if other := stationConflict(codes, dbBroker, loc.ID, loc.Brands); other != "" {
		rlog.Warn("keeping station brands that would clash with another station at the location",
			"broker", dbBroker, "code", loc.ID, "brands", loc.Brands, "location_id", locationID, "other_code", other)
		return nil
	}

	if err := q.UpdateLocationBrokerCodeBrands(ctx, db.UpdateLocationBrokerCodeBrandsParams{
		Brands:           loc.Brands,
		Broker:           dbBroker,
		BrokerLocationID: loc.ID,
	}); err != nil {
		return fmt.Errorf("update location broker code brands: %w", err)
	}

	return nil
}

// resolveLocation finds or creates the canonical location a supplier's station belongs to. Existing
// locations are the source of truth, so a match by IATA keeps its name and only has missing fields
// filled in; a station without one is matched by country and name.
func resolveLocation(ctx context.Context, q db.Querier, loc broker.Location) (int64, error) {
	iata := normalizeIata(loc.Iata)

	if iata != "" {
		id, err := q.FillLocationByIATA(ctx, db.FillLocationByIATAParams{
			Country:     loc.Country,
			CountryCode: loc.CountryCode,
			City:        loc.City,
			Iata:        iata,
		})
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, db.ErrNoRows) {
			return 0, err
		}
	}

	return q.UpsertLocationByCountryCodeName(ctx, db.UpsertLocationByCountryCodeNameParams{
		Country:     loc.Country,
		CountryCode: loc.CountryCode,
		City:        loc.City,
		Name:        loc.Name,
		Iata:        iata,
	})
}

func toDbBroker(sn broker.Name) (db.Broker, error) {
	switch sn {
	case broker.BrokerFlex:
		return db.BrokerFlex, nil
	case broker.BrokerHertz:
		return db.BrokerHertz, nil
	case broker.BrokerAvance:
		return db.BrokerAvance, nil
	case broker.BrokerHertzGR:
		return db.BrokerHertzGr, nil
	default:
		return "", fmt.Errorf("unknown broker: %s", sn)
	}
}
