package auth

import (
	"net/http"
	"time"

	"encore.dev"
)

const (
	refreshCookieName = "refresh_token"

	prodSessionCookieName = "gr_session"
	devSessionCookieName  = "gr_session_dev"

	prodCookieDomain = "aibookingtravel.com"
	devCookieDomain  = "dev.aibookingtravel.com"

	refreshCookieMaxAge = 30 * 24 * time.Hour
)

// isLocal reports whether the app is running in local development.
func isLocal() bool {
	return encore.Meta().Environment.Cloud == encore.CloudLocal
}

// isProduction reports whether the app is running in production.
func isProduction() bool {
	return encore.Meta().Environment.Type == encore.EnvProduction
}

// sessionCookieName returns a different name for the dev session hint so the
// production cookie, which is visible to all aibookingtravel.com subdomains,
// does not conflict with the dev session hint.
func sessionCookieName() string {
	if !isLocal() && !isProduction() {
		return devSessionCookieName
	}

	return prodSessionCookieName
}

// sessionCookieDomain allows the frontend to read the non-sensitive session
// hint cookie.
func sessionCookieDomain() string {
	if isLocal() {
		return ""
	}

	if isProduction() {
		return prodCookieDomain
	}

	return devCookieDomain
}

// authCookies returns the Set-Cookie header values issued on successful auth:
// the httpOnly refresh token and the JS-readable session hint.
func authCookies(refreshToken string) []string {
	return []string{
		newRefreshCookie(refreshToken).String(),
		newSessionHintCookie().String(),
	}
}

// clearedCookies returns Set-Cookie header values that expire the auth cookies.
func clearedCookies() []string {
	return []string{
		expiredCookie(refreshCookieName, "", true).String(),
		expiredCookie(sessionCookieName(), sessionCookieDomain(), false).String(),
	}
}

// newRefreshCookie builds the httpOnly refresh token cookie.
func newRefreshCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/",
		Secure:   !isLocal(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(refreshCookieMaxAge.Seconds()),
		Expires:  time.Now().Add(refreshCookieMaxAge),
	}
}

// newSessionHintCookie builds the JS-readable cookie used by the frontend to
// detect that a session likely exists before attempting a refresh.
func newSessionHintCookie() *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName(),
		Value:    "1",
		Path:     "/",
		Domain:   sessionCookieDomain(),
		Secure:   !isLocal(),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(refreshCookieMaxAge.Seconds()),
		Expires:  time.Now().Add(refreshCookieMaxAge),
	}
}

// expiredCookie builds a cookie that immediately expires the named cookie.
func expiredCookie(name, domain string, httpOnly bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		Secure:   !isLocal(),
		HttpOnly: httpOnly,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
	}
}
