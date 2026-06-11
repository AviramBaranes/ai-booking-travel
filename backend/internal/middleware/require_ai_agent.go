package middleware

import (
	"encore.app/internal/api_errors"
	"encore.dev/middleware"
)

var secrets struct {
	AIAgentToken string
}

// RequireAIAgent is a middleware that restricts access to EP only to AI agents via a token.
// encore:middleware global target=tag:ai_agent
func RequireAIAgent(req middleware.Request, next middleware.Next) middleware.Response {
	token := req.Data().Headers.Get("X-AI-Agent-Token")
	if token != secrets.AIAgentToken {
		return middleware.Response{
			Err: api_errors.ErrNotFound, // Return 404 to avoid revealing the existence of the endpoint
		}
	}

	return next(req)
}
