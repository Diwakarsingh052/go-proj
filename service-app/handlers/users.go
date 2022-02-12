package handlers

import (
	"context"
	"fmt"
	"net/http"
	"service-app/auth"
	"service-app/data/user"
	"service-app/web"
)


type userHandlers struct {
	*user.DbService
	auth *auth.Auth
}

func (h *userHandlers) SignUp(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValue).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}
	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("%w", err)
	}

	usr, err := h.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("user: %+v %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)

}

func (h *userHandlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValue).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := web.Decode(r, &login); err != nil {
		return fmt.Errorf("%w", err)
	}
	claims, err := h.Authenticate(ctx, v.Now, login.Email, login.Password)

	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticating %w", err)
		}
	}
	var tkn struct {
		Token string `json:"token"`
	}

	tkn.Token, err = h.auth.GenerateToken(claims)

	if err != nil {
		return fmt.Errorf("generating token %w", err)
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)

}