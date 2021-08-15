package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/tmkelly28/gowiki/routes"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}

	return m[2], nil
}

func makeHandler(handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
		handler(w, r, title)
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(routes.ViewHandler))
	http.HandleFunc("/edit/", makeHandler(routes.EditHandler))
	http.HandleFunc("/save/", makeHandler(routes.SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
