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
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	/*
	Route{
		"Help",
		"GET",
		"/help",
		Help,
	}, */
	Route{
		"VerIndex",
		"GET",
		"/" + Version,
		VerIndex,
	},
	Route{
		"LogIndex",
		"GET",
		"/" + Version + "/get",
		LogIndex,
	},
	Route{
		"LogShow",
		"GET",
		"/" + Version + "/get/{logId}",
		LogShow,
	},
	Route{
		"Log",
		"POST",
		"/" + Version + "/log",
		LogEntry,
	},
	Route{
		"ReqAuth",
		"POST",
		"/" + Version + "/auth",
		ReqAuth,
	},
}




