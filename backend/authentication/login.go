package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Could not generate key.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := h.sessionsStore.Get(r, SessionName)
	if err != nil {
		log.Println("Could not retrieve user session.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values[StateSessionVar] = state

	err = session.Save(r, w)
	if err != nil {
		log.Println("Could not save user session for request. ", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Redirect user to login page of OAuth provider")
	http.Redirect(w, r, h.auth0Client.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
