package main 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"time"
	"log"

	"precision-analytics/auth"
	"precision-analytics/data"
)

// API Index
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Precision Analytics API \n\nHere be dragons...")
}

// Index for API Version
func VersionIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Precision Analytics API \n\nVersion: " + Version)
}

/*	Log message on server (requires token from ReqAuth first)	*/
func LogEntry(w http.ResponseWriter, r *http.Request) {
	var entry data.Entry
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &entry); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}
	// Add date of receival
	entry.Date = time.Now()

	// Validate token & entry
	/*
	if valid, _ := ValidateEntry(entry); !valid {
		SendErr(w, data.Errors["tokenInvalid"])
	}
	*/

	// Submit to DB
	data.SetLog(entry)

	// Should reply if successful?
}
/* Submit API Key to receive an Auth Token */
func ReqAuth(w http.ResponseWriter, r *http.Request) {
	var authReq auth.AuthReq
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &authReq); err != nil {
		SendErr(w, data.Errors["jsonInvalid"])
	}
	// Generate Token
	var response auth.AuthResponse
	var err error
	response.Token, err = auth.GetToken(authReq.UserId, authReq.ApiKey)
	if err != nil {
		SendErr(w, data.Errors["tokenErr"])
	}
	// Send Response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(response)
	CheckErr(err)
}

/*	API Key Management	*/

func GetKeys(w http.ResponseWriter, r *http.Request) {
	var getMsg data.GetMsg
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &getMsg); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}

	// Validate token & entry
	/*
	if valid, _ := auth.ValidateToken(getMsg.Token); !valid {
		SendErr(w, data.Errors["tokenInvalid"])
	}
	*/

	keys := data.GetKeys()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	encErr := json.NewEncoder(w).Encode(keys)
	CheckErr(encErr)
}

func SetKeys(w http.ResponseWriter, r *http.Request) {
	var keys data.Keys
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &keys); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}

	// TODO: Check Token in header has auth for this func

	data.SetKeys(keys)
	
	// TODO: Respond
}

func RemoveKeys(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

/*	Group Management	*/

func GetGroups(w http.ResponseWriter, r *http.Request) {
	var getMsg data.GetMsg
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &getMsg); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}

	// Validate token & entry
	/*
	if valid, _ := auth.ValidateToken(getMsg.Token); !valid {
		SendErr(w, data.Errors["tokenInvalid"])
	}
	*/

	groups := data.GetGroups()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	encErr := json.NewEncoder(w).Encode(groups)
	CheckErr(encErr)
}

func SetGroups(w http.ResponseWriter, r *http.Request) {
	var groups data.Groups
	body := Read(w, r)

	if err := json.Unmarshal([]byte(body), &groups); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}

	// TODO: Check Token in header has auth for this func

	data.SetGroups(groups)
	
	// TODO: Respond
}

func RemoveGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

/* 		Utils 		*/

func Read(w http.ResponseWriter, r *http.Request) []byte {
	// Limit upload size to prevent malicious attacks
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil  { // Read error
		panic(err)
	}
	err = r.Body.Close()
	CheckErr(err)
	return body
}

// TODO: use builtin jwt middleware tools
/*
func ReadAuth(r *http.Request) {
	authKey := r.Header.Get("Authorization")
}
*/

// Send JSON error to w
func SendErr(w http.ResponseWriter, err data.ErrorItem) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(err.Code)
	reply := data.ErrorReply{Errors: []data.ErrorItem{err} }
	encErr := json.NewEncoder(w).Encode(reply)
	CheckErr(encErr)
}


