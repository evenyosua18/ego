package cryptox

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/evenyosua18/ego/code"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", code.Wrap(err, code.EncryptionError)
	}

	return string(bytes), nil
}

func VerifyHashedPassword(hashedPassword, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return code.Wrap(err, code.EncryptionError).SetMessage("password is invalid")
	}

	return nil
}
