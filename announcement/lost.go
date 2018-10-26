package announcement

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"regexp"
	"reunion/cache"
	"reunion/configuration"
	"reunion/websocket"
	"strconv"
	"strings"
)

const MAX_RESULT = 300

var ANNOUNCEMENT_TYPE = []string{"perdu", "errant", "adopter"}

func GetAnnouncementPage(writer http.ResponseWriter, request *http.Request) {
	announcementId := mux.Vars(request)["announcementId"]

	if announcementId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tmpl, page := GetOnlyAnouncementPage(announcementId)
	writer.Header().Set("Cache-control", "public, max-age=3600")
	error := tmpl.Execute(writer, page)
	if error != nil {
		log.Panicf("Erreur lors de la construction du template : %s", error.Error())
	}
}

func GetOnlyAnouncementPage(id string) (tmpl *template.Template, page Page) {
	announcement := FindById(id)

	tmpl = template.Must(template.New("base").Delims("[[", "]]").ParseFiles(configuration.Conf.GetFilePath("static/html/location.html"), configuration.Conf.GetFilePath("static/html/contact.html"), configuration.Conf.GetFilePath("static/html/auth.html"), configuration.Conf.GetFilePath("static/html/result.html"), configuration.Conf.GetFilePath("static/html/announcement.html"), configuration.Conf.GetFilePath("static/html/base.html")))
	page = Page{name: "annonce", Title: "AnimaLove - " + announcement.Description, Description: getMetaDescription(announcement), Image: announcement.GetImageLink(), Link: announcement.GetLink()}
	return tmpl, page
}

func getMetaDescription(announcement Announcement) string {
	return "Animal perdu - " + announcement.Description + "\n" + announcement.GetLink()
}

func GetAnimalPicture(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Printf("Erreur id manquant")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("id manquant"))
		return
	}
	pictureContent := FindById(id).Picture
	compile := regexp.MustCompile("data:(image/.+);.*$")
	contentType := compile.FindStringSubmatch(pictureContent)[1]
	content, err := base64.StdEncoding.DecodeString(strings.Replace(pictureContent, "data:"+contentType+";base64,", "", 1))
	if err != nil {
		log.Panicf("Erreur lors du décodage de l'image : %s", err.Error())
	}
	writer.Header().Add("Content-Type", contentType)
	writer.Write(content)
}

func GetAnnouncementsPage(writer http.ResponseWriter, request *http.Request) {
	GetLostTemplatePage(writer, request, "AnimaLove - Liste des animaux perdus à la Réunion", "loss_announcement", "Liste des chiens et chats perdus, errants et à adopter à la Réunion. Catégorisation des animaux.")
}

func GetLostFormPage(writer http.ResponseWriter, request *http.Request) {
	GetLostTemplatePage(writer, request, "AnimaLove - Signaler un animal perdu ou trouvé à la Réunion", "create_loss_announcement", "Formulaire de signalement d'un chien ou un chat perdu. Ce formulaire permet de décrire l'animal et d'indiquer le lieux où il a été vu pour la dernière fois.")
}

func GetLostTemplatePage(writer http.ResponseWriter, request *http.Request, title string, htmlFilename string, description string) {
	tmpl := template.Must(template.New("base").Delims("[[", "]]").ParseFiles(configuration.Conf.GetFilePath("static/html/location.html"), configuration.Conf.GetFilePath("static/html/contact.html"), configuration.Conf.GetFilePath("static/html/auth.html"), configuration.Conf.GetFilePath("static/html/result.html"), configuration.Conf.GetFilePath("static/html/search.html"), configuration.Conf.GetFilePath("static/html/"+htmlFilename+".html"), configuration.Conf.GetFilePath("static/html/base.html")))
	page := Page{name: "Animal perdu", Title: title, Description: description, Image: "/static/images/logo.png", Link: configuration.Conf.GetUrl()}
	writer.Header().Set("Cache-control", "public, max-age=3600")
	error := tmpl.Execute(writer, page)
	if error != nil {
		log.Panicf("Erreur lors de la construction du template : %s", error.Error())
	}
}

func GetLostActionHandler(writer http.ResponseWriter, request *http.Request) {
	requestToken := request.Header.Get("Authorization")
	id := bson.ObjectIdHex(mux.Vars(request)["announcementId"])
	if action := request.URL.Query().Get("action"); action != "" {
		if CheckTokenOnAnnouncement(id, requestToken) {
			UpdateState(id, action)
		} else {
			log.Printf("Erreur dans la vérification du token : %s, les tokens ne correspondent pas", requestToken)
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Erreur lors de l'identification"))
			return
		}
	} else {
		log.Printf("Paramètre manquant : %s", request.URL.RequestURI())
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func getId(uri string, idValueIndex int) (id bson.ObjectId, err error) {
	splittedUri := strings.Split(uri, "/")
	if len(splittedUri) > idValueIndex {
		endUri := strings.Split(splittedUri[idValueIndex], "?")
		id = bson.ObjectIdHex(endUri[0])
		return id, nil
	}
	err = errors.New("Pas d'identifiant dans l'URL")
	return id, err
}

func GetLocationsHandler(writer http.ResponseWriter, request *http.Request) {
	var announcement Announcement
	err := json.NewDecoder(request.Body).Decode(&announcement)
	if err != nil {
		log.Panicf("Erreur lors de la reconversion JSON : %s", err.Error())
	}
	AddLocation(announcement)
	writer.WriteHeader(http.StatusOK)
}

func GetAllHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		announcements := GetAll(1, 1000, request.URL.Query().Get("state"))
		jsonAnnouncements, err := json.Marshal(announcements)
		if err != nil {
			log.Panicf("Erreur lors de la conversion JSON : %s", err.Error())
		}
		writer.Write(jsonAnnouncements)
	} else {
		log.Printf("Opération non autorisée : %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetAnnouncementsHandler(writer http.ResponseWriter, request *http.Request) {
	typeAnnouncement := mux.Vars(request)["announcementType"]
	if !validateAnnouncementType(typeAnnouncement) {
		log.Panicf("Type incorrect : %s", typeAnnouncement)
	}

	var announcements []Announcement
	var total int
	offset, limit := findAndValidateLimit(request)
	if queryParam := request.URL.Query().Get("criteria"); queryParam != "" {
		announcements, total = Find(queryParam, offset, limit)
	} else if city := request.URL.Query().Get("city"); city != "" {
		announcements, total = FindByCity(city, offset, limit)
	} else if id := request.URL.Query().Get("id"); id != "" {
		announcements = append(announcements, FindById(id))
	} else {
		announcements, total = GetAnnouncementsByType(typeAnnouncement, offset, limit)
	}

	jsonAnnouncements, err := json.Marshal(map[string]interface{}{"announcements": announcements, "total": total})
	if err != nil {
		log.Panicf("Erreur lors de la conversion JSON : %s", err.Error())
	}

	cache.AddHttpCacheContent(writer, request, jsonAnnouncements)
}

func CreateAnnouncementsHandler(writer http.ResponseWriter, request *http.Request) {
	typeAnnouncement := mux.Vars(request)["announcementType"]
	if !validateAnnouncementType(typeAnnouncement) {
		log.Panicf("Type incorrect : %s", typeAnnouncement)
	}

	var announcement Announcement
	err := json.NewDecoder(request.Body).Decode(&announcement)
	if err != nil {
		log.Panicf("Erreur lors de la reconversion JSON : %s", err.Error())
	}
	if announcement.Locations[0] == nil {
		announcement.Locations = []*Location{}
	}
	Create(&announcement)
	notifySeenCreated(writer, request, announcement)
	writer.WriteHeader(http.StatusCreated)
}

func notifySeenCreated(writer http.ResponseWriter, request *http.Request, announcement Announcement) {
	if announcement.Type == "seen" {
		websocket.SendNotification(writer, request, map[string]interface{}{"action": "seen-created", "announcement": announcement})
	}
}

func validateAnnouncementType(typeAnnouncement string) bool {
	if typeAnnouncement == "" {
		return true
	}
	for _, value := range ANNOUNCEMENT_TYPE {
		if value == typeAnnouncement {
			return true
		}
	}
	return false
}

func findAndValidateLimit(request *http.Request) (int, int) {
	offsetParam := request.URL.Query().Get("offset")
	limitParam := request.URL.Query().Get("limit")
	if offsetParam == "" || limitParam == "" {
		log.Panicf("Paramètre manquant limit ou offset")
	}
	limit := convertToInt(limitParam)
	if limit > MAX_RESULT {
		limit = MAX_RESULT
	}
	offset := convertToInt(offsetParam)
	return offset, limit
}

func convertToInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("Erreur lors de la conversion en entier : %s", err.Error())
	}
	return intValue
}
