package location

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrInvalidContentType = api_errors.NewValidationError("invalid content type: expected multipart/form-data")
	ErrParseMultipartForm = api_errors.NewValidationError("failed to parse multipart form")
	ErrGetFileFromForm    = api_errors.NewValidationError("failed to get file from form data")
)

// InsertFlexLocations fetches all Flex locations from the broker and upserts them into the database.
func (s *LocationService) InsertFlexLocations(ctx context.Context) error {
	flex := broker.NewFlex()
	err := s.InsertLocations(ctx, flex)
	if err != nil {
		rlog.Error("failed to insert Flex locations", "error", err)
		return api_errors.ErrInternalError
	}
	return nil
}

// InsertHertzLocations reads a CSV file from the multipart request and upserts Hertz locations.
func (s *LocationService) InsertHertzLocations(w http.ResponseWriter, req *http.Request) {
	file, err := ExtractFile(req)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}
	defer file.Close()

	hertz := broker.NewHertzWithReader(file)

	ctx := req.Context()
	err = s.InsertLocations(ctx, hertz)
	if err != nil {
		rlog.Error("failed to insert locations", "error", err)
		http.Error(w, "failed to insert locations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// InsertLocations fetches locations from the given broker and inserts them into the database.
func (s *LocationService) InsertLocations(ctx context.Context, b broker.LocationSearcher) error {
	var cursor string
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

// ExtractFile parses a multipart/form-data request and returns the uploaded file.
func ExtractFile(req *http.Request) (multipart.File, error) {
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return nil, ErrInvalidContentType
	}

	err := req.ParseMultipartForm(2 << 20) // 2 MB max memory
	if err != nil {
		return nil, ErrParseMultipartForm
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		return nil, ErrGetFileFromForm
	}

	return file, nil
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
		iata := normalizeIata(loc.Iata)

		var locationID int64

		if iata != "" {
			locationID, err = q.UpsertLocationByIATA(ctx, db.UpsertLocationByIATAParams{
				Country:     loc.Country,
				CountryCode: loc.CountryCode,
				City:        loc.City,
				Name:        loc.Name,
				Iata:        iata,
			})
		} else {
			locationID, err = q.UpsertLocationByCountryCodeName(ctx, db.UpsertLocationByCountryCodeNameParams{
				Country:     loc.Country,
				CountryCode: loc.CountryCode,
				City:        loc.City,
				Name:        loc.Name,
			})
		}
		if err != nil {
			return fmt.Errorf("failed to insert for supplier %s, locationId %s: %w", brokerName, loc.ID, err)
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

func toDbBroker(sn broker.Name) (db.Broker, error) {
	switch sn {
	case broker.BrokerFlex:
		return db.BrokerFlex, nil
	case broker.BrokerHertz:
		return db.BrokerHertz, nil
	default:
		return "", fmt.Errorf("unknown broker: %s", sn)
	}
}
