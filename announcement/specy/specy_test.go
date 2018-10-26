package specy

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"reunion/configuration"
	"testing"
)

func TestGetSpeciesOfDog(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	collection := GetSpecyCollection("species")
	collection.DropCollection()
	collection.Insert(Specy{Name: "Race 1", Picture: "test.jpg", Animal: "Chien"})
	collection.Insert(Specy{Name: "Race 2", Picture: "test.jpg", Animal: "Chat"})

	species := GetSpecies()

	assert.Len(test, species, 2)
	assert.Equal(test, Specy{Name: "Race 1", Picture: "test.jpg", Animal: "Chien"}, species[0])
	assert.Equal(test, Specy{Name: "Race 2", Picture: "test.jpg", Animal: "Chat"}, species[1])
}

func GetSpecyCollection(collectionName string) *mgo.Collection {
	dbConf := configuration.GetConfiguration().GetDatabase()
	session, _ := mgo.Dial(dbConf.Url)
	database := session.DB(dbConf.Name)
	database.Login(dbConf.Username, dbConf.Password)
	collection := database.C(collectionName)
	return collection
}
