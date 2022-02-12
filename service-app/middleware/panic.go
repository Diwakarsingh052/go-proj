package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"service-app/web"
)

func (m *Mid) Panic(next web.HandlerFunc) web.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

		// If the context is missing this value, request the service
		// to be shutdown gracefully.
		v, ok := ctx.Value(web.KeyValue).(*web.Values)
		if !ok {
			return web.NewShutdownError("web value missing from context")
		}
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("PANIC :%v", r))
				// Log the Go stack trace for this panic'd goroutine.
				log.Printf("%s :\n%s", v.TraceID, debug.Stack())

			}

		}()
		return next(ctx, w, r)

	}
}
