package main

import (
	"github.com/dlyahov/startuplink-web-go/backend/app"
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

func main() {
	_, err := bolt.Open("app.db", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Cannot open database. ", err)
	}
	app.Init()
	StartServer()
}
