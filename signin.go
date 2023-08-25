package main

import (
	"log"
	"net/http"
)

// signin signs a user in. It's a response to an XMLHttpRequest (AJAX request)
// containing the user credentials. It responds with a map[string]string that
// can be converted to JSON by the client. The client expects a boolean
// indicating success or error, and a possible error string.
func signin(w http.ResponseWriter, r *http.Request) {
	// Marshal the Credentials into a credentials struct
	c, err := marshalCredentials(r)
	if err != nil {
		log.Println(err)
		ajaxResponse(w, map[string]string{
			"success": "false",
			"error":   "Invalid Credentials",
		})
		return
	}

	// Get the passwords hash from the database by looking up the users
	// name
	hash, err := rdb.Get(rdx, c.Name).Result()
	if err != nil {
		log.Println(err)
		ajaxResponse(w, map[string]string{
			"success": "false",
			"error":   "User doesn't exist",
		})
		return
	}

	// Check if password matches by hashing it and comparing the hashes
	doesMatch := checkPasswordHash(c.Password, hash)
	if doesMatch {
		newClaims(w, r, c)
		ajaxResponse(w, map[string]string{
			"success": "true",
			"error":   "false",
		})
		return
	}
	ajaxResponse(w, map[string]string{"success": "false", "error": "Bad Password"})
}
