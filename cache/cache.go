package cache

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
)

type CacheHeader struct {
	MaxAge int
}

func checkETag(request *http.Request, newETagHash string) bool {
	return newETagHash == request.Header.Get("If-None-Match")
}

func AddHttpCacheContent(writer http.ResponseWriter, request *http.Request, responseContent []byte) {
	newETagHash := addCacheHeader(writer, CacheHeader{MaxAge: 0}, responseContent)
	if checkETag(request, newETagHash) {
		writer.WriteHeader(http.StatusNotModified)
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write(responseContent)
	}
}

func addCacheHeader(writer http.ResponseWriter, cacheHeader CacheHeader, responseContent []byte) (eTagHash string) {
	writer.Header().Add("Cache-Control", "max-age="+strconv.Itoa(cacheHeader.MaxAge))
	hash := md5.New()
	hash.Write(responseContent)
	eTagHash = hex.EncodeToString(hash.Sum(nil))
	writer.Header().Add("ETag", eTagHash)
	return eTagHash
}
