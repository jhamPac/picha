package controller

import "github.com/jhampac/picha/view"

// Static represents all pages that renders static pages with no model bounded to them
type Static struct {
	Home    *view.View
	Contact *view.View
	Error   *view.View
}

// NewStatic instantiates a *Static controller
func NewStatic() *Static {
	return &Static{
		Home:    view.New("appcontainer", "templates/static/home.gohtml"),
		Contact: view.New("appcontainer", "templates/static/contact.gohtml"),
		Error:   view.New("appcontainer", "templates/static/404.gohtml"),
	}
}
