package main

import (
	"log"
	"time"
	
	"github.com/dgrijalva/jwt-go"
)

type AuthReq struct {
	// Valid API Key
	ApiKey 		string 		`json:"apiKey"`
	UserId 		string 		`json:"userId"`
}

type AuthResponse struct {
	Token 		string 		`json:"token"`
}

var apiKey = "61529673-6c86-4f54-9bdd-838bf12360a6"
var tokenSecret = "secret"

func getToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    	"userId": id,
    	"dat": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if(err != nil) {
		return "", err
	}

	log.Printf("Generating token: %s", tokenSecret)
	return tokenString, nil
}

func ValidateApiKey(key string) bool {
	if(key == apiKey) {
		return true
	} else {
		return false
	}
}

