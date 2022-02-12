package handlers

import (
	"log"
	"net/http"
	"os"
	"service-app/auth"
	"service-app/data/user"
	"service-app/middleware"
	"service-app/web"
)

func API(shutdown chan os.Signal, log *log.Logger, a *auth.Auth, udb *user.DbService) http.Handler {
	//app := mux.NewRouter()
	//http.HandlerFunc()
	//app.HandleFunc("/ready", check)
	app := web.NewApp(shutdown)
	m := middleware.Mid{
		Log: log,
		A:   a,
	}
	uh := userHandlers{
		DbService: udb,
		auth:      a,
	}

	app.HandleFunc(http.MethodGet, "/ready", m.Logger(m.Error(m.Panic(m.Authenticate(m.HasRole(check, []string{"ADMIN"}))))))
	app.HandleFunc(http.MethodPost, "/create", m.Logger(m.Error(m.Panic(uh.SignUp))))
	app.HandleFunc(http.MethodPost, "/login", m.Logger(m.Error(m.Panic(uh.Login))))
	return app
}
