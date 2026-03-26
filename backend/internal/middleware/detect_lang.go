package middleware

import (
	"context"

	"encore.dev/middleware"
)

type contextKey string

const LangContextKey contextKey = "lang"

// encore:middleware global target=all
func DetectLangMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	lang := req.Data().Headers.Get("X-Lang")
	if lang == "" {
		lang = "en"
	}

	// set the lang in the context
	ctx := req.Context()
	ctx = context.WithValue(ctx, LangContextKey, lang)
	req = req.WithContext(ctx)
	return next(req)
}
