package app

import (
	"encoding/gob"
	"fmt"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	"github.com/jessevdk/go-flags"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"time"
)

const sessionName = "auth-session"

type Auth0Config struct {
	Auth0ClientId     string `long:"auth0 client id" env:"AUTH0_CLIENT_ID" description:"client id of auht0"`
	Auth0ClientSecret string `long:"auth0 client secret" env:"AUTH0_CLIENT_SECRET" description:"client secret of auht0"`
	Auth0Domain       string `long:"auth0 domain" env:"AUTH0_DOMAIN" description:"domain of auth0 client"`
	RedirectUrl       string `long:"redirect callback url" env:"AUTH0_CALLBACK_URL" description:"redirect callback url"`
}

type App struct {
	renderer    *render.Renderer
	session     sessions.Store
	storage     store.Storage
	auth0Config *Auth0Config
}

var (
	revision = "unknown"
	app      App
)

func Init() {
	fmt.Printf("remark42-memory module %s\n", revision)

	session := sessions.NewFilesystemStore("", []byte("secret"))

	storage, err := store.NewStorage(&bolt.Options{Timeout: 1 * time.Second})
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
	app = App{
		renderer:    renderer,
		session:     session,
		storage:     storage,
		auth0Config: auth0Config,
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
	if auth0Config.RedirectUrl == "" {
		panic("Redirect URL is not specified")
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
