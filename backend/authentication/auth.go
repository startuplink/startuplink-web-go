package authentication

import (
	"github.com/dlyahov/startuplink-web-go/backend/oauth"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	RedirectUrlSessionVar = "redirect_url"
	UserIdSessionVar      = "user_id"
	StateSessionVar       = "state"
	SessionName           = "auth-session"
)

type Handler struct {
	sessionsStore sessions.Store
	storage       store.Storage
	auth0Client   oauth.Auth0Client

	adminPassword string
}

func NewHandler(sessionStore sessions.Store, storage store.Storage, auth0Client oauth.Auth0Client, adminPassword string) *Handler {
	return &Handler{
		auth0Client:   auth0Client,
		sessionsStore: sessionStore,
		storage:       storage,
		adminPassword: adminPassword,
	}
}

type Autenticator interface {
	IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
}

func (h Handler) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := h.sessionsStore.Get(r, SessionName)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values[UserIdSessionVar]; !ok {
		log.Println("User is not authenticated. Redirect to login page")

		session.AddFlash(r.URL.Path, RedirectUrlSessionVar)
		err := session.Save(r, w)
		if err != nil {
			log.Println("Can not save user session: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		log.Println("Proceed user request")
		next(w, r)
	}
}
