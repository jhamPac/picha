package middleware

import (
	"fmt"
	"net/http"

	"github.com/jhampac/picha/model"
)

type RequireUser struct {
	model.UserService
}

// ApplyFn chains to the next call
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}
		fmt.Println("User found: ", user)

		// pushes to next call
		next(w, r)
	})
}
