package admin

import (
	json "encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"log"
	"net/http"
)

type AdminHandler interface {
	GetAllData(writer http.ResponseWriter, request *http.Request)
}

type handler struct {
	storage store.Storage
}

func NewHandler(storage store.Storage) AdminHandler {
	return &handler{storage: storage}
}

func (h handler) GetAllData(writer http.ResponseWriter, request *http.Request) {
	users, err := h.storage.GetAllUsers()
	if err != nil {
		log.Printf("failed get all users. %v\n", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		log.Printf("Cannot parse users to json. %v\n", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		log.Printf("Cannot write users to output. %v\n", err)
	}
}
