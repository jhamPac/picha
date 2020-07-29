package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/view"
)

var homeTemplate *view.View
var contactTemplate *view.View

func main() {
	homeTemplate = view.New("templates/home.gohtml")
	contactTemplate = view.New("templates/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.NotFoundHandler = h

	http.ListenAndServe(":9000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>ARE YOU LOST?</h1>")
}

var h http.Handler = http.HandlerFunc(notfound)

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}
