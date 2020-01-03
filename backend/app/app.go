package app

import (
	"encoding/gob"
	"github.com/dlyahov/startuplink-web-go/backend/render"
	"github.com/gorilla/sessions"
	"net/http"
)

const sessionName = "auth-session"

var (
	Renderer *render.Renderer
	store    *sessions.FilesystemStore
)

func Init() {
	store = sessions.NewFilesystemStore("", []byte("secret"))
	Renderer = render.NewRenderer()
	gob.Register(map[string]interface{}{})
}

func GetSession(request *http.Request) (*sessions.Session, error) {
	return store.Get(request, sessionName)
}
