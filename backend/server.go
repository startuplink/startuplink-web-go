package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dlyahov/startuplink-web-go/backend/auth"
	"github.com/dlyahov/startuplink-web-go/backend/controller"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const authKey = "32-byte-long-auth-key"

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/callback", auth.CallbackHandler)
	r.Handle("/", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(controller.ShowLinks)),
	))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/favicon.ico", faviconHandler)

	http.Handle("/", r)
	log.Print("Server listening on http://localhost:8080/")

	csrfMiddleware := csrf.Protect(
		[]byte(authKey),

		// todo: remove from prod version
		csrf.Secure(false),
	)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(csrfMiddleware)

	api.Handle("/save", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(controller.SaveLinks)),
	)).Methods("POST")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", csrfMiddleware(r)))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
