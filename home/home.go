package home

import (
	"html/template"
	"log"
	"net/http"
	. "reunion/announcement"
	. "reunion/configuration"
)

type HomePage struct {
	name, Title, Description, Image, Link string
	Announcements                         []Announcement
}

func GetPage(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.New("base").Delims("[[", "]]").ParseFiles(Conf.GetFilePath("static/html/search.html"), Conf.GetFilePath("static/html/index.html"), Conf.GetFilePath("static/html/base.html")))
	announcements, _ := GetAnnouncementsByType("", 0, 15)
	page := HomePage{Announcements: announcements, name: "index", Title: "AnimaLove - Recherche d'animaux perdus, trouvés et à adopter", Description: "Service de signalement et de recherche d'animaux perdus, errants et à adopter à la Réunion. Chien ou chat à aider.", Image: "/static/images/logo.png", Link: Conf.GetUrl()}
	writer.Header().Set("Cache-control", "public, max-age=3600")
	error := tmpl.Execute(writer, page)
	if error != nil {
		log.Printf("Erreur lors de la construction du template.", error)
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}
}
