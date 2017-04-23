package main 

import (
	"log"
	"net/http"
	"time"
	"strings"
	//"encoding/base64"

	"precision-analytics/auth"
	"precision-analytics/data"
)
// TODO: pass route name to handler func for permissions 
// Copy Logger func but for auth
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s\t",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func AuthRoute(inner http.Handler, route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Route does not require Authentication - pass to handler
			if !route.RequiresAuth {
				inner.ServeHTTP(w,r)
				return
			}

			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			log.Printf(r.Header.Get("Authorization"))
			log.Printf("Header Token: " + s[1])
			if len(s) != 2 {
				SendErr(w, data.Errors["authHeaderInvalid"])
				return
			}
			/*
			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				SendErr(w, data.Errors["authHeaderInvalid"])
				return
			} */
			groupName, err := auth.ValidateToken(s[1])
			if err != nil {
				// TODO: More verbose errors for why (expired, etc)
				SendErr(w, data.Errors["tokenInvalid"])
				return
			}
			group, err := data.GetGroup(groupName)
			if (err != nil) {
				SendErr(w, data.Errors["authHeaderInvalid"])
				return
			}
			allGroup, err := data.GetGroup("*")
			if(err == nil) {
				if strings.Contains(allGroup.Perms, route.Name) {// All groups have permission
					inner.ServeHTTP(w,r)
					return
				}
			}
			if strings.Contains(group.Perms, route.Name) ||     // Group has permission
				strings.Contains(group.Perms, "*") {		    // Group has Wildcard (all perms)
				inner.ServeHTTP(w,r)
			} else {
				SendErr(w, data.Errors["authHeaderInvalid"])
			}
		})
}