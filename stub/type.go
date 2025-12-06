package stub

import (
	"crypto/rsa"
	"time"
)

// clock

type (
	TimeNowFunc func() time.Time
)

// cryptox

type (
	VerifyHashPasswordFunc func(hashedPassword, plainPassword string) error
	GetRSAPrivateKeyFunc   func(privateKey string) (*rsa.PrivateKey, error)
	GetRSAPublicKeyFunc    func(publicKey string) (*rsa.PublicKey, error)
)

// uuid

type (
	UuidFunc = func() string
)

// str

type (
	GenerateRandomStringFunc func(length int) string
)
