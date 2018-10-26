package announcement

import (
	"html/template"
	"log"
	"net/http"
	. "reunion/configuration"
)

func GetAdoptFormPage(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.New("base").Delims("[[", "]]").ParseFiles(Conf.GetFilePath("static/html/create_adopt_announcement.html"), Conf.GetFilePath("static/html/base.html")))
	description := "Formulaire de signalement d'un chien ou un chat à adopter. Ce formulaire permet de décrire l'animal à la recherche d'un nouveau propriétaire."
	page := Page{name: "Animal à adopter", Title: "AnimaLove - Signaler un chien ou un chat à adopter à la Réunion", Description: description, Image: "/static/images/logo.png", Link: Conf.GetUrl()}
	writer.Header().Set("Cache-control", "public, max-age=3600")
	err := tmpl.Execute(writer, page)
	if err != nil {
		log.Panicf("Erreur lors de la construction du template : %s", err.Error())
	}
}
