package model

import (
	"errors"

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

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete
	ErrInvalidID = errors.New("model: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password is provided
	ErrInvalidPassword = errors.New("model: incorrect password provided")
)

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
	hmac hash.HMAC
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

// NewUserService instantiates a new service with the provided connection information
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}

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
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}

// Create runs through the validation and normalization layer first
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember)

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
		uv.bcryptPassword,
		uv.hmacRemember)

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
			return ErrInvalidID
		}
		return nil
	}
	return fn
}
