package admin

import (
	json "encoding/json"
	"fmt"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AdminHandler interface {
	GetAllData(writer http.ResponseWriter, request *http.Request)
	DumpDatabase(w http.ResponseWriter, request *http.Request)
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

func (h handler) DumpDatabase(w http.ResponseWriter, request *http.Request) {
	fileName := fmt.Sprintf("%s.db", time.Now().Format("2006_01_02T15_04"))

	err, size := h.storage.DumpDatabase(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Set("Content-Length", strconv.Itoa(size))
}
