package view

import "html/template"

// New instantiates a *View type and returns it
func New(files ...string) *View {
	files = append(files, "templates/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}

// View represents a view created by combining n amount of templates
type View struct {
	Template *template.Template
}
