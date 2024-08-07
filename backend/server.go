package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	authKey     = "32-byte-long-auth-key"
	defaultPort = "8080"
)

func StartServer() {
	r := mux.NewRouter()

	authHandler := app.GetAuthenticationHandler()
	adminAuthHandler := app.GetAdminAuthenticationHandler()
	linksHanlder := app.GetLinksHandler()
	greetingHandler := app.GetGreetingHandler()

	r.HandleFunc("/login", authHandler.LoginHandler)
	r.HandleFunc("/logout", authHandler.LogoutHandler)
	r.HandleFunc("/callback", authHandler.CallbackHandler)

	r.Handle("/get-links", authHandler.AuthenticationMiddleware(linksHanlder.GetUserLinks))

	r.HandleFunc("/", greetingHandler.GreetingHandler)

	r.Handle("/home", authHandler.AuthenticationMiddleware(linksHanlder.ShowLinks))

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))),
	)

	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/ping", healthCheck)

	http.Handle("/", r)
	log.Print("Server listening on http://localhost:8080/")

	var csrfMiddleware mux.MiddlewareFunc

	if app.GetProfile() == app.LOCAL {
		csrfMiddleware = csrf.Protect(
			[]byte(authKey),
			csrf.Secure(false),
		)
	} else {
		csrfMiddleware = csrf.Protect(
			[]byte(authKey),
		)
	}
	api := r.PathPrefix("/api").Subrouter()
	api.Use(csrfMiddleware)

	api.Handle("/save", authHandler.AuthenticationMiddleware(linksHanlder.SaveLinks)).Methods("POST")

	adminHandler := app.GetAdminHandler()
	adminApi := r.PathPrefix("/admin/api").Subrouter()
	adminApi.Use(adminAuthHandler.IsAdminAuthenticated)
	adminApi.HandleFunc("/all-data", adminHandler.GetAllData)
	adminApi.HandleFunc("/export-db", adminHandler.DumpDatabase).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           csrfMiddleware(r),
	}
	log.Println("Start app on port: ", port)
	log.Fatal(server.ListenAndServe())
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Error occurred during healthcheck handling request")
	}
}
