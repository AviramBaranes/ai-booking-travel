package middleware

import (
	"encore.app/internal/api_errors"
	"encore.dev/middleware"
)

var secrets struct {
	serviceClientToken string
}

// RequireServiceClient is a middleware that restricts access to EP only to Service clients via a token.
// encore:middleware global target=tag:service_client
func RequireServiceClient(req middleware.Request, next middleware.Next) middleware.Response {
	token := req.Data().Headers.Get("X-Service-Client-Token")
	if token != secrets.serviceClientToken {
		return middleware.Response{
			Err: api_errors.ErrNotFound, // Return 404 to avoid revealing the existence of the endpoint
		}
	}

	return next(req)
}
