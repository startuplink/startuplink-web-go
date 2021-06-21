package greeting

import (
	"github.com/dlyahov/startuplink-web-go/backend/authentication"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type Handler struct {
	sessionStore sessions.Store
	renderer     *render.Renderer
}

func NewHandler(sessionStore sessions.Store, renderer *render.Renderer) *Handler {
	return &Handler{
		sessionStore: sessionStore,
		renderer:     renderer,
	}
}

func (h Handler) GreetingHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := h.sessionStore.Get(request, authentication.SessionName)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values[authentication.UserIdSessionVar]; !ok {
		err := h.renderer.RenderTemplate("greeting.html", writer, map[string]interface{}{})
		if err != nil {
			log.Println(err)
			http.Error(writer, "Cannon render page", http.StatusInternalServerError)
		}
	} else {
		log.Println("Proceed user request")
		http.Redirect(writer, request, "/home", http.StatusSeeOther)
	}

}
