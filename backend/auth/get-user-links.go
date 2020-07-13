package auth

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"log"
	"net/http"
)

func GetUserLinks(w http.ResponseWriter, r *http.Request) {
	session, err := app.GetSession(r)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values[UserIdSessionVar].(string)
	user, err := app.GetStorage().FindUser(userId)
	if err != nil {
		log.Println("Can not get user info")
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	links, err := json.Marshal(user.Links)
	if err != nil {
		log.Println("Cannot parse JSON from user object: ", user.Links)
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(links)
	if err != nil {
		log.Println("Can write response")
		log.Printf("Error: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
