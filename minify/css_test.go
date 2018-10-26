package minify

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reunion/configuration"
	"testing"
)

func TestMinifyContentCss(test *testing.T) {
	minified := Minify(".test { color: red; }", urlCss)

	assert.Equal(test, ".test{color:red}", minified)
}

func TestMinifyCssFromFile(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	file, _ := ioutil.TempFile("", "")
	file.WriteString(".test { color: red; }")

	MinifyCss(file.Name())

	assert.Equal(test, ".test{color:red}", MinifiedContents["styles"])
}
