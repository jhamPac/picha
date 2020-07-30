package controller

import (
	"net/http"

	"github.com/jhampac/picha/view"
)

// NewUser instantiates and returns a *User type
func NewUser() *User {
	return &User{
		NewView: view.New("appcontainer", "templates/user/new.gohtml"),
	}
}

// User represents a user in our application
type User struct {
	NewView *view.View
}

// SignUp is the handler used to sign a new user up
func (u *User) SignUp(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
