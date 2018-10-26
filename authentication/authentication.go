package authentication

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"reunion/configuration"
	"reunion/token"
)

type User struct {
	Email    string
	Password string
}

func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user User
	jsonErr := json.NewDecoder(request.Body).Decode(&user)
	if jsonErr != nil {
		log.Printf("Erreur lors de la désérialisation pour l'authentification : %s", jsonErr)
		return
	}
	defineAuthentication(writer, user)
}

func EncryptSHA256(password string) string {
	salt, err := ioutil.ReadFile(configuration.Conf.GetFilePath(configuration.GetConfiguration().GetKeys().Private))
	if err != nil {
		log.Printf("Erreur lors récupération du salt : %s", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(computeHasherWithSalt(password, salt).Sum([]byte(password)))
}

func computeHasherWithSalt(data string, salt []byte) (hasher hash.Hash) {
	hasher = sha256.New()
	hasher.Write([]byte(data))
	hasher.Write(salt)
	return hasher
}

func isCorrectUserAndPassword(user User) bool {
	number, err := configuration.GetUserCollection().Find(bson.M{"email": user.Email, "password": EncryptSHA256(user.Password)}).Count()
	if err != nil {
		log.Printf("Erreur lors vérification du username et password en bdd : %s", err)
		return false
	}
	return number == 1
}

func defineAuthentication(writer http.ResponseWriter, user User) {
	if isCorrectUserAndPassword(user) {
		writer.Header().Set("Authorization", token.Create("Authorization"))
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}
