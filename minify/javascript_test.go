package minify

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestMinifyContentJs(test *testing.T) {
	minified := Minify("var test = 'essai' ; console.log( test ) ;", urlJs)

	assert.Equal(test, "var test=\"essai\";console.log(test);", minified)
}

func TestMinifyJsFromFile(test *testing.T) {
	file, _ := ioutil.TempFile("", "")
	file.WriteString("var test = 'essai' ; console.log( test ) ;")

	MinifyJs(file.Name())

	assert.Equal(test, "var test=\"essai\";console.log(test);", MinifiedContents["javascript"])
}
