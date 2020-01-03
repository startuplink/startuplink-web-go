package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dlyahov/startuplink-web-go/backend/auth"
	"github.com/dlyahov/startuplink-web-go/backend/view"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/login", auth.LoginHandler)
	//r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/callback", auth.CallbackHandler)
	r.Handle("/", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(view.ShowLinks)),
	))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/", r)
	log.Print("Server listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
