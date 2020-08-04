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
		NewView:   view.New("appcontainer", "user/new"),
		LoginView: view.New("appcontainer", "user/login"),
		us:        us,
	}
}

// User represents a user in our application
type User struct {
	NewView   *view.View
	LoginView *view.View
	us        *model.UserService
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

// Login authenticates a user
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case model.ErrNotFound:
		fmt.Fprintln(w, "invalid email address")
	case model.ErrInvalidPassword:
		fmt.Fprintln(w, "invalid password provided")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SignupForm captures user input from the sign up forms
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm captures user input from the log in form
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
