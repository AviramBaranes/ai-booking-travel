package availability

import (
	"encore.app/internal/api_errors"
	"encore.app/internal/currency"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
)

var (
	errCurrencyRatesUnavailable = api_errors.NewErrorWithDetail(
		errs.Unavailable,
		"Currency rates are unavailable",
		api_errors.ErrorDetails{Code: api_errors.CodeCurrencyRatesUnavailable},
	)

	errAvailabilityPricingFailed = api_errors.NewErrorWithDetail(
		errs.Internal,
		"Failed to price available vehicles",
		api_errors.ErrorDetails{Code: api_errors.CodeAvailabilityPricingFailed},
	)
)

type AvailabilityService struct {
	query db.Querier
	tc    TranslationCacher
	c     *currency.CurrenciesCache
	cfg   *AvailableVehiclesConfig
}

func NewAvailabilityService(query db.Querier, tc TranslationCacher, c *currency.CurrenciesCache, cfg *AvailableVehiclesConfig) *AvailabilityService {
	return &AvailabilityService{query: query, tc: tc, c: c, cfg: cfg}
}

type TranslationCacher interface {
	// GetVerified retrieves the translated text for a given source text, if it exists in the verified cache.
	GetVerified(sourceText string) (string, bool)

	// Exists reports whether the source text exists in the database, regardless of translation status.
	Exists(sourceText string) bool

	// AddKnown adds a source text to the known set, indicating it exists in the database even if not yet translated.
	AddKnown(sourceText string)

	// NormalizeSentence replaces numeric tokens with indexed placeholders like {num1}.
	NormalizeSentence(sentence string) (string, map[string]string)

	// InsertValuesToSentence inserts values back into a translated template.
	InsertValuesToSentence(template string, values map[string]string) string
}
