package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Home page</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "to get in touch, please send us an email "+"to <a href=\"mailto:support@picha.com\">"+"support@picha.com</a>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":9000", nil)
}
