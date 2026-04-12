package booking

import (
	"context"
	"sync"
	"time"

	"encore.app/services/booking/db"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
	t     *TranslationCache
}

var bookingsDB = sqldb.NewDatabase("bookings", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

// TranslationCache provides thread-safe access to the in-memory translation map.
type TranslationCache struct {
	mu           sync.RWMutex
	translations map[string]string
	known        map[string]struct{}
}

// GetVerified retrieves the translated text for a given source text, if it exists in the verified cache.
func (c *TranslationCache) GetVerified(sourceText string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	targetText, exists := c.translations[sourceText]
	return targetText, exists
}

// Exists reports whether the source text exists in the database, regardless of translation status.
func (c *TranslationCache) Exists(sourceText string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.known[sourceText]
	return exists
}

// Replace swaps the in-memory verified and known translation sets in a thread-safe manner.
func (c *TranslationCache) Replace(newTranslations map[string]string, knownTranslations map[string]struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.translations = newTranslations
	c.known = knownTranslations
}

// initService initializes the booking service by setting up the database connection and loading translations into the cache.
func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](bookingsDB)
	query = db.New(pgxdb)

	svc := &Service{
		query: query,
		t: &TranslationCache{
			translations: make(map[string]string),
			known:        make(map[string]struct{}),
		},
	}

	if err := svc.refreshTranslations(context.Background()); err != nil {
		return nil, err
	}

	return svc, nil
}

// refresh translations reloads translations from the database and updates the in-memory cache.
func (s *Service) refreshTranslations(ctx context.Context) error {
	translations, err := s.query.GetAllVerifiedTranslations(ctx)
	if err != nil {
		rlog.Error("failed to fetch translations from database", "error", err)
		return err
	}

	knownSourceTexts, err := s.query.GetAllTranslationSourceTexts(ctx)
	if err != nil {
		rlog.Error("failed to fetch known translation source texts from database", "error", err)
		return err
	}

	translationMap := make(map[string]string, len(translations))
	for _, t := range translations {
		translationMap[t.SourceText] = *t.TargetText
	}

	knownTranslations := make(map[string]struct{}, len(knownSourceTexts))
	for _, sourceText := range knownSourceTexts {
		knownTranslations[sourceText] = struct{}{}
	}

	go s.startBackgroundRefresh()

	s.t.Replace(translationMap, knownTranslations)
	return nil
}

// startBackgroundRefresh starts a background ticker that periodically refreshes the translations cache at the specified interval.
func (s *Service) startBackgroundRefresh() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.refreshTranslations(context.Background()); err != nil {
			rlog.Error("failed to refresh translations", "error", err)
		}
	}
}
