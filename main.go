package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/controller"
	"github.com/jhampac/picha/view"
)

var (
	homeView    *view.View
	contactView *view.View
)

func main() {
	homeView = view.New("appcontainer", "templates/home.gohtml")
	contactView = view.New("appcontainer", "templates/contact.gohtml")
	userC := controller.NewUser()

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", userC.SignUp)
	r.NotFoundHandler = h

	http.ListenAndServe(":9000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>ARE YOU LOST?</h1>")
}

var h http.Handler = http.HandlerFunc(notfound)

// must as in error must be nil
func must(err error) {
	if err != nil {
		panic(err)
	}
}
