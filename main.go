package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Picha, share your photos securely</h2>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":9000", nil)
}
