package minify

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reunion/configuration"
)

var MinifiedContents = make(map[string]string)

func Store(filename string, minified string) {
	MinifiedContents[filename] += minified
}

func MinifyFromFile(filepath string, minifyUrl string, fileType string) {
	content, err := ioutil.ReadFile(configuration.Conf.GetFilePath(filepath))
	if err != nil {
		log.Panicf("Erreur de compression de fichier : %s", err.Error())
	}
	Store(fileType, Minify(string(content), minifyUrl))
}

func Minify(content string, minifyUrl string) string {
	response, err := http.PostForm(minifyUrl, url.Values{"input": {content}})
	if err != nil || response.StatusCode != 200 {
		log.Printf("Erreur de compression de fichier lors de l'appel au WS : %s", minifyUrl)
		return content
	}
	defer response.Body.Close()
	var minifiedContent []byte
	if minifiedContent, err = ioutil.ReadAll(response.Body); err != nil {
		log.Panicf("Erreur de compression de fichier : %s", err.Error())
	}
	return string(minifiedContent)
}

func write(writer http.ResponseWriter, filename string) {
	writer.Write([]byte(MinifiedContents[filename]))
}
