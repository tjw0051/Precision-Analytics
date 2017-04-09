package main 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"time"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Precision Analytics API \n\nHere be dragons...")
}

func VerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Precision Analytics API \n\nVersion: " + Version)
}

func LogIndex(w http.ResponseWriter, r *http.Request) {

	// Tell client what data to expect before sending
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK) // Before sending data, tell client the request was OK

	// Kinda like Swift's 'if let'
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

func LogShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logId := vars["logId"]
	fmt.Fprintln(w, "Todo show: ", logId)
	entry := RepoFindEntry(logId)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		panic(err)
	}
}

func LogEntry(w http.ResponseWriter, r *http.Request) {
	var entry Entry
	// Limit upload size to prevent malicious attacks
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil  { // Read error
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	/*
	if err := json.Unmarshal(body, &entry); err != nil { // error unserializing
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	*/

	if err := json.Unmarshal([]byte(body), &entry); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}


	entry.Date = time.Now()
	t := RepoCreateEntry(entry)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func ReqAuth(w http.ResponseWriter, r *http.Request) {
	var authReq AuthReq
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil  { // Read error
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(body), &authReq); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	// Validate API Key
	if(ValidateApiKey(authReq.ApiKey) == false) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(401)
		panic("Invalid API Key")
	}
	// Generate Token
	var response AuthResponse
	response.Token, err = getToken(authReq.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(401)
		panic(err)
	}
	// Send Response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}	

}
/*	ERROR: Initialization loop (routes-handlers-routes)
func Help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Current API routes for latest version(%s): ", Version)
	var routeCpy := routes
	for _, route := range routeCpy {
		fmt.Fprintln(w, "%s - %s - %s", route.Name, route.Method, route.Pattern)
	}
}
*/



