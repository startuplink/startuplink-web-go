package app

import (
	"encoding/gob"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
)

type Auth0Config struct {
	Auth0ClientId     string `long:"auth0 client id" env:"AUTH0_CLIENT_ID" description:"client id of auht0"`
	Auth0ClientSecret string `long:"auth0 client secret" env:"AUTH0_CLIENT_SECRET" description:"client secret of auht0"`
	Auth0Domain       string `long:"auth0 domain" env:"AUTH0_DOMAIN" description:"domain of auth0 client"`
}

type App struct {
	renderer    *render.Renderer
	session     sessions.Store
	storage     store.Storage
	auth0Config *Auth0Config
	profile     Profile
}

var (
	//revision = "unknown"
	app App
)

const (
	sessionName         = "auth-session"
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

	auth0Config := &Auth0Config{}
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
	app = App{
		renderer:    renderer,
		session:     session,
		storage:     storage,
		auth0Config: auth0Config,
		profile:     profile,
	}

	gob.Register(map[string]interface{}{})
	gob.Register(&model.User{})
}

func checkConfig(auth0Config *Auth0Config) {
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

func GetSession(request *http.Request) (*sessions.Session, error) {
	return app.session.Get(request, sessionName)
}

func GetStorage() store.Storage {
	return app.storage
}

func GetRenderer() *render.Renderer {
	return app.renderer
}

func GetAuth0Config() *Auth0Config {
	return app.auth0Config
}

func GetProfile() Profile {
	return app.profile
}
