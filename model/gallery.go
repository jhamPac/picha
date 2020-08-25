package model

import "github.com/jinzhu/gorm"

// Gallery contains images to view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
