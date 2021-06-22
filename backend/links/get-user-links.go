package links

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/authentication"
	"log"
	"net/http"
)

func (h Handler) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionsStore.Get(r, authentication.SessionName)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values[authentication.UserIdSessionVar].(string)
	user, err := h.storage.FindUser(userId)
	if err != nil {
		log.Println("Can not get user info")
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	links, err := json.Marshal(user.Links)
	if err != nil {
		log.Println("Cannot parse JSON from user object: ", user.Links)
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(links)
	if err != nil {
		log.Println("Can write response")
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
