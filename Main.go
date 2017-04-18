package main

import (
	"log"
	"net/http"
	
	//"github.com/dgrijalva/jwt-go"
)

// Define API version
var Version = "v1"

func main() {
	router := NewRouter()

	log.Printf("Precision API starting...")
	//CreateDB()
	log.Fatal(http.ListenAndServe(":8080", router))	
}
