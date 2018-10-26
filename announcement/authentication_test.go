package announcement

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/http/httptest"
	"reunion/configuration"
	"strings"
	"testing"
)

func TestLoginSuccessWithEncryptedPassword(test *testing.T) {
	initCollections()
	insertUserAndAnnouncement("toto@test", "totoPassword")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication?id=54f070908d7d4f1bfdb7a9dd", strings.NewReader("{\"email\": \"toto@test\", \"password\": \"totoPassword\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusOK, writer.Code)
	assert.Equal(test, "Authentification réussie pour : toto@test", writer.Body.String())
	assert.NotEmpty(test, writer.Header().Get("Authorization"))
}

func TestLoginIsMethodNotAllowedWhenItIsNotPostMethod(test *testing.T) {
	initCollections()
	insertUserAndAnnouncement("toto@test", "totoPassword")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "localhost:8080/authentication?id=54f070908d7d4f1bfdb7a9dd", strings.NewReader("{\"email\": \"toto@test\", \"password\": \"totoPassword\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusMethodNotAllowed, writer.Code)
	assert.Equal(test, "Opération non autorisée ou paramètre manquant: PUT id = 54f070908d7d4f1bfdb7a9dd", writer.Body.String())
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func TestLoginFailedOnIncorrectUsername(test *testing.T) {
	initCollections()
	insertUserAndAnnouncement("toto@test", "tutu")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication?id=54f070908d7d4f1bfdb7a9dd", strings.NewReader("{\"email\": \"otherUser\", \"password\": \"tutu\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusUnauthorized, writer.Code)
	assert.Equal(test, "Erreur lors de l'authentification pour le user : otherUser", writer.Body.String())
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func TestLoginFailedOnIncorrectPassword(test *testing.T) {
	initCollections()
	insertUserAndAnnouncement("toto@test", "tutu")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication?id=54f070908d7d4f1bfdb7a9dd", strings.NewReader("{\"email\": \"toto@test\", \"password\": \"incorrectPassword\"}"))

	Login(writer, request)

	assert.Equal(test, "Erreur lors de l'authentification pour le user : toto@test", writer.Body.String())
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func insertUserAndAnnouncement(username string, password string) {
	collections, session := configuration.GetAnnouncementCollection()
	defer session.Close()
	collections.Insert(Announcement{Id: bson.ObjectIdHex("54f070908d7d4f1bfdb7a9dd"), Account: &Account{Email: username, Password: EncryptSHA256(password)}})
}

func initCollections() {
	configuration.Conf = new(configuration.TestConfiguration)
	collections, session := configuration.GetAnnouncementCollection()
	defer session.Close()
	collections.DropCollection()
}
