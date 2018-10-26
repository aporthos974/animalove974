package configuration

import (
	"labix.org/v2/mgo"
	"log"
	"os"
)

const ERROR_MESSAGE = "Erreur lors du dialogue avec la base de donnée : %s"

func GetAnnouncementCollection() (*mgo.Collection, *mgo.Session) {
	db := GetConfiguration().GetDatabase()
	session, err := mgo.Dial(os.Getenv("MONGOLAB_URI"))
	if err != nil {
		log.Printf(ERROR_MESSAGE, err.Error())
	}
	database := session.DB(db.Name)
	collection := database.C("announcement")
	return collection, session
}

func GetAccountCollection() (*mgo.Collection, *mgo.Session) {
	db := GetConfiguration().GetDatabase()
	session, err := mgo.Dial(os.Getenv("MONGOLAB_URI"))
	if err != nil {
		log.Printf(ERROR_MESSAGE, err.Error())
		return nil, nil
	}
	database := session.DB(db.Name)
	collection := database.C("account")
	return collection, session
}

func GetUserCollection() *mgo.Collection {
	db := GetConfiguration().GetDatabase()
	session, err := mgo.Dial(os.Getenv("MONGOLAB_URI"))
	if err != nil {
		log.Printf("Erreur lors du dialogue avec la base de donnée : ", err)
	}
	database := session.DB(db.Name)
	collection := database.C("users")
	return collection
}

func GetSpecyCollection() (*mgo.Collection, *mgo.Session) {
	db := GetConfiguration().GetDatabase()
	session, err := mgo.Dial(os.Getenv("MONGOLAB_URI"))
	if err != nil {
		log.Panicf(ERROR_MESSAGE, err.Error())
	}
	database := session.DB(db.Name)
	collection := database.C("species")
	return collection, session
}
