package model

import "github.com/jinzhu/gorm"

// User represents our customers
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
