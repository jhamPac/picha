package model

import "github.com/jinzhu/gorm"

// Gallery contains images to view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// GalleryService provides an interface to the Gallery model
type GalleryService interface {
	GalleryDB
}

// GalleryDB is the DB connection for galleries
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return nil
}
