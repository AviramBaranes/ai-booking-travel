package availability

import "encore.app/services/booking/db"

type AvailabilityService struct {
	query db.Querier
	tc    TranslationCacher
	cfg   *AvailableVehiclesConfig
}

func NewAvailabilityService(query db.Querier, tc TranslationCacher, cfg *AvailableVehiclesConfig) *AvailabilityService {
	return &AvailabilityService{query: query, tc: tc, cfg: cfg}
}

type TranslationCacher interface {
	// GetVerified retrieves the translated text for a given source text, if it exists in the verified cache.
	GetVerified(sourceText string) (string, bool)

	// Exists reports whether the source text exists in the database, regardless of translation status.
	Exists(sourceText string) bool

	// NormalizeSentence replaces numeric tokens with indexed placeholders like {num1}.
	NormalizeSentence(sentence string) (string, map[string]string)

	// InsertValuesToSentence inserts values back into a translated template.
	InsertValuesToSentence(template string, values map[string]string) string
}
