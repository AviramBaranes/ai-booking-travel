package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateResetToken returns a cryptographically secure random token (in hex)
// and its SHA-256 hash (also hex-encoded), or an error.
func GenerateResetToken() (rawToken string, tokenHash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	rawToken = hex.EncodeToString(b)
	tokenHash = HashToken(rawToken)

	return rawToken, tokenHash, nil
}

// HashToken hashes the provided token using SHA-256 and returns the hex-encoded result.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
