package middleware

import (
	"context"

	"encore.app/internal/lang"
	"encore.dev/middleware"
)

// encore:middleware global target=all
func DetectLangMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	langCode := req.Data().Headers.Get("X-Lang")
	if langCode == "" {
		langCode = "en"
	}

	// set the lang in the context
	ctx := req.Context()
	ctx = context.WithValue(ctx, lang.ContextKey, langCode)
	req = req.WithContext(ctx)
	return next(req)
}
