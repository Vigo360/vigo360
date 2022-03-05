package routes

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

//go:embed includes
var includes embed.FS

// Map each file to its etag
var etags map[string]string = make(map[string]string)

func includesHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	inm := r.Header.Get("If-None-Match")
	if etags[file] == inm && inm != "" {
		w.WriteHeader(304)
		return
	}

	ext := regexp.MustCompile(`\.[A-Za-z]+$`).FindString(file)
	bytes, err := includes.ReadFile("includes/" + file)
	if err != nil {
		// TODO error handling
		log.Fatalf(err.Error())
	}

	// ETag for file not calculated
	if _, ok := etags[file]; !ok {
		etags[file] = GenerateEtag(bytes)
	}

	w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	w.Header().Add("Cache-Control", "max-age=3600")
	w.Header().Add("ETag", etags[file])
	w.Write(bytes)
}

func GenerateEtag(body []byte) string {
	hasher := sha1.New()
	hasher.Write(body)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
