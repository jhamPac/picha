package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jhampac/picha/model"
	"github.com/jhampac/picha/rand"
	"github.com/jhampac/picha/view"
)

// LoginForm captures user input from the log in form
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// SignupForm captures user input from the sign up forms
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// User represents a user in our application
type User struct {
	NewView   *view.View
	LoginView *view.View
	us        model.UserService
}

// NewUser instantiates and returns a *User type
func NewUser(us model.UserService) *User {
	return &User{
		NewView:   view.New("appcontainer", "user/new"),
		LoginView: view.New("appcontainer", "user/login"),
		us:        us,
	}
}

// New is the handler used to sign a new user up
func (u *User) New(w http.ResponseWriter, r *http.Request) {
	type Alert struct {
		Level   string
		Message string
	}

	alert := Alert{
		Level:   "success",
		Message: "Successfully rendered a dynamic alert!",
	}

	if err := u.NewView.Render(w, alert); err != nil {
		panic(err)
	}
}

// Create a new user by handling the request with form data
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	// parse the form and place the results at the address *form
	// gorilla mux schema
	// I like pointers at call-site
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}

	// instantiate a user model with the values from the form
	user := model.User{
		Name:     strings.ToLower(form.Name),
		Email:    form.Email,
		Password: form.Password,
	}

	// create the user in the db with the provided UserService
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remember me token
	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login authenticates a user
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	// vs
	// form := LoginForm{}

	// parse the form and place the results at the address *form
	// gorilla mux schema
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}

	// authenticate the user with the UserService
	user, err := u.us.Authenticate(form.Email, form.Password)

	// check for errors
	if err != nil {
		switch err {
		case model.ErrNotFound:
			fmt.Fprintln(w, "invalid email address")
		case model.ErrPasswordIncorrect:
			fmt.Fprintln(w, "invalid password provided")
		case nil:
			fmt.Fprintln(w, user)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// remember token
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// CookieTest is a debug route for cookies
func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)
}

func (u *User) signIn(w http.ResponseWriter, user *model.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
