package announcement

import (
	"github.com/stretchr/testify/assert"
	"reunion/configuration"
	"testing"
)

func TestGetSeenAnnouncement(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	first := Announcement{Description: "description test", PhoneNumber: "0606060606", Type: "seen", State: "validated"}
	second := Announcement{Description: "second description test", PhoneNumber: "0606060606", Type: "seen", State: "validated"}
	collection := GetAdsCollection()
	collection.DropCollection()
	collection.Insert(first)
	collection.Insert(second)

	announcements, total := GetAnnouncementsByType("seen", 0, 2)

	assert.Len(test, announcements, 2)
	assert.Equal(test, 2, total)
	assert.Equal(test, first.Description, announcements[0].Description)
	assert.Equal(test, second.Description, announcements[1].Description)
}
