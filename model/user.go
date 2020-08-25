package model

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jhampac/picha/hash"
	"github.com/jhampac/picha/rand"
	"github.com/jinzhu/gorm"

	// driver for postgres gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	ErrNotFound = errors.New("model: resource not found")

	// ErrIDInvalid is returned when an invalid ID is provided to a method like Delete
	ErrIDInvalid = errors.New("model: ID provided was invalid")

	// ErrPasswordIncorrect is returned when an invalid password is provided
	ErrPasswordIncorrect = errors.New("model: incorrect password provided")

	// ErrPasswordTooShort is returned when the password provided does not meet the 8 character minimum
	ErrPasswordTooShort = errors.New("model: passwords must be at least 8 characters long")

	// ErrPasswordRequired is returned when a create is attempted without a user password provided
	ErrPasswordRequired = errors.New("model: password is required")

	// ErrEmailRequired is returned when an email address is not provided when creating a user
	ErrEmailRequired = errors.New("model: email address is required")

	// ErrEmailInvalid is returned when an email address provided does not match our regex
	ErrEmailInvalid = errors.New("model: email address is not valid")

	// ErrEmailTaken is returned when an update or create is attempted with an email address that is already in use
	ErrEmailTaken = errors.New("model: email address is already taken")

	// ErrRememberRequired is returned when a create or update is attempted without a user remember token hash
	ErrRememberRequired = errors.New("model: remember token is required")

	// ErrRememberTooShort is returned when a token does not meet the 32 byte minimum
	ErrRememberTooShort = errors.New("model: remember token must be at least 32 bytes")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "model: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

const userPwPepper = "secret-dev-pepper"

const hmacSecretKey = "not-really-a-secret"

// UserDB is an interface to interact with the users db
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// db utilities
	Close() error
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// User represents our customers
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userService implements the UserService interface
type userService struct {
	UserDB
}

// userValidator implements the UserDB; It is a layer that validates and normalizes data before passing it on to the next UserDB layer
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// userGorm implements the UserDB interface
type userGorm struct {
	db *gorm.DB
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

func newUserValidator(orm UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     orm,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// NewUserService instantiates a new service with the provided connection information
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	// interface chaining; validator first then to the gorm/db layer
	return &userService{
		UserDB: uv,
	}, nil
}

// Authenticate users into the app
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}

// Create runs through the validation and normalization layer first
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.emailIsAvail)

	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Create inserts the normalized data into the db
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// ByID queries and returns a user by the id provided
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail validation and normalization layer that passes it to the db layer
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByEmail queries and returns a user by the email provided
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByRemember is the first deferment in the chain to validate and normalize
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// ByRemember looks up a user with the given rememberHash provided by the validation layer
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update is the first deferment in the chain to validate and normalize
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.emailIsAvail)

	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Update will update the provided user with all of the data in the provided user object from the validation layer
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete validate the ID first then pass it to the next in chain
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// Delete will delete a user from the db
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close the db connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// AutoMigrate will attempt to automatically migrate the user table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset tears and rebuilds the user db
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

type userValFn func(*User) error

// this is a lot like the functional programing I am used to in JavaScript
// it is like pipe or flow. (...fns) => (x) => fns.reduce()
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	saltNpepper := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltNpepper, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFn {
	fn := func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	}
	return fn
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}

	if uv.emailRegex.MatchString(user.Email) == false {
		return ErrEmailInvalid
	}

	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)

	// this would not happen during an Update; only on a Create
	if err == ErrNotFound {
		return nil
	}

	// ensure we did not get an err other than ErrNotFound
	if err != nil {
		return err
	}

	if user.ID != existing.ID {
		return ErrEmailTaken
	}

	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	// ignore the validation
	if user.Password == "" {
		return nil
	}

	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}
