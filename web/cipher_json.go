package web

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"github.com/gin-gonic/gin/binding"
	logger "github.com/sirupsen/logrus"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"github.com/xiaoxue1272/club-5fw-backend/utils"
	"io"
	"net/http"
)

var cipherJsonPrivateKey *rsa.PrivateKey

var cipherJsonPublicKey *rsa.PublicKey

func initJsonCipher(cipherJsonConfig *config.CipherJsonConfiguration) {
	if cipherJsonConfig.Rsa.PrivateKey != "" && cipherJsonConfig.Rsa.PublicKey != "" {
		logger.Info("Loading cipher json rsa key pair from cipher json configuration")
		var err error
		cipherJsonPrivateKey, err = utils.ResolveRsaPrivateKey([]byte(cipherJsonConfig.Rsa.PrivateKey))
		if err != nil {
			logger.Panicf("Failed to load cipher json rsa private key from pem file, casue %v", err)
		}
		cipherJsonPublicKey, err = utils.ResolveRsaPublicKey([]byte(cipherJsonConfig.Rsa.PublicKey))
		if err != nil {
			logger.Panicf("Failed to load cipher json rsa public key from pem file, casue %v", err)
		}
	} else {
		logger.Panic("Cipher json RSA key pair pem is empty")
	}
}

type internalCipherJson struct {
	binding binding.BindingBody
}

func (r *internalCipherJson) Name() string {
	return "cipher-json"
}

func (r *internalCipherJson) Bind(req *http.Request, obj any) error {
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return r.BindBody(bytes, obj)
}

func (r *internalCipherJson) BindBody(body []byte, obj any) error {
	bytes, err := encoder.DecodeString(string(body))
	if err != nil {
		return err
	}
	bytes, err = rsa.DecryptPKCS1v15(rand.Reader, cipherJsonPrivateKey, bytes)
	if err != nil {
		return err
	}
	return r.binding.BindBody(bytes, obj)
}

var CipherJson = &internalCipherJson{
	binding: binding.JSON,
}

var encoder = base64.StdEncoding
