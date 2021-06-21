package app

import (
	"encoding/gob"
	"github.com/dlyahov/startuplink-web-go/backend/authentication"
	"github.com/dlyahov/startuplink-web-go/backend/greeting"
	"github.com/dlyahov/startuplink-web-go/backend/links"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/oauth"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type App struct {
	renderer *render.Renderer
	session  sessions.Store
	storage  store.Storage

	authenticationHandler      authentication.Autenticator
	adminAuthenticationHandler authentication.AdminAuthentication
	greetingHandler            *greeting.Handler
	linksHandler               *links.Handler
	auth0Config                *oauth.Auth0Config
	profile                    Profile
}

var app App

const (
	cookieAuthKeyName   = "COOKIE_AUTH_KEY"
	cookieSecretKeyName = "COOKIE_SECRET_KEY"
)

func Init() {
	session := sessions.NewCookieStore(
		[]byte(os.Getenv(cookieAuthKeyName)),
		[]byte(os.Getenv(cookieSecretKeyName)),
	)
	session.Options = &sessions.Options{
		HttpOnly: true,
	}

	storage, err := store.NewStorage()
	if err != nil {
		log.Fatal("couldn't open storage")
	}
	renderer := render.NewRenderer()

	auth0Config := &oauth.Auth0Config{}
	_, err = flags.Parse(auth0Config)
	if err != nil {
		log.Fatal("couldn't parse authentication config. ", err.Error())
	}
	checkConfig(auth0Config)

	profileName := os.Getenv("PROFILE")
	profile, err := getProfile(profileName)
	if err != nil {
		log.Fatal("couldn't get profile. ", err.Error())
	}

	log.Println("profile is active: ", profile)

	host := os.Getenv("HOST")
	auth0Client, err := oauth.NewAuth0Client(host, auth0Config)
	if err != nil {
		log.Fatal("can't connect to authenticator. ", err)
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("Admin password is not specified!")
	} else {
		log.Println("Admin password is set")
	}

	authenticationHandler := authentication.NewHandler(session, storage, auth0Client, adminPassword)
	greetingHandler := greeting.NewHandler(session, renderer)
	linkHandler := links.NewHandler(session, storage, renderer)

	app = App{
		renderer:                   renderer,
		session:                    session,
		storage:                    storage,
		auth0Config:                auth0Config,
		profile:                    profile,
		authenticationHandler:      authenticationHandler,
		adminAuthenticationHandler: authenticationHandler,
		greetingHandler:            greetingHandler,
		linksHandler:               linkHandler,
	}

	gob.Register(map[string]interface{}{})
	gob.Register(&model.User{})
}

func checkConfig(auth0Config *oauth.Auth0Config) {
	if auth0Config.Auth0ClientId == "" {
		panic("Client id is not specified")
	}
	if auth0Config.Auth0ClientSecret == "" {
		panic("Client secret is not specified")
	}
	if auth0Config.Auth0Domain == "" {
		panic("Domain is not specified")
	}
}

func GetRenderer() *render.Renderer {
	return app.renderer
}

func GetProfile() Profile {
	return app.profile
}

func GetAuthenticationHandler() authentication.Autenticator {
	return app.authenticationHandler
}

func GetAdminAuthenticationHandler() authentication.AdminAuthentication {
	return app.adminAuthenticationHandler
}

func GetGreetingHandler() *greeting.Handler {
	return app.greetingHandler
}

func GetLinksHandler() *links.Handler {
	return app.linksHandler
}
