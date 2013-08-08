package gsix

import (
	"os"
	"log"
	"html/template"
	"path/filepath"
	"io/ioutil"
)

type TemplateEngine struct {
	BaseDir string
	Templates map[string] *template.Template
}


func NewTemplateEngine(basedir string) (*TemplateEngine) {
	templateEngine := new(TemplateEngine)
	templateEngine.BaseDir = basedir
	templateEngine.Templates = make(map[string] *template.Template)

	filepath.Walk(basedir, templateEngine.TemplateWalker)

	return templateEngine
}

func (t *TemplateEngine)TemplateWalker(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		return nil
	}

	namespace, _ := filepath.Rel(t.BaseDir, path)

	templates := t.Templates[namespace]
	if templates == nil {
		var files []string
		fileInfos, _ := ioutil.ReadDir(path)
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				files = append(files, path + "/" + fileInfo.Name())
			}
		}

		if len(files) > 0 {
			templates, err = template.ParseFiles(files...)
			if err != nil {
				log.Fatal(err)
			}
			t.Templates[namespace] = templates
		}
	}

	return nil
}

func(t *TemplateEngine) Render (path string, options map[string]string, callback ViewCallback) {
	callback(nil, "<h1>This is the template engine</h1>")
}
