package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jhampac/picha/model"
	"github.com/jhampac/picha/view"
)

// NewUser instantiates and returns a *User type
func NewUser(us *model.UserService) *User {
	return &User{
		NewView: view.New("appcontainer", "user/new"),
		us:      us,
	}
}

// User represents a user in our application
type User struct {
	NewView *view.View
	us      *model.UserService
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
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}
	user := model.User{
		Name:     strings.ToLower(form.Name),
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "user is", user)
}

// SignupForm captures user input from forms
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
