package main 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"time"
	"log"

	"Precision-Analytics/auth"
	"Precision-Analytics/db"
	"Precision-Analytics/data"
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
	// Limit upload size to prevent malicious attacks
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil  { // Read error
		panic(err)
	}
	err = r.Body.Close()
	CheckErr(err)

	if err := json.Unmarshal([]byte(body), &entry); err != nil {
		SendErr(w, data.Errors["msgFormat"])
	}
	// Add date of receival
	entry.Date = time.Now()

	// Validate token & entry
	if valid, _ := ValidateEntry(entry); !valid {
		SendErr(w, data.Errors["tokenInvalid"])
	}

	// Submit to DB
	db.SetLog(entry)

	// Should reply if successful?
}
/* Submit API Key to receive an Auth Token */
func ReqAuth(w http.ResponseWriter, r *http.Request) {
	var authReq auth.AuthReq
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	CheckErr(err)
	err = r.Body.Close()
	CheckErr(err)

	if err := json.Unmarshal([]byte(body), &authReq); err != nil {
		SendErr(w, data.Errors["jsonInvalid"])
	}
	// Validate API Key
	if(auth.ValidateApiKey(authReq.ApiKey) == false) {
		SendErr(w, data.Errors["keyInvalid"])
	}
	// Generate Token
	var response auth.AuthResponse
	response.Token, err = auth.GetToken(authReq.UserId)
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
	log.Printf("Not Implemented!")
}

func SetKeys(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

func RemoveKeys(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

/*	Group Management	*/

func GetGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

func SetGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}

func RemoveGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Implemented!")
}



// Send JSON error to w
func SendErr(w http.ResponseWriter, err data.ErrorItem) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(err.Code)
	reply := data.ErrorReply{Errors: []data.ErrorItem{err} }
	encErr := json.NewEncoder(w).Encode(reply)
	CheckErr(encErr)
}


