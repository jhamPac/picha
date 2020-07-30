package controller

import (
	"fmt"
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

// New is the handler used to sign a new user up
func (u *User) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create a new user by handling the request with form data
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "coming soon")
}
