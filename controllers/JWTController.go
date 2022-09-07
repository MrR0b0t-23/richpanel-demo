package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	EmailId string `json:"emailId"`
	jwt.StandardClaims
}

func GenerateJwtToken(emailId string) string {

	currentTime := time.Now()
	expirationTime := currentTime.Add(300 * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		EmailId: emailId,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func VerifyJwtToken(tokenString string) (bool, string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, "", err
	}

	// parse the jwt token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, "", fmt.Errorf("can't convert token's claims to standard claims")
	}

	var exp time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		exp = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		exp = time.Unix(v, 0)
	}

	if exp.Before(time.Now()) {
		return false, "", fmt.Errorf("token expired at %v", exp)
	} else {
		return true, claims["emailId"].(string), nil
	}
}
