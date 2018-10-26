package cache

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMaxAgeInHttpResponse(test *testing.T) {
	writer := httptest.NewRecorder()

	addCacheHeader(writer, CacheHeader{MaxAge: 600}, nil)

	assert.Equal(test, "max-age=600", writer.Header().Get("Cache-Control"))
}

func TestAddETagInHttpResponse(test *testing.T) {
	writer := httptest.NewRecorder()
	responseContent := []byte("{ toto: 'fait du vélo'}")

	eTag := addCacheHeader(writer, CacheHeader{MaxAge: 600}, responseContent)

	assert.Equal(test, "4804f7dad78e1d3487651eb471fdc8f1", writer.Header().Get("ETag"))
	assert.Equal(test, "4804f7dad78e1d3487651eb471fdc8f1", eTag)
}

func TestCheckRequestETagIsTheSame(test *testing.T) {
	request := &http.Request{Header: http.Header{}}
	request.Header.Add("If-None-Match", "4804f7dad78e1d3487651eb471fdc8f1")

	isTheSame := checkETag(request, "4804f7dad78e1d3487651eb471fdc8f1")

	assert.True(test, isTheSame)
}

func TestCheckRequestETagIsNotTheSame(test *testing.T) {
	request := &http.Request{Header: http.Header{}}
	request.Header.Add("If-None-Match", "4804f7dad78e1d3487651eb471fdc8f2")

	isTheSame := checkETag(request, "4804f7dad78e1d3487651eb471fdc8f1")

	assert.False(test, isTheSame)
}
