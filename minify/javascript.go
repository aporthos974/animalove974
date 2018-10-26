package minify

import (
	"net/http"
)

var urlJs = "http://javascript-minifier.com/raw"

func MinifyJs(filepath string) {
	MinifyFromFile(filepath, urlJs, "javascript")
}

func AppJsHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "javascript")
}
