package web

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"syscall"
	"time"
)

type ctxKey int

const KeyValue ctxKey = 1

type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}
type App struct {
	*mux.Router // field name is router
	shutdown    chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		Router:   mux.NewRouter(),
		shutdown: shutdown,
	}
}

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
type sum func(int, int) int

// HandleFunc is custom implementation of handlers.
func (a *App) HandleFunc(method string, path string, handler HandlerFunc) {

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValue, &v)
		err := handler(ctx, w, r) // wrapping the handler inside the h and later we pass h to gorilla router
		if err != nil {
			a.SignalShutdown()
		}
	}

	a.Router.HandleFunc(path, h).Methods(method)

}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
