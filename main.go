package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/controller"
	"github.com/jhampac/picha/model"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "testpassword"
	dbname   = "picha_dev"
)

func main() {
	// db connection and service creation; data layer
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	us, err := model.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	// instatantiate controllers
	staticC := controller.NewStatic()
	userC := controller.NewUser(us)

	// // routing
	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")

	r.Handle("/login", userC.LoginView).Methods("GET")
	r.HandleFunc("/login", userC.Login).Methods("POST")

	r.HandleFunc("/cookietest", userC.CookieTest).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		staticC.Error.ServeHTTP(w, r)
	})

	// initiate app; serve app; accept connections
	http.ListenAndServe(":9000", r)
}
