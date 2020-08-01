package controller

import "github.com/jhampac/picha/view"

// Static represents all pages that renders static pages with no model bounded to them
type Static struct {
	HomeView    *view.View
	ContactView *view.View
}

// NewStatic instantiates a *Static controller
func NewStatic() *Static {
	return &Static{
		HomeView:    view.New("appcontainer", "templates/static/home.gohtml"),
		ContactView: view.New("appcontainer", "templates/static/contact.gohtml"),
	}
}
