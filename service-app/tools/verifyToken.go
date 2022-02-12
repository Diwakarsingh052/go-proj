package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
)

var tokenStr = `JhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkgcHJvamVjdCIsInN1YiI6IjEyMzQ1NjciLCJleHAiOjE2NDE5OTY3MDksImlhdCI6MTY0MTk5MzcwOSwiUm9sZXMiOlsiQURNSU4iXX0.Hxx64RIa7zcukCK2uwi4YaAHeI65VkhyHWozqz9KjJTA-95S2X9bf4Y-rQMTE-9atO6L3LOh8DeTVMYR1J6xBWAgViw4Za66fXW_x85UiEI5xRTYwoUjH_znlxjcJOs_4Rgrc6jSCIwxH7w1ftvvUlR8FrEZmUYHmGoytHzsQJPMS5AFdBwyLvqk7PGrfhFGk4DSU7jkB1I2_4xPjkwIZks0elLGZ4oAc1fNqE-LTPKif6ki7ZuZO9mS27ctVTa90khHx_XkZougzJRYN6iqkp93hV1bSxBGNalVO0y_SIT51RatRnd6dEh9yLqJQ_aMZVbd1CI3BsbfqBoRhR_m2Q`

func main() {

	type claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
		Roles []string `json:roles`
	}
	PrivatePEM, err := os.ReadFile("private.pem")
	if err != nil {
		log.Fatalln("not able to read pem file")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(PrivatePEM)
	_ = privateKey
	if err != nil {
		log.Fatalln("parsing private key")
	}

	var c claims

	token, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
		return privateKey.Public(), nil
		//return []byte("any key"), nil
	})

	if err != nil {
		fmt.Println("parsing token", err)
		return
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return
	}
	fmt.Println(c)

}
