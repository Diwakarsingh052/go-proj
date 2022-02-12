package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"time"
)

func main() {

	PrivatePEM, err := os.ReadFile("private.pem")
	if err != nil {
		log.Fatalln("not able to read pem file")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(PrivatePEM)

	if err != nil {
		log.Fatalln("parsing private key")
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "api project",
			Subject:   "1234567",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(50 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Roles: []string{"ADMIN"},
	}
	//generate the token
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	str, err := tkn.SignedString(privateKey)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(str)
}
