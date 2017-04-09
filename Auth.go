package main

import (
	"fmt"
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

type CustomClaims struct {
	UserId 		string 		`json:"userId"`
	jwt.StandardClaims
}

var apiKey = "61529673-6c86-4f54-9bdd-838bf12360a6"
var tokenSecret = "secret"
var tokenIssuer = "precision"

func getToken(id string) (string, error) {
	log.Printf("Token for: %s", id)
	// Create the Claims
	claims := CustomClaims{
	    id,
	    jwt.StandardClaims{
	        Issuer:    tokenIssuer,
	        IssuedAt: time.Now().Unix(),
	        ExpiresAt: time.Now().Unix() + 3600, // 1hr = 3600
	        NotBefore: time.Now().Unix(),
	    },
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if(err != nil) {
		return "", err
	}

	log.Printf("Generating token: %s", tokenString)
	return tokenString, nil
}

func ValidateApiKey(key string) bool {
	if(key == apiKey) {
		return true
	} else {
		return false
	}
}

//TODO: Implement
func ValidateToken(tokenString string) (bool, error) {
	log.Printf("Validating token: %s", tokenString)

	//at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			secret := []byte(tokenSecret)
		    // Don't forget to validate the alg is what you expect:
		    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		        return false, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		    }
		    return secret, nil
		})

		log.Printf("Finished parse...")
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			log.Printf("Token OK")
			return false, fmt.Errorf("Invalid issuer on token claims")
		} else {
			return false, err
		}
	//})
}

