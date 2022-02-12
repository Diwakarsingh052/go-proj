package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	v, ok := ctx.Value(KeyValue).(*Values)

	if !ok {
		return NewShutdownError("web value missing from the context")
	}
	v.StatusCode = statusCode
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil

}

// RespondError sends an error response back to the client.
//this should not be directly called by handler func. handler func should return errors and let middleware handle that
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	var webErr *Error
	if ok := errors.As(err, &webErr); ok {
		er := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}

func Decode(r *http.Request, val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&val)
	return err
}
