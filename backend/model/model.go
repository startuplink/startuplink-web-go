package model

import (
	"time"
)

type User struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Links        []Link    `json:"links"`
	LastModified time.Time `json:"last_modified"`
}

type Link struct {
	Url    string `json:"url"`
	Pinned bool   `json:"pinned"`
}
