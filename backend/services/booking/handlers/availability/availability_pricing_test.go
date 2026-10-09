package availability

import (
	"slices"
	"testing"
)

func TestEnglishIfTranslated(t *testing.T) {
	english := []string{"Unlimited mileage", "Theft Waiver with Excess"}

	if got := englishIfTranslated(english, []string{"קילומטראז' בלתי מוגבל", "Theft Waiver with Excess"}); !slices.Equal(got, english) {
		t.Errorf("translated: got %v, want the English source", got)
	}
	// An English search, or a Hebrew one with nothing translated, stores nothing extra.
	if got := englishIfTranslated(english, english); got != nil {
		t.Errorf("not translated: got %v, want nil", got)
	}
	if got := englishIfTranslated(nil, english); got != nil {
		t.Errorf("no source kept: got %v, want nil", got)
	}
}
