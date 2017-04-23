package main 

import (
	"net/http"
)

type Route struct {
	Name			string
	Method			string
	Pattern			string
	RequiresAuth	bool
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	// Index of the API - Could serve static page from here
	Route{
		"Index",
		"GET",
		"/",
		false,
		Index,
	},
	// Index of the API version
	Route{
		"VerIndex",
		"GET",
		"/" + Version,
		false,
		VersionIndex,
	},
	// Log an analytics message
	Route{
		"Log",
		"POST",
		"/" + Version + "/log",
		true,
		LogEntry,
	},
	// Request an auth token
	Route{
		"ReqAuth",
		"POST",
		"/" + Version + "/auth",
		false,
		ReqAuth,
	},

	/****	API Key Management	****/

	// List API Keys
	Route{
		"ShowKeys",
		"GET",
		"/" + Version + "/key/get",
		true,
		GetKeys,
	},
	// Create/Modify API group Key
	Route{
		"SetKeys",
		"POST",
		"/" + Version + "/key/set",
		true,
		SetKeys,
	},
	// Create API group Key
	Route{
		"RemoveKeys",
		"POST",
		"/" + Version + "/key/remove",
		true,
		RemoveKeys,
	},

	/****	Group Management	****/

	// List groups
	Route{
		"ShowGroups",
		"GET",
		"/" + Version + "/group/get",
		true,
		GetGroups,
	},
	// Create/Modify group
	Route{
		"SetGroups",
		"POST",
		"/" + Version + "/group/set",
		true,
		SetGroups,
	},
	// Remove group
	Route{
		"RemoveGroups",
		"POST",
		"/" + Version + "/group/remove",
		true,
		RemoveGroups,
	},
}




