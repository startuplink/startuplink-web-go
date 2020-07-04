package auth

import (
	"context"
	oidc "github.com/coreos/go-oidc"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

const RedirectUrlSessionVar = "redirect_url"

func newAuthenticator(host string) (*Authenticator, error) {
	ctx := context.Background()
	config := app.GetAuth0Config()

	provider, err := oidc.NewProvider(ctx, "https://"+config.Auth0Domain+"/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     config.Auth0ClientId,
		ClientSecret: config.Auth0ClientSecret,
		RedirectURL:  "http://" + host + "/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	log.Println("Callback url: ", conf.RedirectURL)

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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
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
