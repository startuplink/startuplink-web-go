package render

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
}

var cachedTemplates = []string{
	"login",
	"main-page",
}

func NewRenderer() *Renderer {
	cwd, _ := os.Getwd()
	for i, templateName := range cachedTemplates {
		cachedTemplates[i] = filepath.Join(cwd, "./template/"+templateName+".html")
	}
	renderer := &Renderer{}

	renderer.templates = template.Must(template.ParseFiles(cachedTemplates...))
	return renderer
}

func (renderer *Renderer) RenderTemplate(viewName string, writer http.ResponseWriter, data interface{}) error {
	return renderer.templates.ExecuteTemplate(writer, viewName, data)
}
