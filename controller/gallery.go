package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jhampac/picha/context"
	"github.com/jhampac/picha/model"
	"github.com/jhampac/picha/view"
)

const (
	ShowGallery = "show_gallery"
)

// Gallery controller for all related resources
type Gallery struct {
	NewView  *view.View
	ShowView *view.View
	gs       model.GalleryService
	r        *mux.Router
}

// NewGallery instantiates a new controller for the gallery resource
func NewGallery(gs model.GalleryService, r *mux.Router) *Gallery {
	return &Gallery{
		NewView:  view.New("appcontainer", "gallery/new"),
		ShowView: view.New("appcontainer", "gallery/show"),
		gs:       gs,
		r:        r,
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

	user := context.User(r.Context())

	gallery := model.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}

// Show will display a gallery that matches the provided ID
func (g *Gallery) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case model.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Uh oh! something went wrong", http.StatusInternalServerError)
		}
		return
	}

	var vd view.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}
