package middleware

import (
	"net/http"

	"github.com/jhampac/picha/context"
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
			http.Redirect(w, r, "/login", http.StatusFound) // 302 Found as in redirected to login page
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		// pushes to next call
		next(w, r)
	})
}

// Apply middleware step to routes that are configured with http.Handler (ServeHTTP)
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
