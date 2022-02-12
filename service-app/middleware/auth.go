package middleware

import (
	"context"
	"errors"
	"net/http"
	"service-app/auth"
	"service-app/web"
	"strings"
)

// token ""
// bearer "" // Bearer

var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

func (m *Mid) Authenticate(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		//var a auth.Auth // here a will have 0 values set
		// Parse the authorization header. Expected header is of
		// the format `Bearer <token>`.
		parts := strings.Split(r.Header.Get("Authorization"), " ")

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			return web.NewRequestError(err, http.StatusUnauthorized)
		}
		claims, err := m.A.ValidateToken(parts[1])
		if err != nil {
			return web.NewRequestError(err, http.StatusUnauthorized)
		}
		// Add claims to the context so they can be retrieved later.
		ctx = context.WithValue(ctx, auth.Key, claims)

		return next(ctx, w, r)
	}
}

func (m *Mid) HasRole(next web.HandlerFunc, roles []string) web.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, ok := ctx.Value(auth.Key).(auth.Claims)

		if !ok {
			return errors.New("claims missing from context: HasRole called without/before Authenticate")
		}

		if !claims.HasRole(roles...) {
			return ErrForbidden
		}
		return next(ctx, w, r)
	}

}
