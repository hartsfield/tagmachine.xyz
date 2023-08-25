package main

import (
	"log"
	"net/http"
	"regexp"
)

// signup signs a user up. It's a response to an XMLHttpRequest (AJAX request)
// containing new user credentials. It responds with a map[string]string that
// can be converted to JSON. The client expects a boolean indicating success or
// error, and a possible error string.
func signup(w http.ResponseWriter, r *http.Request) {
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

	// Make sure the username doesn't contain forbidden symbols
	match, err := regexp.MatchString("^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$", c.Name)
	if err != nil {
		log.Println(err)
		ajaxResponse(w, map[string]string{
			"success": "false",
			"error":   "Invalid Username",
		})
		return
	}

	// Make sure the username is longer than 3 characters and shorter than
	// 25, and the password is longer than 7.
	if match && (len(c.Name) < 25) && (len(c.Name) > 3) && (len(c.Password) > 7) {
		// Check if user already exists
		_, err = rdb.Get(rdx, c.Name).Result()
		if err != nil {
			// If username is unique and valid, we attempt to hash
			// the password
			hash, err := hashPassword(c.Password)
			if err != nil {
				log.Println(err)
				ajaxResponse(w, map[string]string{
					"success": "false",
					"error":   "Invalid Password",
				})
				return
			}

			// Add the user the the USERS set in redis. This
			// associates a score with the user that can be
			// incremented or decremented
			_, err = rdb.ZAdd(rdx, "USERS", makeZmem(c.Name)).Result()
			if err != nil {
				log.Println(err)
				ajaxResponse(w, map[string]string{
					"success": "false",
					"error":   "Error ",
				})
				return
			}

			// If the password is hashable, and we were able to add
			// the user to the redis ZSET, we store the hash in the
			// database with the username as the key and the hash
			// as the value thats returned by the key.
			_, err = rdb.Set(rdx, c.Name, hash, 0).Result()
			if err != nil {
				log.Println(err)
				ajaxResponse(w, map[string]string{
					"success": "false",
					"error":   "Error ",
				})
				return
			}

			// Set user token/credentials
			newClaims(w, r, c)

			// success response
			ajaxResponse(w, map[string]string{
				"success": "true",
				"error":   "false",
			})
			return
		}
		ajaxResponse(w, map[string]string{
			"success": "false",
			"error":   "User Exists",
		})
		return
	}
	ajaxResponse(w, map[string]string{
		"success": "false",
		"error":   "Invalid Username",
	})
}
