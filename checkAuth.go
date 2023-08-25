package main

import (
	"context"
	"log"
	"net/http"
)

// checkAuth parses and renews the authentication token, and adds it to the
// context. checkAuth is used as a middleware function for routes that allow or
// require authentication.
func checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a generic user object that not signed in to be used
		// as a placeholder until credentials are verified.
		user := credentials{IsLoggedIn: false}
		// ctx is a user who isn't logged in
		ctx := context.WithValue(r.Context(), ctxkey, user)

		// get the "token" cookie
		token, err := r.Cookie("token")
		if err != nil {
			log.Println(err)
			// anonSignin(w, r)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// parse the "token" cookie, making sure it's valid, and
		// obtaining user credentials if it is
		c, err := parseToken(token.Value)
		if err != nil {
			log.Println(err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// check if "token" cookie matches the token stored in the
		// database
		tkn, err := rdb.Get(ctx, c.Name+":token").Result()
		if err != nil {
			log.Println(err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// if the tokens match we renew the token and mark the user as
		// logged in
		if tkn == token.Value {
			c.IsLoggedIn = true
			ctxx := renewToken(w, r, c)
			next.ServeHTTP(w, r.WithContext(ctxx))
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
