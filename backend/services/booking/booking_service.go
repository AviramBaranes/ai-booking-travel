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
}

// Get retrieves the translated text for a given source text, if it exists.
func (c *TranslationCache) Get(sourceText string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	targetText, exists := c.translations[sourceText]
	return targetText, exists
}

// Replace swaps the entire translations map with a new one in a thread-safe manner.
func (c *TranslationCache) Replace(newTranslations map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.translations = newTranslations
}

// initService initializes the booking service by setting up the database connection and loading translations into the cache.
func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](bookingsDB)
	query = db.New(pgxdb)

	svc := &Service{
		query: query,
		t:     &TranslationCache{translations: make(map[string]string)},
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

	translationMap := make(map[string]string, len(translations))
	for _, t := range translations {
		translationMap[t.SourceText] = *t.TargetText
	}

	go s.startBackgroundRefresh()

	s.t.Replace(translationMap)
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
