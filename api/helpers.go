package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) ErrorMessage {

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):

			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return newErrorMessage(msg, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return newErrorMessage("Request body contains badly type-fomred JSON", http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return newErrorMessage(msg, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return newErrorMessage(msg, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			return newErrorMessage("Request body must not be empty.", http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			return newErrorMessage("Request body must not be larger than 1MB", http.StatusRequestEntityTooLarge)

		default:
			log.Printf(
				"[server] There was an non expected error trying to parse the request body: %v",
				err.Error(),
			)

			return newErrorMessage(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return newErrorMessage("Request body must only contain a single JSON object. ", http.StatusBadRequest)
	}

	return ErrorMessage{}
}
