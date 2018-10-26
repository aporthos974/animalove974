package specy

import (
	"encoding/json"
	. "labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"reunion/configuration"
	"strings"
)

var SPECIES_URL = map[string]string{"Chien": "http://wamiz.com/chiens/race-chien/races", "Chat": "http://wamiz.com/chats/race-chat/races"}

type Species []Specy

type Specy struct {
	Name, Picture, Animal string
}

func (species Species) Len() int {
	return len(species)
}

func (species Species) Swap(i, j int) {
	species[i], species[j] = species[j], species[i]
}

func (species Species) Less(i, j int) bool {
	return strings.Replace(species[i].Name, "É", "E", 1) < strings.Replace(species[j].Name, "É", "E", 1)
}

func GetSpeciesHandler(writer http.ResponseWriter, request *http.Request) {
	jsonSpecies, err := json.Marshal(GetSpecies())
	if err != nil {
		log.Panicf("Erreur lors de la conversion JSON : %s", err.Error())
	}
	writer.Write(jsonSpecies)
}

func GetSpecies() Species {
	if isAlreadyFetchedFromSite() {
		return getSpeciesFromCollection()
	}
	log.Panicf("Aucune donnée concernant les races d'animaux")
	return nil
}

func isAlreadyFetchedFromSite() bool {
	collection, session := configuration.GetSpecyCollection()
	defer session.Close()

	count, err := collection.Find(M{}).Count()
	logError(err)
	return count > 0
}

func getSpeciesFromCollection() (species Species) {
	collection, session := configuration.GetSpecyCollection()
	defer session.Close()

	collection.Find(M{}).All(&species)
	return species
}

func logError(err error) {
	if err != nil {
		log.Panicf("Erreur lors du parsing de la page des races : %s", err.Error())
	}
}
