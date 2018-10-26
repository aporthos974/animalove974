package announcement

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"reunion/configuration"
	"reunion/token"
)

type SecureHandler func(writer http.ResponseWriter, request *http.Request)

func (fn SecureHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	currentToken := request.Header.Get("Authorization")
	if !token.Verify(currentToken) {
		log.Printf("Opération non autorisée : %s", currentToken)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	fn(writer, request)
}

func GetAdminHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		tmpl := template.Must(template.New("base").Delims("[[", "]]").ParseFiles(configuration.Conf.GetFilePath("static/html/admin.html"), configuration.Conf.GetFilePath("static/html/base.html")))
		page := Page{name: "annonce", Title: "Administration"}
		error := tmpl.Execute(writer, page)
		if error != nil {
			log.Panicf("Erreur lors de la construction du template : %s", error.Error())
		}
	} else {
		log.Printf("Opération non autorisée : %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func GetAdminAnnouncementsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		writer.WriteHeader(http.StatusOK)
		writer.Write(GetAnnouncementsTemplate())
	} else {
		log.Printf("Opération non autorisée : %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetAnnouncementsTemplate() []byte {
	content, err := ioutil.ReadFile(configuration.Conf.GetFilePath("static/html/announcements_admin.html"))
	if err != nil {
		log.Panicf("Erreur lors de la récupération du template : %s", err.Error())
	}
	return content
}

func GetAdminLostActionHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "PUT" {
		id, err := getId(request.URL.RequestURI(), 5)
		if state := request.URL.Query().Get("state"); err == nil && state != "" {
			UnrestrictedUpdateState(id, state)
		} else {
			log.Printf("Paramètre manquant : %s", request.URL.RequestURI())
			if err != nil {
				log.Printf("Paramètre manquant : %s", err.Error())
			}
			writer.WriteHeader(http.StatusBadRequest)
		}
	} else {
		log.Printf("Opération non autorisée : %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
