package configuration

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetRobotsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(GetRobots()))
}

func GetRobots() string {
	content, err := ioutil.ReadFile(Conf.GetFilePath("static/robots.txt"))
	if err != nil {
		log.Panicf("Erreur lors la lecture du fichier robots.txt : %s", err.Error())
	}
	return string(content)
}

func GetSitemapHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(GetSitemap()))
	writer.Header().Set("Content-Type", "application/xml")
}

func GetSitemap() string {
	content, err := ioutil.ReadFile(Conf.GetFilePath("static/sitemap.xml"))
	if err != nil {
		log.Panicf("Erreur lors la lecture du fichier robots.txt : %s", err.Error())
	}
	return string(content)
}
