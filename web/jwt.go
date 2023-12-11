package web

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"time"
)

type ClubGlobalClaims struct {
	*jwt.RegisteredClaims
	Data any `json:"data,omitempty"`
}

const Issuer = "5fw.club"

var rsaKey *rsa.PrivateKey

var SignMethod = jwt.SigningMethodRS512

func initJwt() {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logger.Panicf("Generate RSA key failed.\nError: %v", err)
	}
	rsaKey = key
}

var jwtParse = jwt.NewParser(
	jwt.WithExpirationRequired(),
	jwt.WithIssuer(Issuer),
	jwt.WithLeeway(5*time.Second),
	jwt.WithPaddingAllowed(),
	jwt.WithValidMethods([]string{"RS512"}))

func generateJwt(data any) (string, error) {
	claims := &ClubGlobalClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   "Auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
		Data: data,
	}
	return jwt.NewWithClaims(SignMethod, claims).SignedString(rsaKey)
}

func resolveJwt(tokenString string) (*any, error) {
	token, err := jwtParse.ParseWithClaims(tokenString, &ClubGlobalClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method == SignMethod {
			return &rsaKey.PublicKey, nil
		}
		return nil, errors.New("Only RS256 alg can be parsed")
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*ClubGlobalClaims)
	if ok {
		return &claims.Data, nil
	}
	return &claims.Data, errors.New("Unknown Claims Type")
}
