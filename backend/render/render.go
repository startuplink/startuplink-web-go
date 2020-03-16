package render

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Renderer struct {
	templates map[string]*template.Template
}

var cachedTemplates = []string{
	"login.html",
	"main-page.html",
}

func NewRenderer() *Renderer {
	cwd, _ := os.Getwd()
	var templatePaths = make([]string, len(cachedTemplates))

	for i, templateName := range cachedTemplates {
		templatePaths[i] = filepath.Join(cwd, "./template/"+templateName)
	}
	renderer := &Renderer{
		templates: make(map[string]*template.Template, len(cachedTemplates)),
	}

	for i, templateName := range cachedTemplates {
		parseFiles, err := template.New(templateName).Funcs(template.FuncMap{
			"inc": func(n int) int {
				return n + 1
			},
		}).ParseFiles(templatePaths[i])

		if err != nil {
			log.Fatal("Cannot parse input templates")
		}

		renderer.templates[templateName] = template.Must(parseFiles, err)
	}

	return renderer
}

func (renderer *Renderer) RenderTemplate(viewName string, writer http.ResponseWriter, data interface{}) error {
	return renderer.templates[viewName].ExecuteTemplate(writer, viewName, data)
}
