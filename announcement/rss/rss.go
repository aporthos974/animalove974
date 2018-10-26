package rss

import (
	"fmt"
	"net/http"
	"reunion/announcement"
	"reunion/configuration"
)

func GetFile(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(GetRss()))
	writer.Header().Set("Content-Type", "application/xml")
}

func GetRss() string {
	announcements, _ := announcement.GetAnnouncementsByType("", 0, 25)
	return Get(announcements)
}

func Get(announcements []announcement.Announcement) string {
	base := `<?xml version="1.0" ?><rss version="2.0"><channel><title>AnimaLove - Liste des animaux perdus à la Réunion</title><link>http://` + configuration.Conf.GetUrl() + `</link>
	<description>AnimaLove propose un système de recherche par mot clef d'animaux perdus à la Réunion. Les personnes qui auraient aperçu un animal perdu peuvent indiquer sur une carte satellite sa localisation, ce qui permettra aux propriétaires de le retrouver plus rapidement.</description>%s</channel></rss>`
	var xmlItems string
	for _, ads := range announcements {
		xmlAnnouncement := fmt.Sprintf("<guid>%s</guid>", ads.Id.Hex())
		xmlAnnouncement += fmt.Sprintf("<pubDate>%s</pubDate>", ads.CreationDate)
		xmlAnnouncement += fmt.Sprintf("<title>%s</title>", getTitle(ads))
		xmlAnnouncement += fmt.Sprintf("<thumbnail>%s</thumbnail>", ads.GetImageLink())
		xmlDescription := ads.Description
		xmlDescription += "&lt;br&gt;Lien : " + ads.GetLink()
		xmlDescription += "&lt;br&gt;C'est un " + ads.Animal
		if ads.Specy != nil {
			xmlDescription += "&lt;br&gt;" + ads.Specy.Name
		}
		xmlDescription += "&lt;br&gt;Couleur : " + ads.Color
		xmlDescription += "&lt;br&gt;" + ads.PhoneNumber
		xmlAnnouncement += fmt.Sprintf("<description>%s</description>", xmlDescription)
		xmlItems += fmt.Sprintf("<item>%s</item>", xmlAnnouncement)
	}
	return fmt.Sprintf(base, xmlItems)
}

func getTitle(ads announcement.Announcement) string {
	title := "Animal " + ads.Type
	if ads.Name != "" {
		title += " - " + ads.Name
	}
	return title
}
