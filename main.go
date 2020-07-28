package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Picha, share photos securely!</h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "to get in touch, please send us an email "+"to <a href=\"mailto:support@picha.com\">"+"support@picha.com</a>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":9000", nil)
}
