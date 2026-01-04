package cryptox

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString generates a cryptographically secure random string of length n (in bytes, so 2*n in hex chars)
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
