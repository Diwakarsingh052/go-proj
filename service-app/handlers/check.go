package handlers

import (
	"context"
	"net/http"
	"service-app/web"
)

func check(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//return errors.New("any error")
	//return web.NewRequestError(errors.New("any err"), http.StatusBadRequest)
	//panic("i need to panic")

	status := struct {
		Status string
	}{
		Status: "ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
	//w.Header().Set("Content-Type", "application/json")
	//return json.NewEncoder(w).Encode(status)
}
