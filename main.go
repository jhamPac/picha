package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/controller"
	"github.com/jhampac/picha/middleware"
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
	services, err := model.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// instatantiate controllers
	staticC := controller.NewStatic()
	userC := controller.NewUser(services.User)
	galleryC := controller.NewGallery(services.Gallery)

	// middleware
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	// // routing
	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")

	r.Handle("/login", userC.LoginView).Methods("GET")
	r.HandleFunc("/login", userC.Login).Methods("POST")

	newGallery := requireUserMw.Apply(galleryC.NewView)
	createGallery := requireUserMw.ApplyFn(galleryC.Create)
	r.Handle("/gallery/new", newGallery).Methods("GET")
	r.HandleFunc("/gallery", createGallery).Methods("POST")

	r.HandleFunc("/cookietest", userC.CookieTest).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		staticC.Error.ServeHTTP(w, r)
	})

	// initiate app; serve app; accept connections
	http.ListenAndServe(":9000", r)
}
