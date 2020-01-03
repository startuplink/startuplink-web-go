package app

import (
	"encoding/gob"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"github.com/gorilla/sessions"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"time"
)

const sessionName = "auth-session"

type App struct {
	renderer *render.Renderer
	session  *sessions.FilesystemStore
	storage  store.Storage
}

var app App

func Init() {
	session := sessions.NewFilesystemStore("", []byte("secret"))

	storage, err := store.NewStorage(&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("couldn't open storage")
		return
	}
	renderer := render.NewRenderer()
	app = App{
		renderer: renderer,
		session:  session,
		storage:  storage,
	}

	gob.Register(map[string]interface{}{})
}

func GetSession(request *http.Request) (*sessions.Session, error) {
	return app.session.Get(request, sessionName)
}

func GetStorage() *store.Storage {
	return &app.storage
}

func GetRenderer() *render.Renderer {
	return app.renderer
}
