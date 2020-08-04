package controller

import (
	"fmt"
	"net/http"

	"github.com/jhampac/picha/view"
)

// NewUser instantiates and returns a *User type
func NewUser() *User {
	return &User{
		NewView: view.New("appcontainer", "user/new"),
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
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Name is", form.Name)
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}

// SignupForm captures user input from forms
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
