package templates

import (
	"html/template"
	"io"
	"os"
)

type Templates struct {
	Templates map[string]*template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.Templates[name].Execute(w, data)
}

func NewTemplates() *Templates {
	tmplFilesDir := "pkg/delivery/rest/views/"
	tmplFiles, err := os.ReadDir(tmplFilesDir)
	if err != nil {
		panic(err)
	}

	parsedTemplates := make(map[string]*template.Template, len(tmplFiles))
	for _, v := range tmplFiles {
		parsedTemplates[v.Name()] = template.Must(template.ParseFiles(tmplFilesDir+"base.tmpl", tmplFilesDir+v.Name()))
	}

	return &Templates{
		Templates: parsedTemplates,
	}
}
