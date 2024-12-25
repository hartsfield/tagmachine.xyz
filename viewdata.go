package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// credentials are user credentials and are used in the HTML templates and also
// by handlers that do authorized requests
type credentials struct {
	Name       string   `json:"username"`
	Password   string   `json:"password"`
	IsLoggedIn bool     `json:"isLoggedIn"`
	Posts      []string `json:"posts"`
	Score      uint     `json:"score"`
	jwt.StandardClaims
}

// post is the structure of a user post. Posts are created by users and stored
// in redis.
type post struct {
	Title  string `json:"title" redis:"title"`
	Id     string `json:"id" redis:"id"`
	Author string `json:"author,name" redis:"author"`
	// timestamp
	TS time.Time `json:"ts" redis:"ts"`
	// formatted time stamp
	FTS      string `json:"fts" redis:"fts"`
	BodyText string `json:"bodytext" redis:"bodytext"`
	// TODO: implement nonce
	Nonce      string  `json:"nonce" redis:"nonce"`
	Children   []*post `json:"children" redis:"children"`
	ChildCount int     `json:"childCount" redis:"childCount"`
	Parent     string  `json:"parent" redis:"parent"`
	// used for pagification
	PostCount string `json:"postCount" redis:"postCount"`
	Media     string `json:"media" redis:"media"`
	MediaType string `json:"mediaType" redis:"mediaType"`
}

// viewData represents the root model used to dynamically update the app
type viewData struct {
	ViewType    string `json:"viewType" redis:"viewType"`
	PageTitle   string
	CompanyName string
	Stream      []*post
	Nonce       string
	Order       string `json:"order" redis:"order"`
	UserData    *credentials
}
