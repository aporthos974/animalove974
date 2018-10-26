package announcement

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"github.com/disintegration/imaging"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	. "labix.org/v2/mgo/bson"
	"log"
	"os/exec"
	"regexp"
	"reunion/announcement/specy"
	"reunion/configuration"
	"strings"
	"time"
)

type Announcement struct {
	XMLName                xml.Name     `bson:"-" xml:"item"`
	Id                     ObjectId     `bson:"_id,omitempty" xml:"-"`
	Name                   string       `xml:"title,omitempty"`
	Description            string       `xml:"description,omitempty"`
	City, PhoneNumber, Sex string       `xml:",omitempty"`
	Color, Animal          string       `xml:",omitempty"`
	Type                   string       `xml:"-"`
	Picture                string       `xml:"-"`
	State                  string       `xml:"-"`
	CreationDate, LostDate *time.Time   `xml:",omitempty"`
	Account                *Account     `xml:",omitempty"`
	Specy                  *specy.Specy `xml:",omitempty"`
	Link                   string       `xml:"link,omitempty"`
	Locations              []*Location  `bson:"locations,omitempty" json:",omitempty"`
}

type Location struct {
	Latitude, Longitude float64
	Date                time.Time `json:",omitempty"`
}

type Page struct {
	name, Title, Description, Image, Link string
}

func (announcement *Announcement) GetLink() string {
	return configuration.Conf.GetUrl() + "/animaux/id/" + announcement.Id.Hex()
}

func (announcement *Announcement) GetImageLink() string {
	return configuration.Conf.GetUrl() + "/photos/" + announcement.Picture
}

func GetAnnouncementsByType(typeAnnouncement string, offset int, limit int) ([]Announcement, int) {
	announcements := []Announcement{}

	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	filter := M{"state": M{"$in": []string{"validated", "waiting for validation", "found"}}}
	if typeAnnouncement != "" {
		filter["type"] = typeAnnouncement
	}
	if err := collection.Find(M{"$and": []interface{}{filter}}).Sort("-creationdate").Skip(offset).Limit(limit).All(&announcements); err != nil {
		log.Panicf("Erreur lors de la récupération des annonces : %s", err.Error())
	}
	total, err := collection.Find(filter).Count()
	if err != nil {
		log.Panicf("Erreur lors de la récupération du nombre total d'annonces : %s", err.Error())
	}
	return announcements, total
}

func GetAll(offset int, limit int, state string) []Announcement {
	announcements := []Announcement{}

	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.Find(M{"state": state}).Sort("-creationdate").Skip(offset).Limit(limit).All(&announcements); err != nil {
		log.Panicf("Erreur lors de la récupération des annonces : %s", err.Error())
	}
	return announcements
}

func Find(criteria string, offset int, limit int) (announcements []Announcement, total int) {
	regexp := RegEx{Pattern: `.*` + criteria + `.*`, Options: "i"}
	matches := M{"$and": []interface{}{M{"state": M{"$in": []string{"validated", "waiting for validation", "found"}}}, M{"$or": []interface{}{M{"description": regexp}, M{"name": regexp}, M{"city": regexp}, M{"type": regexp}, M{"color": regexp}, M{"specy.name": regexp}, M{"animal": regexp}, M{"phonenumber": regexp}}}}}
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.Find(matches).Sort("-creationdate").Skip(offset).Limit(limit).All(&announcements); err != nil {
		log.Panicf("Erreur lors de la récupération des annonces : %s", err.Error())
	}
	total, err := collection.Find(matches).Count()
	if err != nil {
		log.Panicf("Erreur lors de la récupération du nombre total d'annonces : %s", err.Error())
	}

	return announcements, total
}

func FindByCity(city string, offset int, limit int) (announcements []Announcement, total int) {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.Find(M{"$and": []interface{}{M{"state": M{"$in": []string{"validated", "waiting for validation"}}}, M{"city": city}}}).Sort("-creationdate").Skip(offset).Limit(limit).All(&announcements); err != nil {
		log.Panicf("Erreur lors de la récupération des annonces : %s", err.Error())
	}
	total, err := collection.Find(M{"$and": []interface{}{M{"state": M{"$in": []string{"validated", "waiting for validation"}}}, M{"city": city}}}).Count()
	if err != nil {
		log.Panicf("Erreur lors de la récupération du nombre total d'annonces : %s", err.Error())
	}

	return announcements, total
}

func FindById(id string) (announcement Announcement) {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.Find(M{"_id": ObjectIdHex(id)}).One(&announcement); err != nil {
		log.Panicf("Erreur lors de la récupération de l'annonce : %s", err.Error())
	}

	return announcement
}

func Create(announcement *Announcement) {
	now := time.Now()
	announcement.CreationDate = &now
	announcement.State = "waiting for validation"
	announcement.Account.Password = EncryptSHA256(announcement.Account.Password)
	announcement.Picture = saveImage(announcement.Picture)

	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.Insert(announcement); err != nil {
		log.Panicf("Erreur lors de l'ajout d'une nouvelle annonce : %s", err.Error())
	}
	sendMail(Mail{Sender: configuration.Conf.GetMail().Contact, Recipient: configuration.Conf.GetMail().Admin}, "Un nouvel animal perdu a été publié : "+announcement.CreationDate.String()+" \nVoici son email :"+announcement.Account.Email)
	sendMail(Mail{Sender: configuration.Conf.GetMail().Contact, Recipient: announcement.Account.Email}, "Bonjour,\n\nNous vous confirmons la publication de l'animal perdu sur AnimaLove.\nVoici le lien vers cet animal : "+announcement.GetLink())
}

func saveImage(imageContent string) string {
	if imageContent == "" {
		return ""
	}

	uuid, err := exec.Command("uuidgen", "-r").Output()
	if err != nil {
		log.Panicf("Erreur dans la génération du nom du fichier : %s", err.Error())
	}
	pictureName := strings.Replace(string(uuid), "\n", "", 1) + ".jpeg"
	picture := compressAndResizeImage(imageContent)
	ioutil.WriteFile("/opt/animalove/photos/"+pictureName, picture, 0664)
	return pictureName
}

func compressAndResizeImage(pictureInString string) []byte {
	if pictureInString == "" {
		return nil
	}
	compile := regexp.MustCompile("data:image/.+;base64,")
	pictureWithoutHeader := compile.ReplaceAllString(pictureInString, "")
	pictureInBytes := make([]byte, base64.StdEncoding.DecodedLen(len(pictureWithoutHeader)))
	_, err := base64.StdEncoding.Decode(pictureInBytes, []byte(pictureWithoutHeader))
	if err != nil {
		log.Panicf("Erreur lors de la compression de l'image : %s", err.Error())
	}
	picture, _, err := image.Decode(bytes.NewReader(pictureInBytes))
	if err != nil {
		log.Panicf("Erreur lors de la compression de l'image : %s", err.Error())
	}
	resizedImage := picture
	ratio := picture.Bounds().Dx() / 400
	if ratio > 0 {
		resizedImage = imaging.Resize(picture, picture.Bounds().Dx()/ratio, picture.Bounds().Dy()/ratio, imaging.NearestNeighbor)
	}

	var buffer bytes.Buffer
	jpeg.Encode(&buffer, resizedImage, &jpeg.Options{Quality: 60})
	return buffer.Bytes()
}

func AddLocation(announcement Announcement) {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	err := collection.UpdateId(announcement.Id, M{"$set": M{"locations": announcement.Locations}})
	if err != nil {
		log.Panicf("Erreur lors de la modification des localisations : %s", err.Error())
	}
	mailContent := Mail{Sender: configuration.Conf.GetMail().Contact, Recipient: announcement.Account.Email}
	sendMail(mailContent, "Bonjour,\n\nQuelqu'un a localisé votre animal sur la carte \nVous pouvez la consulter sur ce lien : "+announcement.GetLink())
}

func UpdateState(announcementId ObjectId, state string) {
	if state != "deleted" && state != "found" && state != "adopted" {
		log.Panicf("Statut de la demande de mise à jour invalide")
	}
	UnrestrictedUpdateState(announcementId, state)
}

func UnrestrictedUpdateState(announcementId ObjectId, state string) {
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()

	if err := collection.UpdateId(announcementId, M{"$set": M{"state": state}}); err != nil {
		log.Panicf("Erreur lors de la modification de l'état d'une annonce : %s", err.Error())
	}
}
