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
		log.Println("Can not get user session")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	value := session.Values["user"]
	var user *model.User
	var ok bool
	if user, ok = value.(*model.User); !ok {
		log.Println("Can not get user from session")
		http.Error(writer, "Can not get user from session", http.StatusInternalServerError)
		return
	}

	err = app.GetRenderer().RenderTemplate("main-page.html", writer, user)
	if err != nil {
		log.Println(err)
	}
}
