package auth

import (
	"fmt"
	"log"
	"time"
	
	"precision-analytics/data"

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
	Group 		string 		`json:"group"`
	jwt.StandardClaims
}

var apiKey = "61529673-6c86-4f54-9bdd-838bf12360a6"
var tokenSecret = "secret"
var tokenIssuer = "precision"
// TODO: Delete cache when new key is added
var keyCache data.Keys

func GetToken(id string, apikey string) (string, error) {
	log.Printf("Token for: %s", id)

	key, err := ValidateApiKey(apikey)
	if err != nil {
		return "", err
	}

	// Create the Claims | TODO: put claims in standard claims
	claims := CustomClaims{
	    id,
	    key.Group,
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

// TODO: check keys in DB - errors for expired, deactive, limit-reached, not found, etc.
func ValidateApiKey(key string) (data.Key, error) {
	keyCache = data.GetKeys()

	for i := 0; i < len(keyCache.Keys); i++ {
		if keyCache.Keys[i].Key == key {
			// Key has been Deactivated
			if !keyCache.Keys[i].Active {
				return data.Key{}, fmt.Errorf("Key has been deactivated.")
			}
			// Key has Expired
			if keyCache.Keys[i].Expires && time.Now().After(keyCache.Keys[i].ExpDate) {
				return data.Key{}, fmt.Errorf("Key has expired.")

			}
			return keyCache.Keys[i], nil
		}
	}
	return data.Key{}, fmt.Errorf("Unknown API Key.")
}
/*	Checks if the route requires Authentication. If it does:
	- Token is extracted from http Header
	- Token is validated
	- Group name is extracted from Token claims
	- Checks if group has permission to access API route

	NOTE: 	Groups with the '*' wildcard permission can access all routes.
			All groups inherit from the wildcard group with name '*'.
				(this group can be modified or removed as necessary)
			The '*' group CANNOT also contain the wildcard '*' permission.
*/

func ValidateToken(tokenString string) (string, error) {
	log.Printf("Validating token: %s", tokenString)

	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
		return claims.Group, nil
		//return false, fmt.Errorf("Invalid issuer on token claims")
	} else {
		return "", fmt.Errorf("Invalid Token")
	}
}

