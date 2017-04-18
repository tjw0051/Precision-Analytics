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

	/****	API Key Management	****/

	// List API Keys
	Route{
		"ShowKeys",
		"POST",
		"/" + Version + "/key/get",
		GetKeys,
	},
	// Create/Modify API group Key
	Route{
		"SetKey",
		"POST",
		"/" + Version + "/key/set",
		SetKeys,
	},
	// Create API group Key
	Route{
		"RemoveKey",
		"POST",
		"/" + Version + "/key/remove",
		RemoveKeys,
	},

	/****	Group Management	****/

	// List groups
	Route{
		"ShowGroups",
		"POST",
		"/" + Version + "/group/get",
		GetGroups,
	},
	// Create/Modify group
	Route{
		"SetGroups",
		"POST",
		"/" + Version + "/group/set",
		SetGroups,
	},
	// Remove group
	Route{
		"RemoveGroup",
		"POST",
		"/" + Version + "/group/remove",
		RemoveGroups,
	},
}




