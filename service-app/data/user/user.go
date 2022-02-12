package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"service-app/auth"
	"strconv"
	"time"
)

// ErrAuthenticationFailure occurs when a user attempts to authenticate but
// anything goes wrong.
var ErrAuthenticationFailure = errors.New("authentication failed")

type DbService struct {
	*sql.DB
}

func (db *DbService) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("generating password hash %w", err)
	}

	u := User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users
		(name, email, password_hash, roles, date_created, date_updated)
		VALUES ( $1, $2, $3, $4, $5, $6)
		Returning id`
	var id int
	if err = db.QueryRowContext(ctx, q, u.Name, u.Email, u.PasswordHash, u.Roles, u.DateCreated, u.DateUpdated).Scan(&id); err != nil {
		return User{}, fmt.Errorf("inserting user %w", err)
	}

	u.ID = strconv.Itoa(id)
	return u, nil
}

func (db *DbService) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	const q = `SELECT id,name,email,roles,password_hash FROM users WHERE email = $1`
	var u User

	err := db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email, &u.Roles, &u.PasswordHash)
	if err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, fmt.Errorf("selecting single user %w", err)
	}
	// Compare the provided password with the saved hash. Use the bcrypt
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   u.ID,
			Audience:  jwt.ClaimStrings{"students"},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: u.Roles,
	}
	return claims, nil

}
