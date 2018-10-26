package announcement

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
	"reunion/configuration"
	"testing"
)

func TestUpdateState(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetAdsCollection()
	collection.DropCollection()
	insertedAnnouncement := Announcement{Description: "updated", PhoneNumber: "0606060606", City: "stdenis", Type: "lost", State: "waiting for validation"}
	collection.Insert(insertedAnnouncement)
	collection.Find(bson.M{"description": "updated"}).One(&insertedAnnouncement)

	UnrestrictedUpdateState(insertedAnnouncement.Id, "validated")

	var announcement Announcement
	collection.Find(bson.M{}).One(&announcement)
	assert.Equal(test, "validated", announcement.State)
}
