package authentication

import (
	"context"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"log"
	"net/http"
)

func (h Handler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionsStore.Get(r, SessionName)
	if err != nil {
		log.Println("Could not retrieve user session", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	state := session.Values[StateSessionVar].(string)
	if r.URL.Query().Get("state") != state {
		log.Println("State is not valid from user!")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	session.Values[StateSessionVar] = nil

	token, err := h.auth0Client.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("id_token field in oauth2 token.")
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	idToken, err := h.auth0Client.Verify(context.TODO(), rawIDToken)

	if err != nil {
		log.Println("Failed to verify ID Token: ", err)
		http.Error(w, "Failed to verify ID Token", http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		log.Println("Failed to verify ID Token: " + err.Error())
		http.Error(w, "Failed to verify ID Token", http.StatusInternalServerError)
		return
	}

	sub, ok := profile["sub"]
	if !ok {
		log.Println("Failed to get 'sub' property.")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.storage.FindUser(sub.(string))
	if err == store.ErrUserNotFound {
		userName, ok := profile["name"]
		if !ok {
			log.Println("Cannot get user name from parsed jwt token.")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user = &model.User{
			Id:    sub.(string),
			Name:  userName.(string),
			Links: nil,
		}
		err := h.storage.SaveUser(user)
		if err != nil {
			log.Println("Can not save new user: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	session.Values[UserIdSessionVar] = user.Id
	err = session.Save(r, w)
	if err != nil {
		log.Println("Can not save user session: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("User was found with id %s\n", user.Id)

	redirectUrl := "/home"
	if flashes := session.Flashes(RedirectUrlSessionVar); len(flashes) == 1 {

		// Redirect to the requested page if needed
		log.Println("redirect url found")
		redirectUrl = flashes[0].(string)
	} else {
		log.Println("redirect url not found")
	}

	log.Println("Redirect user after login to: ", redirectUrl)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
