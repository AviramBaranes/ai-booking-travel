package lang

import "context"

type contextKey string

const ContextKey contextKey = "lang"

// FromContext returns the language stored in ctx or def when unavailable.
func FromContext(ctx context.Context, def string) string {
	if v, ok := ctx.Value(ContextKey).(string); ok && v != "" {
		return v
	}
	return def
}
