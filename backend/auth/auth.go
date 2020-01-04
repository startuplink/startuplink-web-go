package auth

import (
	"context"
	oidc "github.com/coreos/go-oidc"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func newAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://"+os.Getenv("AUTH0_DOMAIN")+"/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := app.GetSession(r)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		log.Println("User is not authenticated. Redirect to login page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		log.Println("Proceed user request")
		next(w, r)
	}
}
