package rss

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reunion/announcement"
	"reunion/configuration"
	"testing"
)

func TestGetAnnouncementInRSS(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)

	xmlContent := Get([]announcement.Announcement{announcement.Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Name: "Junko", Description: "description de test"}})

	actual := "<?xml version=\"1.0\" ?><rss version=\"2.0\"><channel><title>AnimaLove - Liste des animaux perdus à la Réunion</title><link>http://localhost:8080</link>\n\t<description>AnimaLove propose un système de recherche par mot clef d'animaux perdus à la Réunion. Les personnes qui auraient aperçu un animal perdu peuvent indiquer sur une carte satellite sa localisation, ce qui permettra aux propriétaires de le retrouver plus rapidement.</description><item><guid>54f47ae176a10ee8daace822</guid><pubDate><nil></pubDate><title>Animal  - Junko</title><thumbnail>localhost:8080/photos/</thumbnail><description>description de test&lt;br&gt;Lien : localhost:8080/announcements/id/54f47ae176a10ee8daace822&lt;br&gt;C'est un &lt;br&gt;Couleur : &lt;br&gt;</description></item></channel></rss>"
	assert.Equal(test, actual, xmlContent)
}

func TestGetSeveralAnnouncementsInRSS(test *testing.T) {
	announcements := []announcement.Announcement{announcement.Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Name: "Junko", Description: "description de test", Type: "perdu"}, announcement.Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace823"), Name: "Qiu", Description: "description de test qiu"}}

	xmlContent := Get(announcements)

	actual := "<?xml version=\"1.0\" ?><rss version=\"2.0\"><channel><title>AnimaLove - Liste des animaux perdus à la Réunion</title><link>http://localhost:8080</link>\n\t<description>AnimaLove propose un système de recherche par mot clef d'animaux perdus à la Réunion. Les personnes qui auraient aperçu un animal perdu peuvent indiquer sur une carte satellite sa localisation, ce qui permettra aux propriétaires de le retrouver plus rapidement.</description><item><guid>54f47ae176a10ee8daace822</guid><pubDate><nil></pubDate><title>Animal perdu - Junko</title><thumbnail>localhost:8080/photos/</thumbnail><description>description de test&lt;br&gt;Lien : localhost:8080/announcements/id/54f47ae176a10ee8daace822&lt;br&gt;C'est un &lt;br&gt;Couleur : &lt;br&gt;</description></item><item><guid>54f47ae176a10ee8daace823</guid><pubDate><nil></pubDate><title>Animal  - Qiu</title><thumbnail>localhost:8080/photos/</thumbnail><description>description de test qiu&lt;br&gt;Lien : localhost:8080/announcements/id/54f47ae176a10ee8daace823&lt;br&gt;C'est un &lt;br&gt;Couleur : &lt;br&gt;</description></item></channel></rss>"
	assert.Equal(test, actual, xmlContent)
}

func TestGetAnnouncementRSSInDB(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	announcement := announcement.Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Name: "Junko", Description: "description de test", State: "validated", Type: "perdu"}
	collection := GetAnnouncementCollection()
	collection.DropCollection()
	collection.Insert(announcement)

	xmlContent := GetRss()

	actual := "<?xml version=\"1.0\" ?><rss version=\"2.0\"><channel><title>AnimaLove - Liste des animaux perdus à la Réunion</title><link>http://localhost:8080</link>\n\t<description>AnimaLove propose un système de recherche par mot clef d'animaux perdus à la Réunion. Les personnes qui auraient aperçu un animal perdu peuvent indiquer sur une carte satellite sa localisation, ce qui permettra aux propriétaires de le retrouver plus rapidement.</description><item><guid>54f47ae176a10ee8daace822</guid><pubDate><nil></pubDate><title>Animal perdu - Junko</title><thumbnail>localhost:8080/photos/</thumbnail><description>description de test&lt;br&gt;Lien : localhost:8080/announcements/id/54f47ae176a10ee8daace822&lt;br&gt;C'est un &lt;br&gt;Couleur : &lt;br&gt;</description></item></channel></rss>"
	assert.Equal(test, actual, xmlContent)
}

func GetAnnouncementCollection() *mgo.Collection {
	dbConf := configuration.GetConfiguration().GetDatabase()
	session, _ := mgo.Dial(dbConf.Url)
	database := session.DB(dbConf.Name)
	database.Login(dbConf.Username, dbConf.Password)
	collection := database.C("announcement")
	return collection
}
