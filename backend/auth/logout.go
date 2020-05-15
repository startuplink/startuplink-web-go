package auth

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	"log"
	"net/http"
	"net/url"
	"os"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("AUTH0_DOMAIN")

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		log.Println("Cannot parse url.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		log.Println(w, "Cannot parse url.", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	session, err := app.GetSession(r)
	if err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// clean session
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Println("Can not save user session: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
