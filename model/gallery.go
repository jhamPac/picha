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

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

// NewGalleryService instantiates a new GalleryService
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}
