package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dlyahov/startuplink-web-go/backend/auth"
	"github.com/dlyahov/startuplink-web-go/backend/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/callback", auth.CallbackHandler)
	r.Handle("/", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(controller.ShowLinks)),
	))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// todo: add CSRF handling
	r.Handle("/save", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(controller.SaveLinks)),
	)).Methods("POST")
	r.HandleFunc("/favicon.ico", faviconHandler)

	http.Handle("/", r)
	log.Print("Server listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
