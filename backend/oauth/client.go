package oauth

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"os"
)

type Auth0Config struct {
	Auth0ClientId     string `long:"auth0 client id" env:"AUTH0_CLIENT_ID" description:"client id of auth0"`
	Auth0ClientSecret string `long:"auth0 client secret" env:"AUTH0_CLIENT_SECRET" description:"client secret of auth0"`
	Auth0Domain       string `long:"auth0 domain" env:"AUTH0_DOMAIN" description:"domain of auth0 client"`
}

type Auth0Client interface {
	AuthCodeURL(state string) string
	Verify(ctx context.Context, idToken string) (*oidc.IDToken, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type client struct {
	provider   *oidc.Provider
	oidcConfig *oidc.Config
	config     oauth2.Config
	ctx        context.Context
}

func (c client) AuthCodeURL(state string) string {
	return c.config.AuthCodeURL(state)
}

func (c client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

func (c client) Verify(ctx context.Context, idToken string) (*oidc.IDToken, error) {
	return c.provider.Verifier(c.oidcConfig).Verify(ctx, idToken)
}

func NewAuth0Client(host string, authConfig *Auth0Config) (Auth0Client, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://"+authConfig.Auth0Domain+"/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH0_CLIENT_ID"),
	}

	conf := oauth2.Config{
		ClientID:     authConfig.Auth0ClientId,
		ClientSecret: authConfig.Auth0ClientSecret,
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	log.Println("Callback url: ", conf.RedirectURL)

	return &client{
		provider:   provider,
		config:     conf,
		ctx:        ctx,
		oidcConfig: oidcConfig,
	}, nil
}
