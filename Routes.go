package main 

import (
	"net/http"
)

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	// Index of the API - Could serve static page from here
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	// Index of the API version
	Route{
		"VerIndex",
		"GET",
		"/" + Version,
		VersionIndex,
	},
	// Log an analytics message
	Route{
		"Log",
		"POST",
		"/" + Version + "/log",
		LogEntry,
	},
	// Request an auth token
	Route{
		"ReqAuth",
		"POST",
		"/" + Version + "/auth",
		ReqAuth,
	},
}




