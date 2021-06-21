package links

import (
	"github.com/dlyahov/startuplink-web-go/backend/authentication"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
)

func (h Handler) ShowLinks(writer http.ResponseWriter, request *http.Request) {
	session, err := h.sessionsStore.Get(request, authentication.SessionName)
	if err != nil {
		log.Println("Can not get user session")
		log.Printf("Error: %s", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values[authentication.UserIdSessionVar].(string)
	user, err := h.storage.FindUser(userId)
	if err != nil {
		log.Println("Can not get user info")
		log.Printf("Error: %s", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.renderer.RenderTemplate("main-page.html", writer, map[string]interface{}{
		"csrfToken": csrf.Token(request),
		"user":      user,
	})
	if err != nil {
		log.Println(err)
	}
}
