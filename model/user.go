package model

import (
	"errors"

	"github.com/jinzhu/gorm"

	// driver for postgres gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User represents our customers
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

var (
	// ErrNotFound is returned when a resource cannot be found
	ErrNotFound = errors.New("model: resource not found")
)

// UserService is the DB abstraction layer
type UserService struct {
	db *gorm.DB
}

// NewUserService instantiates a new service with the provided connection information
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

// Create a user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// ByID queries and returns a user by the id provided
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail queries and returns a user by the email provided
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Close the db connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset tears and rebuilds the user db
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
