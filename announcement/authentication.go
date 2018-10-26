package announcement

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"reunion/configuration"
	"reunion/token"
)

type Account struct {
	Email, Password string `xml:"-"`
}

func Login(writer http.ResponseWriter, request *http.Request) {
	announcementId := request.FormValue("id")
	if request.Method != "POST" || announcementId == "" {
		doResponse(writer, http.StatusMethodNotAllowed, "Opération non autorisée ou paramètre manquant: %s", request.Method+" id = "+announcementId)
		return
	}
	var account Account
	jsonErr := json.NewDecoder(request.Body).Decode(&account)
	if jsonErr != nil {
		doResponse(writer, http.StatusBadRequest, "Erreur lors de la désérialisation pour l'authentification : %s", jsonErr.Error())
		return
	}
	defineAuthenticated(writer, account, announcementId)
}

func EncryptSHA256(password string) string {
	salt, err := ioutil.ReadFile(configuration.Conf.GetFilePath(configuration.GetConfiguration().GetKeys().Private))
	if err != nil {
		log.Panicf("Erreur lors de la récupération du salt : %s", err)
	}

	return base64.StdEncoding.EncodeToString(computeHasherWithSalt(password, salt).Sum([]byte(password)))
}

func doResponse(writer http.ResponseWriter, httpCode int, logMessage string, logComplement string) {
	writer.WriteHeader(httpCode)
	writer.Write([]byte(fmt.Sprintf(logMessage, logComplement)))
	log.Printf(logMessage, logComplement)
}

func computeHasherWithSalt(data string, salt []byte) (hasher hash.Hash) {
	hasher = sha256.New()
	hasher.Write([]byte(data))
	hasher.Write(salt)
	return hasher
}

func isCorrectUserAndPasswordForAnnouncement(account Account, announcementId string) bool {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	number, err := collection.Find(bson.M{"_id": bson.ObjectIdHex(announcementId), "account.email": account.Email, "account.password": EncryptSHA256(account.Password)}).Count()
	if err != nil {
		log.Panicf("Erreur lors vérification du username et password en bdd : %s", err)
	}
	return number == 1
}

func CheckTokenOnAnnouncement(announcementId bson.ObjectId, tokenToVerify string) bool {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	number, err := collection.Find(bson.M{"_id": announcementId, "account.token": tokenToVerify}).Count()
	if err != nil {
		log.Panicf("Erreur lors vérification du token en bdd : %s", err.Error())
	}
	return number == 1
}

func defineAuthenticated(writer http.ResponseWriter, account Account, announcementId string) {
	if isCorrectUserAndPasswordForAnnouncement(account, announcementId) {
		generatedToken := token.Create("Authorization")
		collection, session := configuration.GetAnnouncementCollection()
		defer session.Close()

		if err := collection.UpdateId(bson.ObjectIdHex(announcementId), bson.M{"$set": bson.M{"account.token": generatedToken}}); err != nil {
			log.Panicf("Erreur lors la modification du token en bdd : %s", err.Error())
		}

		writer.Header().Set("Authorization", generatedToken)
		doResponse(writer, http.StatusOK, "Authentification réussie pour : %s", account.Email)
	} else {
		doResponse(writer, http.StatusUnauthorized, "Erreur lors de l'authentification pour le user : %s", account.Email)
	}
}
