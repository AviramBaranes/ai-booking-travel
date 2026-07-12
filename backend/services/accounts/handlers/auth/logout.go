package auth

import "context"

// Logout clears the user's auth cookies (refresh token + session hint).
func (s *AuthService) Logout(ctx context.Context) (*LogoutResponse, error) {
	return &LogoutResponse{SetCookies: clearedCookies()}, nil
}
