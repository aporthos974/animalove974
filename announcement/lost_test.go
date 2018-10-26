package announcement

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reunion/configuration"
	"testing"
	"time"
)

func TestCreateLostAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	announcement := &Announcement{Description: "new announcement", PhoneNumber: "0606060606", Type: "seen", State: "validated", Account: &Account{}}

	Create(announcement)

	announcements, _ := GetAnnouncementsByType("seen", 0, 1)
	assert.Equal(test, "new announcement", announcements[0].Description)
	assert.Equal(test, "seen", announcements[0].Type)
}

func TestGetLostAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	first := Announcement{Description: "description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"}
	second := Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"}
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(first)
	collection.Insert(second)

	announcements, total := GetAnnouncementsByType("lost", 0, 2)

	assert.Len(test, announcements, 2)
	assert.Equal(test, 2, total)
	assert.Equal(test, first.Description, announcements[0].Description)
	assert.Equal(test, second.Description, announcements[1].Description)
}

func TestGetAllAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	first := Announcement{Description: "description test", PhoneNumber: "0606060606", Type: "lost", State: "deactivated"}
	second := Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "deleted"}
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(first)
	collection.Insert(second)

	announcements := GetAll(0, 2, "deleted")

	assert.Len(test, announcements, 1)
	assert.Equal(test, second.Description, announcements[0].Description)
	assert.Equal(test, "deleted", announcements[0].State)
}

func TestGetAllWithOffsetAndLimit(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Description: "description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "pre last description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "last description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})

	announcements := GetAll(3, 5, "validated")

	assert.Len(test, announcements, 2)
	assert.Equal(test, "pre last description test", announcements[0].Description)
	assert.Equal(test, "last description test", announcements[1].Description)
}

func TestFindByIdLostAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Description: "description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace823"), Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})

	announcement := FindById("54f47ae176a10ee8daace822")

	assert.Equal(test, "description test", announcement.Description)
}

func TestFindLostAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Description: "description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "pre last description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "pre last description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "last description test", PhoneNumber: "0606060606", Type: "lost", State: "validated"})

	announcements, total := Find("pre", 1, 1)

	assert.Len(test, announcements, 1)
	assert.Equal(test, 2, total)
	assert.Equal(test, "pre last description test", announcements[0].Description)
}

func TestFindOnlyLostAnimalsByCity(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Description: "first description test", PhoneNumber: "0606060606", City: "Saint Denis", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", City: "Saint Denis", State: "validated"})
	collection.Insert(Announcement{Description: "third description test", PhoneNumber: "0606060606", City: "Saint Denis", State: "found"})
	collection.Insert(Announcement{Description: "last description test", PhoneNumber: "0606060606", City: "Saint Denis", State: "deleted"})

	announcements, total := FindByCity("Saint Denis", 0, 2)

	assert.Len(test, announcements, 2)
	assert.Equal(test, 2, total)
	assert.Equal(test, "first description test", announcements[0].Description)
	assert.Equal(test, "second description test", announcements[1].Description)
}

func TestFindByCityLostAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Description: "description test", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "validated"})
	collection.Insert(Announcement{Description: "second description test", PhoneNumber: "0606060606", City: "stemarie", Type: "lost", State: "validated"})

	announcements, total := FindByCity("stemarie", 0, 10)

	assert.Len(test, announcements, 1)
	assert.Equal(test, 1, total)
	assert.Equal(test, "second description test", announcements[0].Description)
	assert.Equal(test, "stemarie", announcements[0].City)
}

func TestUpdateAnnouncementForAnimalFound(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	insertedAnnouncement := Announcement{Description: "updated", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "lost"}
	collection.Insert(Announcement{Description: "not updated", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "lost"})
	collection.Insert(insertedAnnouncement)
	collection.Insert(Announcement{Description: "not updated", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "lost"})
	collection.Find(bson.M{"description": "updated"}).One(&insertedAnnouncement)

	UpdateState(insertedAnnouncement.Id, "found")

	var announcements []Announcement
	collection.Find(bson.M{}).All(&announcements)
	assert.Equal(test, "not updated", announcements[0].Description)
	assert.Equal(test, "lost", announcements[0].State)
	assert.Equal(test, "updated", announcements[1].Description)
	assert.Equal(test, "found", announcements[1].State)
	assert.Equal(test, "not updated", announcements[2].Description)
	assert.Equal(test, "lost", announcements[2].State)
}

func TestDontUpdateWithWrongState(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	insertedAnnouncement := Announcement{Description: "updated", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "waiting for validation"}
	collection.Insert(insertedAnnouncement)
	collection.Find(bson.M{"description": "updated"}).One(&insertedAnnouncement)

	assert.Panics(test, func() { UpdateState(insertedAnnouncement.Id, "toto") }, "Calling UpdateState() should panic")

	var announcement Announcement
	collection.Find(bson.M{}).One(&announcement)
	assert.Equal(test, "waiting for validation", announcement.State)
}

func TestSaveLocationAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()

	Create(&Announcement{Locations: []*Location{&Location{Latitude: -55.1232343, Longitude: 1.12122332, Date: time.Now()}}, Account: &Account{}})

	var announcements []Announcement
	collection.Find(bson.M{}).All(&announcements)
	assert.Equal(test, -55.1232343, announcements[0].Locations[0].Latitude)
	assert.Equal(test, 1.12122332, announcements[0].Locations[0].Longitude)
	assert.NotNil(test, announcements[0].Locations[0].Date)
}

func TestAddLocationAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Locations: []*Location{&Location{Latitude: 1, Longitude: 2}}, Account: &Account{Email: "toto@test"}})

	AddLocation(Announcement{Id: bson.ObjectIdHex("54f47ae176a10ee8daace822"), Locations: []*Location{&Location{Latitude: 1, Longitude: 2}, &Location{Latitude: -55.1232343, Longitude: 1.12122332, Date: time.Now()}}, Account: &Account{Email: "toto@test"}})

	var announcements []Announcement
	collection.Find(bson.M{}).All(&announcements)
	locations := announcements[0].Locations
	assert.Len(test, locations, 2)
	assert.Equal(test, 1, locations[0].Latitude)
	assert.Equal(test, 2, locations[0].Longitude)
	assert.Equal(test, -55.1232343, locations[1].Latitude)
	assert.Equal(test, 1.12122332, locations[1].Longitude)
	assert.NotNil(test, locations[1].Date)
}

func GetAdsCollection() *mgo.Collection {
	dbConf := configuration.GetConfiguration().GetDatabase()
	session, _ := mgo.Dial(dbConf.Url)
	database := session.DB(dbConf.Name)
	database.Login(dbConf.Username, dbConf.Password)
	collection := database.C("announcement")
	return collection
}
