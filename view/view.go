package view

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	// LayoutDir path to template layouts
	LayoutDir string = "templates/layouts/"

	// TemplateDir path to template directory
	TemplateDir string = "templates/"

	// TemplateExt is the extension for templates
	TemplateExt string = ".gohtml"
)

// View represents a view created by combining n amount of templates
type View struct {
	Template *template.Template
	Layout   string
}

// New instantiates a *View type and returns it
func New(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
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

// Render executes a template and writes it to io.Writer
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath prepends the path to templates to a strings in the slice
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
