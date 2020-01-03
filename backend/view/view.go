package view

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"log"
	"net/http"
)

func ShowLinks(writer http.ResponseWriter, request *http.Request) {
	err := app.Renderer.RenderTemplate("main-page.html", writer, nil)
	if err != nil {
		log.Fatal(err)
	}
}
