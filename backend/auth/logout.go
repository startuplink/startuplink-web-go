package auth

import (
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

	// clean session
	cookie := &http.Cookie{
		Name:   "auth-session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
