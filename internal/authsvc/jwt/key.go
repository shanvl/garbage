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

// KeysFromFiles loads private and public keys
func KeysFromFiles(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	prKey, err := PrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := PublicKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}
	return prKey, pubKey, nil
}
