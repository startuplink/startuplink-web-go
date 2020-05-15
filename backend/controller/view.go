package controller

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
)

func ShowLinks(writer http.ResponseWriter, request *http.Request) {
	session, err := app.GetSession(request)
	if err != nil {
		log.Println("Can not get user session")
		log.Printf("Error: %s", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	value := session.Values["user"].(*model.User)
	user, err := app.GetStorage().FindUser(value.Id)
	if err != nil {
		log.Println("Can not get user info")
		log.Printf("Error: %s", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = app.GetRenderer().RenderTemplate("main-page.html", writer, map[string]interface{}{
		"csrfToken": csrf.Token(request),
		"user":      user,
	})
	if err != nil {
		log.Println(err)
	}
}
