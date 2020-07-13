package controller

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/dlyahov/startuplink-web-go/backend/auth"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func SaveLinks(writer http.ResponseWriter, request *http.Request) {
	log.Println("Saving user links")

	storage := app.GetStorage()
	session, err := app.GetSession(request)
	if err != nil {
		log.Println("Could not get user session.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values[auth.UserIdSessionVar].(string)
	user, err := app.GetStorage().FindUser(userId)
	if err != nil {
		log.Println("Can not get user info")
		log.Printf("Error: %s", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Error occurred during reading request body.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	var links []model.Link
	err = json.Unmarshal(body, &links)
	if err != nil {
		log.Println("Error occurred parsing links from user.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Links = links
	user.LastModified = time.Now().UTC()

	err = storage.SaveUser(user)

	if err != nil {
		log.Println("Cannot save user links.", err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("User links '%s' saved successfully.\n", user.Id)

	writer.WriteHeader(http.StatusCreated)
}
