package authentication

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type AdminAuthentication interface {
	IsAdminAuthenticated(next http.Handler) http.Handler
}

func (h Handler) IsAdminAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			log.Println("Cannot get admin auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authHeaderValues := strings.Split(authHeader[0], " ")
		if authHeaderValues[0] != "Basic" {
			log.Println("It is not basic authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		encodedCredentials, err := base64.StdEncoding.DecodeString(authHeaderValues[1])
		if err != nil {
			log.Println("Cannot get admin credentials")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parsedCredentials := strings.Split(string(encodedCredentials), ":")
		if parsedCredentials[0] != "admin" || parsedCredentials[1] != h.adminPassword {
			log.Println("Admin credentials are incorrect")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
