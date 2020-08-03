package model

import (
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

// Close the db connection
func (us *UserService) Close() error {
	return us.db.Close()
}
