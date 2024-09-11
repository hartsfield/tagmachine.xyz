package main

import (
	"log"
	"net/http"
	"time"
)

// logout logs the user out by overwriting the token. It must first validate
// the existing token to get the username to overwrite the old token in the
// database
func logout(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Println(err)
	}

	c, err := parseToken(token.Value)
	if err != nil {
		log.Println(err)
	}
	rdb.Set(rdx, c.Name+":token", "loggedout", 0)

	expire := time.Now()
	cookie := http.Cookie{
		Name:    "token",
		Value:   "loggedout",
		Path:    "/",
		Expires: expire,
		MaxAge:  0,
	}
	http.SetCookie(w, &cookie)

	ajaxResponse(w, map[string]string{"error": "false", "success": "true"})
}
