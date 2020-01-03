package view

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"log"
	"net/http"
)

func ShowLinks(writer http.ResponseWriter, request *http.Request) {
	session, err := app.GetSession(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user := session.Values["user"].(*model.User)

	err = app.GetRenderer().RenderTemplate("main-page.html", writer, user)
	if err != nil {
		log.Fatal(err)
	}
}
