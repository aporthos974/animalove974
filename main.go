package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	. "reunion/announcement"
	"reunion/announcement/rss"
	"reunion/announcement/specy"
	"reunion/authentication"
	. "reunion/compression"
	"reunion/configuration"
	"reunion/home"
	"reunion/minify"
	"reunion/websocket"
)

func main() {
	configuration.NewInstance()

	minify.MinifyJs("static/js/reunionctrl.js")
	minify.MinifyJs("static/js/notify.js")
	minify.MinifyJs("static/js/angular-locale_fr-fr.js")
	minify.MinifyCss("static/css/reunion.css")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(configuration.Conf.GetFilePath("static/")))))
	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("/opt/animalove/photos/"))))
	http.HandleFunc("/static/js/javascript-min.js", GziperHandler(minify.AppJsHandler, "application/javascript"))
	http.HandleFunc("/static/css/styles-min.css", GziperHandler(minify.AppCssHandler, "text/css"))
	http.HandleFunc("/rss.xml", rss.GetFile)
	http.HandleFunc("/robots.txt", configuration.GetRobotsHandler)
	http.HandleFunc("/sitemap.xml", configuration.GetSitemapHandler)

	router := mux.NewRouter()

	router.HandleFunc("/", home.GetPage)
	router.HandleFunc("/animaux", GziperHandler(GetAnnouncementsPage, "text/html")).Methods("GET")
	router.HandleFunc("/animaux/perdu/nouveau", GziperHandler(GetLostFormPage, "text/html")).Methods("GET")
	router.HandleFunc("/animaux/errant/nouveau", GziperHandler(GetSeenFormPage, "text/html")).Methods("GET")
	router.HandleFunc("/animaux/adopter/nouveau", GziperHandler(GetAdoptFormPage, "text/html")).Methods("GET")
	router.HandleFunc("/animaux/id/{announcementId}", GziperHandler(GetAnnouncementPage, "text/html")).Methods("GET")

	router.HandleFunc("/login", Login)
	router.HandleFunc("/admin", GetAdminHandler)
	router.HandleFunc("/admin/login", authentication.Login)
	router.Handle("/admin/announcements", SecureHandler(GetAdminAnnouncementsHandler))
	router.Handle("/ws/admin/announcements/id/", SecureHandler(GetAdminLostActionHandler))
	router.Handle("/ws/admin/announcements/lost/all", SecureHandler(GetAllHandler))

	router.HandleFunc("/socket", websocket.WsHandler).Methods("GET")
	router.HandleFunc("/ws/mail", GetMailHandler).Methods("POST")
	router.HandleFunc("/ws/species", specy.GetSpeciesHandler).Methods("GET")
	router.HandleFunc("/ws/contact/message", GetContactMessageHandler).Methods("POST")
	router.HandleFunc("/ws/animaux", GetAnnouncementsHandler).Methods("GET")
	router.HandleFunc("/ws/animaux/{announcementType}", GetAnnouncementsHandler).Methods("GET")
	router.HandleFunc("/ws/animaux/{announcementType}", CreateAnnouncementsHandler).Methods("POST")
	router.Handle("/ws/animaux/id/{announcementId}", SecureHandler(GetLostActionHandler)).Methods("PUT")
	router.HandleFunc("/ws/animaux/perdu/locations", GetLocationsHandler).Methods("PUT")

	http.Handle("/", router)

	fmt.Printf("Running on port %s...\n", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Printf("Erreur au démarrage du serveur : %s\n", err.Error())
		panic("Erreur au démarrage du serveur")
	}
}
