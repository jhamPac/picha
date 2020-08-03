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

// ByID queries and returns a user by the id provided
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil

	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound

	default:
		return nil, err
	}
}

// Close the db connection
func (us *UserService) Close() error {
	return us.db.Close()
}
