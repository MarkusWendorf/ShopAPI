package token

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"shopApi/model"
	"time"
)

var (
	publicKey  = mustGetPublicKey()
	privateKey = mustGetPrivateKey()
)

const keyPath = "./keys/"
const tokenExpiration = 1 * time.Hour

func ParseAndVerifyJwt(token string) (bool, model.JwtClaims) {

	claims := model.JwtClaims{}

	tok, err := jwt.ParseWithClaims(token, &claims, func(base64decodedToken *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return false, claims
	}

	return tok.Valid, claims
}

func CreateJWT(user model.User) (string, error) {

	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, model.JwtClaims{
		UserId: user.Id,
		Email:  user.Email,
		Exp:    time.Now().Add(tokenExpiration).Unix(),
	})

	token, err := tk.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func mustGetPrivateKey() *rsa.PrivateKey {

	keyFile, _ := os.Open(keyPath + "private_key.pem")
	keyBytes, _ := ioutil.ReadAll(keyFile)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)

	if err != nil {
		panic("couldn't get private_key.pem")
	}

	return privateKey
}

func mustGetPublicKey() *rsa.PublicKey {

	keyFile, _ := os.Open(keyPath + "public_key.pem")
	keyBytes, _ := ioutil.ReadAll(keyFile)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)

	if err != nil {
		panic("couldn't get public_key.pem")
	}

	return publicKey
}
