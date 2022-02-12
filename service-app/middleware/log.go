package middleware

import (
	"context"
	"log"
	"net/http"
	"service-app/auth"
	"service-app/web"
	"time"
)

type Mid struct {
	Log *log.Logger
	A   *auth.Auth
}

func (m *Mid) Logger(next web.HandlerFunc) web.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		v, ok := ctx.Value(web.KeyValue).(*web.Values)
		if !ok {
			return web.NewShutdownError("web values missing from the context")
		}
		m.Log.Printf("%s : started   : %s %s ",
			v.TraceID,
			r.Method, r.URL.Path)
		err := next(ctx, w, r)

		m.Log.Printf("%s : completed : %s %s ->  (%d)  (%s)",
			v.TraceID,
			r.Method, r.URL.Path,
			v.StatusCode, time.Since(v.Now),
		)
		return err
	}
}
