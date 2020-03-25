package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"github.com/dlyahov/startuplink-web-go/backend/auth"
	"github.com/dlyahov/startuplink-web-go/backend/controller"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	authKey            = "32-byte-long-auth-key"
	defaultPort        = "8080"
	PKI_VALIDATION_KEY = "PKI_VALIDATION_KEY"
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

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/ping", healthCheck)
	r.HandleFunc("/.well-known/pki-validation", pkiValidationFile)

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

	api.Handle("/save", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(controller.SaveLinks)),
	)).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Println("Start app on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, csrfMiddleware(r)))
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

func pkiValidationFile(w http.ResponseWriter, r *http.Request) {
	validationKey := os.Getenv(PKI_VALIDATION_KEY)
	if validationKey == "" {
		log.Println("Validation key not found!")
		_, err := w.Write([]byte("Validation key not found"))
		if err != nil {
			log.Println("Error occurred when handle endpoint of validation key")
		}
		return
	}

	_, err := w.Write([]byte(validationKey))
	if err != nil {
		log.Println("Error occurred when handle endpoint of validation key")
	}
}
