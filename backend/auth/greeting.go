package auth

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"log"
	"net/http"
)

func GreetingHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := app.GetSession(request)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		err := app.GetRenderer().RenderTemplate("greeting.html", writer, map[string]interface{}{})
		if err != nil {
			log.Println(err)
			http.Error(writer, "Cannon render page", http.StatusInternalServerError)
		}
	} else {
		log.Println("Proceed user request")
		http.Redirect(writer, request, "/home", http.StatusSeeOther)
	}

}
