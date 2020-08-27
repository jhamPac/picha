package controller

import (
	"fmt"
	"net/http"

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

// GalleryForm represents the data parsed from the form body
type GalleryForm struct {
	Title string `schema:"title"`
}

// Create parses the form body and create an new gallery
func (g *Gallery) Create(w http.ResponseWriter, r *http.Request) {
	var vd view.Data
	var form GalleryForm

	if err := parseForm(&form, r); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	gallery := model.Gallery{
		Title: form.Title,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
