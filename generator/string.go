package generator

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"github.com/evenyosua18/ego/code"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString returns a random string of the given length
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func SecureCode(length int) (string, error) {
	b := make([]byte, length)

	if _, err := cryptoRand.Read(b); err != nil {
		return "", code.Wrap(err, code.EncryptionError)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
