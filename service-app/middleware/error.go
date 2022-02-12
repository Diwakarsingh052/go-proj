package middleware

import (
	"context"
	"net/http"
	"service-app/web"
)

func (m *Mid) Error(next web.HandlerFunc) web.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Create the handler that will be attached in the middleware chain.

		// If the context is missing this value, request the service
		// to be shutdown gracefully.

		v, ok := ctx.Value(web.KeyValue).(*web.Values)
		if !ok {
			return web.NewShutdownError("web values missing from the context")
		}

		// Run the handler chain and catch any propagated error.
		if err := next(ctx, w, r); err != nil {
			// Log the error.
			m.Log.Printf("%s : ERROR : %v", v.TraceID, err)
			// Respond to the error.

			if err := web.RespondError(ctx, w, err); err != nil {
				return err
			}

			// If we receive the shutdown err we need to return it
			// back to the base handler to shutdown the service.
			if ok := web.IsShutdown(err); ok {
				return err
			}

		}
		return nil
	}
}
