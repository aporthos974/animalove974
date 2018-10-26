package announcement

import (
	"html/template"
	"log"
	"net/http"
	. "reunion/configuration"
)

func GetSeenFormPage(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.New("base").Delims("[[", "]]").ParseFiles(Conf.GetFilePath("static/html/create_seen_announcement.html"), Conf.GetFilePath("static/html/base.html")))
	description := "Formulaire de signalement d'un chien ou chat errant. Ce formulaire permet de décrire l'animal aperçu et d'indiquer le lieux où il a été vu pour la dernière fois."
	page := Page{name: "Animal errant", Title: "AnimaLove - Signaler un chien ou chat errant à la Réunion", Description: description, Image: "/static/images/logo.png", Link: Conf.GetUrl()}
	writer.Header().Set("Cache-control", "public, max-age=3600")
	err := tmpl.Execute(writer, page)
	if err != nil {
		log.Panicf("Erreur lors de la construction du template : %s", err.Error())
	}
}
