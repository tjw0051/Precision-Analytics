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

	log.Fatal(http.ListenAndServe(":8080", router))
	log.Printf("Precision API started...")
}
