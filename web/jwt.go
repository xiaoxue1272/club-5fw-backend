package web

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"github.com/xiaoxue1272/club-5fw-backend/utils"
	"time"
)

type ClubGlobalClaims struct {
	*jwt.RegisteredClaims
	Data any `json:"data,omitempty"`
}

const Issuer = "5fw.club"

var PrivateKey *rsa.PrivateKey

var PublicKey *rsa.PublicKey

var SignMethod jwt.SigningMethod

func initJwt(jwtConfig *config.JwtConfiguration) {
	SignMethod = jwt.GetSigningMethod(jwtConfig.Alg)
	if SignMethod == nil {
		logger.Panic("Unknown jwt alg method, please check it")
	}

	if jwtConfig.Rsa.PrivateKey != "" && jwtConfig.Rsa.PublicKey != "" {
		logger.Info("Loading jwt rsa key pair from jwt configuration")
		var err error
		PrivateKey, err = utils.ResolveRsaPrivateKey([]byte(jwtConfig.Rsa.PrivateKey))
		if err != nil {
			logger.Panicf("Failed to load jwt rsa private key from pem file, casue %v", err)
		}
		PublicKey, err = utils.ResolveRsaPublicKey([]byte(jwtConfig.Rsa.PublicKey))
		if err != nil {
			logger.Panicf("Failed to load jwt rsa public key from pem file, casue %v", err)
		}
	} else {
		logger.Panic("Jwt RSA key pair pem is empty")
	}
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
	return jwt.NewWithClaims(SignMethod, claims).SignedString(PrivateKey)
}

func resolveJwt(tokenString string) (*any, error) {
	token, err := jwtParse.ParseWithClaims(tokenString, &ClubGlobalClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method == SignMethod {
			return PublicKey, nil
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
