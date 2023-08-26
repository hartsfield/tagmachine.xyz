package main

import "net/http"

// registerRoutes registers the routes with the provided *http.ServeMuc
func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", checkAuth(home))
	mux.HandleFunc("/ranked", checkAuth(getByRanked))
	mux.HandleFunc("/chron", checkAuth(getByChron))
	mux.HandleFunc("/post/", checkAuth(viewPost))
	mux.HandleFunc("/submitForm", checkAuth(handleForm))
	mux.HandleFunc("/login", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/logout", logout)
}
