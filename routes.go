package main

import "net/http"

// registerRoutes registers the routes with the provided *http.ServeMux
func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", checkAuth(home))
	mux.HandleFunc("/ranked", checkAuth(getByRanked))
	mux.HandleFunc("/chron", checkAuth(getByChron))
	mux.HandleFunc("/post/", checkAuth(viewPost))
	mux.HandleFunc("/submitRoot", checkAuth(submitRoot))
	mux.HandleFunc("/submitReply", checkAuth(submitReply))
	mux.HandleFunc("/login", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/logout", logout)
}
