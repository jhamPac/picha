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
	EditView *view.View
	gs       model.GalleryService
	r        *mux.Router
}

// NewGallery instantiates a new controller for the gallery resource
func NewGallery(gs model.GalleryService, r *mux.Router) *Gallery {
	return &Gallery{
		NewView:  view.New("appcontainer", "gallery/new"),
		ShowView: view.New("appcontainer", "gallery/show"),
		EditView: view.New("appcontainer", "gallery/edit"),
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

	url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Show will display a gallery that matches the provided ID
func (g *Gallery) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd view.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

// Edit a users gallery
func (g *Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permissions to edit this gallery", http.StatusForbidden)
		return
	}
	var vd view.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

// Update a gallery resource: POST /gallery/:id/update
func (g *Gallery) Update(w http.ResponseWriter, r *http.Request) {
	// retrieve the gallery by ID and the user from the context
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// parse form from the edit POST call
	var vd view.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(&form, r); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	// update the gallery
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &view.Alert{
			Level:   view.AlertLvlSuccess,
			Message: "Gallery successfully updated!",
		}
	}
	g.EditView.Render(w, vd)
}

// Delete a gallery resource: POST /gallery/:id/delete
func (g *Gallery) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd view.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, vd)
	}

	fmt.Fprintln(w, "successfully deleted!")
}

func (g *Gallery) galleryByID(w http.ResponseWriter, r *http.Request) (*model.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case model.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Uh oh! something went wrong", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
