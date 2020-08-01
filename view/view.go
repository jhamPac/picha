package view

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	// LayoutDir path to template layouts
	LayoutDir string = "templates/layouts/"

	// TemplateExt is the extension for templates
	TemplateExt string = ".gohtml"
)

// New instantiates a *View type and returns it
func New(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// View represents a view created by combining n amount of templates
type View struct {
	Template *template.Template
	Layout   string
}

// Render executes a template and writes it to io.Writer
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}
