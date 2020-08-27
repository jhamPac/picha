package controller

import (
	"github.com/jhampac/picha/model"
	"github.com/jhampac/picha/view"
)

// Gallery controller for all related resources
type Gallery struct {
	NewView *view.View
	gs      model.GalleryService
}

// NewGallery instantiates a new controller for the gallery resource
func NewGallery(gs model.GalleryService) *Gallery {
	return &Gallery{
		NewView: view.New("appcontainer", "gallery/new"),
		gs:      gs,
	}
}
