package locations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/supplier"
	"encore.app/services/locations/db"
	"encore.dev/rlog"
)

// encore:api auth method=POST path=/locations/flex tag:admin
func (s *Service) InsertFlexLocations(ctx context.Context) error {
	flex := supplier.NewFlex()
	return insertLocations(ctx, &flex, s.query)
}

// encore:api auth method=POST path=/locations/hertz tag:admin raw
func (s *Service) InsertHertzLocations(w http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		rlog.Warn("invalid content type for Hertz locations upload", "content_type", ct)
		http.Error(w, "content type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	err := req.ParseMultipartForm(2 << 20) // 2 MB max memory
	if err != nil {
		rlog.Error("failed to parse multipart form", "error", err)
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		rlog.Error("failed to get file from form data", "error", err)
		http.Error(w, "failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	hertz := supplier.NewHertzWithReader(file)

	ctx := req.Context()
	err = insertLocations(ctx, hertz, s.query)
	if err != nil {
		rlog.Error("failed to insert locations", "error", err)
		http.Error(w, "failed to insert locations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func insertLocations(ctx context.Context, sup supplier.Supplier, q db.Querier) error {
	var cursor string
	skippedCursors := make([]string, 0)
	for {
		page, err := sup.GetLocationsPage(cursor)
		if err != nil {
			if cursor == "" && page.NextPage == "" {
				rlog.Error("failed to get first page of locations", "supplier", sup.Supplier(), "error", err)
				return api_errors.ErrInternalError
			}

			rlog.Error("failed to get locations page, skipping to next page", "supplier", sup.Supplier(), "cursor", cursor, "error", err)
			skippedCursors = append(skippedCursors, cursor)
			cursor = page.NextPage
			if cursor == "" {
				rlog.Info("no more pages to skip to after error", "supplier", sup.Supplier())
				break
			}
			continue
		}

		if len(page.Locations) > 0 {
			err = insertBatch(ctx, q, page.Locations, sup.Supplier())
			if err != nil {
				rlog.Error("failed to insert locations batch", "supplier", sup.Supplier(), "cursor", cursor, "error", err)
				return api_errors.ErrInternalError
			}
		}

		if page.NextPage == "" {
			rlog.Info("no more pages to fetch", "supplier", sup.Supplier())
			break
		}
		next := page.NextPage
		if next == cursor {
			rlog.Error("supplier returned same next page cursor, stopping to avoid infinite loop",
				"supplier", sup.Supplier(), "cursor", cursor)
			return api_errors.ErrInternalError
		}
		cursor = next
	}

	if len(skippedCursors) > 0 {
		rlog.Warn("skipped some pages of locations due to errors", "supplier", sup.Supplier(), "skipped_cursors", skippedCursors)
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

func insertBatch(ctx context.Context, q db.Querier, locs []supplier.Location, sup supplier.SupplierName) error {
	dbSup, err := toDBSupplier(sup)
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
			return err
		}

		// Only insert supplier mapping if it does not already exist.
		_, err = q.GetLocationSupplierCode(ctx, db.GetLocationSupplierCodeParams{
			Supplier:           dbSup,
			SupplierLocationID: loc.ID,
			LocationID:         locationID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				if _, err := q.InsertLocationSupplierCode(ctx, db.InsertLocationSupplierCodeParams{
					LocationID:         locationID,
					Supplier:           dbSup,
					SupplierLocationID: loc.ID,
				}); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func toDBSupplier(sn supplier.SupplierName) (db.Supplier, error) {
	switch sn {
	case supplier.SupplierFlex:
		return db.SupplierFlex, nil
	case supplier.SupplierHertz:
		return db.SupplierHertz, nil
	default:
		return "", fmt.Errorf("unknown supplier: %s", sn)
	}
}
