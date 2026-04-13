package booking

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"encore.dev/rlog"
)

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

type numberMatch struct {
	value string
	start int
	end   int
}

var numberRegex = regexp.MustCompile(`\d+(?:\.\d+)?`)
var placeholderRegex = regexp.MustCompile(`\{num\d+\}`)

// normalizeSentence replaces numeric tokens with indexed placeholders like {num1}.
func normalizeSentence(sentence string) (string, map[string]string) {
	indexes := numberRegex.FindAllStringIndex(sentence, -1)
	if len(indexes) == 0 {
		return sentence, nil
	}

	values := make(map[string]string, len(indexes))
	var b strings.Builder

	last := 0
	for i, idx := range indexes {
		start, end := idx[0], idx[1]
		key := fmt.Sprintf("num%d", i+1)
		placeholder := "{" + key + "}"

		values[key] = sentence[start:end]

		b.WriteString(sentence[last:start])
		b.WriteString(placeholder)
		last = end
	}

	b.WriteString(sentence[last:])
	return b.String(), values
}

// insertValuesToSentence inserts values back into a translated template.
func insertValuesToSentence(template string, values map[string]string) string {
	indexes := placeholderRegex.FindAllStringIndex(template, -1)
	if len(indexes) == 0 || len(values) == 0 {
		return template
	}

	var b strings.Builder
	last := 0

	for _, idx := range indexes {
		start, end := idx[0], idx[1]
		placeholder := template[start:end]
		key := strings.TrimSuffix(strings.TrimPrefix(placeholder, "{"), "}")

		b.WriteString(template[last:start])

		if value, ok := values[key]; ok {
			b.WriteString(value)
		} else {
			b.WriteString(placeholder)
		}

		last = end
	}

	b.WriteString(template[last:])
	return b.String()
}
