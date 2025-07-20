package cryptox

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

const (
	PKCEMethodS256 = "S256"
)

func IsPKCEValid(codeVerifier, codeChallenge, method string) bool {
	if method == PKCEMethodS256 {
		sum := sha256.Sum256([]byte(codeVerifier))
		expected := base64.RawURLEncoding.EncodeToString(sum[:])
		return subtle.ConstantTimeCompare([]byte(expected), []byte(codeChallenge)) == 1
	}

	return false
}
