package links

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/authentication"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	sessionsStore sessions.Store
	storage       store.Storage
	renderer      *render.Renderer
}

func NewHandler(sessionsStore sessions.Store, storage store.Storage, renderer *render.Renderer) *Handler {
	return &Handler{
		sessionsStore: sessionsStore,
		storage:       storage,
		renderer:      renderer,
	}
}

func (h Handler) SaveLinks(writer http.ResponseWriter, request *http.Request) {
	log.Println("Saving user links")

	session, err := h.sessionsStore.Get(request, authentication.SessionName)
	if err != nil {
		log.Println("Could not get user session.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values[authentication.UserIdSessionVar].(string)
	user, err := h.storage.FindUser(userId)
	if err != nil {
		log.Printf("Can not get user info, %v \n", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Error occurred during reading request body.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	var links []model.Link
	err = json.Unmarshal(body, &links)
	if err != nil {
		log.Println("Error occurred parsing links from user.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Links = links
	user.LastModified = time.Now().UTC()

	err = h.storage.SaveUser(user)

	if err != nil {
		log.Println("Cannot save user links.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("User links '%s' saved successfully.\n", user.Id)

	writer.WriteHeader(http.StatusCreated)
}
