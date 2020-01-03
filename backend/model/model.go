package model

import "net/url"

type User struct {
	Id    string
	Name  string
	Links []Link
}

type Link struct {
	Url    url.URL
	Pinned bool
}
