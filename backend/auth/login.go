package auth

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Could not generate key.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.GetSession(r)
	if err != nil {
		log.Println("Could not retrieve user session.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.AddFlash(state, StateSessionVar)
	err = session.Save(r, w)
	if err != nil {
		log.Println("Could not save user session for request. ", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authenticator, err := newAuthenticator(r.Host)
	if err != nil {
		log.Println("Could not create authenticator object.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Redirect user to login page of OAuth provider")
	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
