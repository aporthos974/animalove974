package token

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"reunion/configuration"
	"time"
)

func Create(subject string) string {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["exp"] = getExpiration()
	token.Claims["sub"] = subject

	tokenString, err := token.SignedString(getPrivateKey())
	if err != nil {
		log.Panicf("Erreur lors de la signature du token : %s", err)
	}
	return tokenString
}

func Verify(myToken string) bool {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return getPublicKey(), nil
	})

	if err != nil {
		log.Printf("Erreur lors du parsing du token : %s", err)
	}
	return token != nil && token.Valid && token.Claims["sub"] != nil
}

func getExpiration() int64 {
	return time.Now().Add(time.Minute * 2).Unix()
}

func getPrivateKey() []byte {
	privateKey, err := ioutil.ReadFile(configuration.Conf.GetFilePath(configuration.GetConfiguration().GetKeys().Private))
	if err != nil {
		log.Panicf("Erreur dans l'ouverture de la clé privée : %s", err)
	}
	return privateKey
}

func getPublicKey() []byte {
	publicKey, err := ioutil.ReadFile(configuration.Conf.GetFilePath(configuration.GetConfiguration().GetKeys().Public))
	if err != nil {
		log.Panicf("Erreur dans l'ouverture de la clé publique : %s", err)
	}
	return publicKey
}
