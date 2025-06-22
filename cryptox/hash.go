package cryptox

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashValue(value string) (strBase64, strHash string) {
	// generate refresh token
	strBase64 = base64.URLEncoding.EncodeToString([]byte(value))

	// hash refresh token
	hashedBase64Value := sha256.Sum256([]byte(strBase64))
	strHash = base64.URLEncoding.EncodeToString(hashedBase64Value[:])

	return
}

func IsHashValid(storedHash, value string) bool {
	sum := sha256.Sum256([]byte(value))
	return base64.URLEncoding.EncodeToString(sum[:]) == storedHash
}
