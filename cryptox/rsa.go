package cryptox

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/evenyosua18/ego/code"
)

func GetRSAPrivateKey(strPrivateKey string) (*rsa.PrivateKey, error) {
	// decode the base64 string
	pemData, err := base64.StdEncoding.DecodeString(strPrivateKey)
	if err != nil {
		return nil, code.Wrap(err, code.EncryptionError)
	}

	// decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, code.Get(code.EncryptionError).SetErrorMessage("failed to decode pem block")
	}

	// parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, code.Wrap(err, code.EncryptionError)
	}

	convertedPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, code.Get(code.EncryptionError).SetErrorMessage("not a rsa private key")
	}

	return convertedPrivateKey, nil
}

func GetRSAPublicKey(strPublicKey string) (*rsa.PublicKey, error) {
	pemData, err := base64.StdEncoding.DecodeString(strPublicKey)
	if err != nil {
		return nil, code.Wrap(err, code.EncryptionError)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, code.Get(code.EncryptionError).SetErrorMessage("failed to decode pem block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, code.Wrap(err, code.EncryptionError)
	}

	rsaPub, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, code.Get(code.EncryptionError).SetErrorMessage("not a rsa public key")
	}
	return rsaPub, nil
}
