package auth

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/dlyahov/startuplink-web-go/backend/store"
	"log"
	"net/http"
	"os"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.GetSession(r)
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

	authenticator, err := newAuthenticator(r.Host)
	if err != nil {
		log.Println("Could not create authenticator object", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
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

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH0_CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		log.Println("Failed to verify ID Token: " + err.Error())
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

	var storage = app.GetStorage()
	user, err := storage.FindUser(profile["sub"].(string))

	if err == store.ErrUserNotFound {
		user = &model.User{
			Id:    profile["sub"].(string),
			Name:  profile["name"].(string),
			Links: nil,
		}
		err := storage.SaveUser(user)
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
