package gsix

import (
	"html/template"
	"bytes"
)

type TemplateEngine struct {
	Templates map[string] *template.Template
}


func NewTemplateEngine() (*TemplateEngine) {
	templateEngine := new(TemplateEngine)
	templateEngine.Templates = make(map[string] *template.Template)
	return templateEngine
}

func(t *TemplateEngine) Render (path string, data interface{}, options map[string]string, callback ViewCallback) {
	var err error

	tmpl, ok := t.Templates[path]
	if !ok {
		tmpl, err= tmpl.ParseFiles(path)
		if err != nil {
			callback(err, "")
			return
		}

		t.Templates[path] = tmpl
	}

	var doc bytes.Buffer 
	if err := tmpl.Execute(&doc, data); err != nil {
			callback(err, "")
			return
	}

	callback(nil, doc.String())
}
