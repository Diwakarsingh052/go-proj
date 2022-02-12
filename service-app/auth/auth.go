package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// create // admin , dev
// token // user , dev // aman
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func (c Claims) Valid() error {

	if err := c.RegisteredClaims.Valid(); err != nil {
		return fmt.Errorf("validating standard claims %w", err)
	}
	return nil

}

func (c Claims) HasRole(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

type Auth struct {
	privateKey *rsa.PrivateKey // nil
	algorithm  string          // "" // "Rs257"
	parser     *jwt.Parser     // nil
}

func NewAuth(privateKey *rsa.PrivateKey, algorithm string) (*Auth, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, fmt.Errorf("unknown algorithm %v", algorithm)
	}
	//parser := jwt.Parser{ValidMethods: []string{algorithm}}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{algorithm}))
	a := Auth{
		privateKey: privateKey,
		algorithm:  algorithm,
		parser:     parser,
	}

	return &a, nil

}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {

	method := jwt.GetSigningMethod(a.algorithm)
	tkn := jwt.NewWithClaims(method, claims)

	tokenStr, err := tkn.SignedString(a.privateKey)

	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}

	return tokenStr, nil
}

func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return a.privateKey.Public(), nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parsing token %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
