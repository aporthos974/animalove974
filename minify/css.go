package minify

import (
	"net/http"
)

var urlCss = "http://cssminifier.com/raw"

func MinifyCss(filepath string) {
	MinifyFromFile(filepath, urlCss, "styles")
}

func AppCssHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "styles")
}
