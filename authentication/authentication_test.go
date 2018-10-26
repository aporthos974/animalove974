package authentication

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reunion/configuration"
	"strings"
	"testing"
)

func TestLoginSuccessWithEncryptedPassword(test *testing.T) {
	initCollections()
	insertUser("toto@toto", "totoPassword")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication", strings.NewReader("{\"Email\": \"toto@toto\", \"Password\": \"totoPassword\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusOK, writer.Code)
	assert.NotEmpty(test, writer.Header().Get("Authorization"))
}

func TestLoginIsMethodNotAllowedWhenItIsNotPostMethod(test *testing.T) {
	initCollections()
	insertUser("toto@toto", "totoPassword")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "localhost:8080/authentication", strings.NewReader("{\"Email\": \"toto@toto\", \"Password\": \"totoPassword\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusMethodNotAllowed, writer.Code)
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func TestLoginFailedOnIncorrectUsername(test *testing.T) {
	initCollections()
	insertUser("toto@toto", "tutu")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication", strings.NewReader("{\"Email\": \"otherUser\", \"Password\": \"tutu\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusUnauthorized, writer.Code)
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func TestLoginFailedOnIncorrectPassword(test *testing.T) {
	initCollections()
	insertUser("toto@toto", "tutu")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost:8080/authentication", strings.NewReader("{\"Email\": \"toto@toto\", \"Password\": \"incorrectPassword\"}"))

	Login(writer, request)

	assert.Equal(test, http.StatusUnauthorized, writer.Code)
	assert.Empty(test, writer.Header().Get("Authorization"))
}

func insertUser(username string, password string) {
	collections := configuration.GetUserCollection()
	collections.Insert(User{Email: "toto@toto", Password: EncryptSHA256(password)})
}

func initCollections() {
	configuration.Conf = new(configuration.TestConfiguration)
	collections := configuration.GetUserCollection()
	collections.DropCollection()
}
