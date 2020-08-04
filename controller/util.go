package controller

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(dst interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
