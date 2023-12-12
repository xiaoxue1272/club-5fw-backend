package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"strings"
)

func GenerateRsaKey(keySize int) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		logger.Panicf("Generate RSA key failed.\nError: %v", err)
	}
	return key
}

func ResolveRsaPrivateKey(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
		return nil, errors.New("PEM file type is not private key")
	}
	var privateKey any
	var err error
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.New("Rsa private key is not pkc1 or pkc8 type")
		}
	}
	return privateKey.(*rsa.PrivateKey), nil
}

func ResolveRsaPublicKey(bytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil || !strings.Contains(block.Type, "PUBLIC KEY") {
		return nil, errors.New("PEM file type is not public key")
	}
	var publicKey any
	var err error
	publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, errors.New("Rsa public key is not pkc1 or pkc8(pkix) type")
		}
	}
	return publicKey.(*rsa.PublicKey), nil
}

func RsaPrivateKeyToString(key *rsa.PrivateKey) (string, error) {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	writer := &strings.Builder{}
	err := pem.Encode(writer, block)
	return writer.String(), err
}

func RsaPublicKeyToString(key *rsa.PublicKey) (string, error) {
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(key),
	}
	writer := &strings.Builder{}
	err := pem.Encode(writer, block)
	return writer.String(), err
}
