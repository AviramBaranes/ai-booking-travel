package auth

import (
	"net/http"
	"time"

	"encore.dev"
)

const (
	refreshCookieName = "refresh_token"
	sessionCookieName = "gr_session"

	// refreshCookieMaxAge is the lifetime of the refresh token cookie (30 days).
	refreshCookieMaxAge = 30 * 24 * time.Hour
)

// isLocal reports whether the app is running in local development.
// Cookie security attributes are relaxed only in this case.
func isLocal() bool {
	return encore.Meta().Environment.Cloud == encore.CloudLocal
}

// cookieDomain returns the parent domain used to share cookies across the
// frontend and API subdomains in dev/prod, and an empty domain locally.
func cookieDomain() string {
	if isLocal() {
		return ""
	}
	return ".aibookingtravel.com"
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
		expiredCookie(refreshCookieName, true).String(),
		expiredCookie(sessionCookieName, false).String(),
	}
}

// newRefreshCookie builds the httpOnly refresh token cookie.
func newRefreshCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/",
		Domain:   cookieDomain(),
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
		Name:     sessionCookieName,
		Value:    "1",
		Path:     "/",
		Domain:   cookieDomain(),
		Secure:   !isLocal(),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(refreshCookieMaxAge.Seconds()),
		Expires:  time.Now().Add(refreshCookieMaxAge),
	}
}

// expiredCookie builds a cookie that immediately expires the named cookie.
func expiredCookie(name string, httpOnly bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   cookieDomain(),
		Secure:   !isLocal(),
		HttpOnly: httpOnly,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
	}
}
