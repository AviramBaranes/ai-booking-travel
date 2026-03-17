package booking

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

// InsertFlexLocations handles the HTTP request to insert Flex locations.
// encore:api auth method=POST path=/locations/flex tag:admin
func (s *Service) InsertFlexLocations(ctx context.Context) error {
	flex := broker.NewFlex()
	err := insertLocations(ctx, &flex, s.query)
	if err != nil {
		rlog.Error("failed to insert Flex locations", "error", err)
		return api_errors.ErrInternalError
	}
	return nil
}

// InsertHertzLocations handles the HTTP request to insert Hertz locations.
// encore:api auth method=POST path=/locations/hertz tag:admin raw
func (s *Service) InsertHertzLocations(w http.ResponseWriter, req *http.Request) {
	file, err := extractFile(req)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}
	defer file.Close()

	hertz := broker.NewHertzWithReader(file)

	ctx := req.Context()
	err = insertLocations(ctx, hertz, s.query)
	if err != nil {
		rlog.Error("failed to insert locations", "error", err)
		http.Error(w, "failed to insert locations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var (
	errInvalidContentType = api_errors.NewValidationError("invalid content type: expected multipart/form-data")
	errParseMultipartForm = api_errors.NewValidationError("failed to parse multipart form")
	errGetFileFromForm    = api_errors.NewValidationError("failed to get file from form data")
)

func extractFile(req *http.Request) (multipart.File, error) {
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return nil, errInvalidContentType
	}

	err := req.ParseMultipartForm(2 << 20) // 2 MB max memory
	if err != nil {
		return nil, errParseMultipartForm
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		return nil, errGetFileFromForm
	}

	return file, nil
}

// insertLocations fetches locations from the given broker and inserts them into the database using the provided querier.
func insertLocations(ctx context.Context, b broker.Broker, q db.Querier) error {
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
			err = insertBatch(ctx, q, page.Locations, b.Name())
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

// normalizeIata normalizes an IATA code by trimming whitespace and converting to uppercase. If the resulting string is not exactly 3 characters long, it returns an empty string.
func normalizeIata(iata string) string {
	s := strings.TrimSpace(strings.ToUpper(iata))
	if len(s) != 3 {
		return ""
	}
	return s
}

// insertBatch inserts a batch of locations into the database, using the IATA code for upsert if available, and associating each location with the broker's location ID. It returns an error if any database operation fails.
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

// toDbBroker converts a broker.Name to a db.Broker, returning an error if the broker is unknown.
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
