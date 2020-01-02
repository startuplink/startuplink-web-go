package main

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type User struct {
	Id    string
	Links []string
}

var templates = template.Must(template.ParseFiles(
	"template/login.html", "template/main-page.html"))

func main() {
	db, err := bolt.Open("app.db", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Cannot open database. ", err)
	}

	user := User{"id1", []string{"https://google.com"}}
	db.Update(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		userJson, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("parsing user object: %s", err)
		}

		err = users.Put([]byte(user.Id), userJson)
		if err != nil {
			return fmt.Errorf("saving user: %s", err)
		}
		return nil
	})

	initServer()
}

func initServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/main", mainHandler)

	log.Println("Listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(writer http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(writer, "main-page.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(writer, "login.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
