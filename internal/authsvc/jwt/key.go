package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

// PrivateKeyFromFile reads private rsa key from the given file
func PrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read key file: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("couldn't get private key: %w", err)
	}
	return key, nil
}

// PrivateKeyFromFile reads public rsa key from the given file
func PublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read key file: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("couldn't get public key: %w", err)
	}
	return key, nil
}
