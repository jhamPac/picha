package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/controller"
)

func main() {
	staticC := controller.NewStatic()
	userC := controller.NewUser()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		staticC.Error.ServeHTTP(w, r)
	})

	http.ListenAndServe(":9000", r)
}
