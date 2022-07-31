package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"urlshort/models"
	"urlshort/utils"

	"github.com/gorilla/mux"
)

type URLHandler struct {
	Port string
	DB   models.DBRepo
}

func (h URLHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome! You are about to shorten a URL. \n"))
}

func (h URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := URLShortenRequest{}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// TODO: JSON Handler should not be here!

	err := dec.Decode(&req)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):

			msg := fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset,
			)

			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):

			msg := "Request body contains badly type-formed JSON"
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset,
			)
			http.Error(w, msg, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf(
				"Request body contains unknown field %s",
				fieldName,
			)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			msg := "Request  body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		default:
			log.Printf(
				"[server] There was an non expected error trying to parse the request body: %v",
				err.Error(),
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}

		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object. "
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if err != nil {
		msg := fmt.Sprintf("Not a valid request: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: Format Errors

	mapping := models.MapURL{
		OriginalURL:    req.OriginalURL,
		UsedCount:      0,
		ExpirationTime: time.Now().Add(time.Hour),
	}

	shortURL := utils.GenerateShortLink(req.OriginalURL)

	err = h.DB.Save(shortURL, mapping)
	if err != nil {
		msg := err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp := SuccessfullyCreated{
		Message:  "Short URL created successfully",
		ShortURL: fmt.Sprintf("http://localhost:%v/%s", h.Port, shortURL),
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[server] Unexpected error in the server: %v", err)
		return
	}
	w.Write(jsonResponse)

}

func (h URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["short_url"]

	result, err := h.DB.Get(shortURL)
	if err != nil {
		return
	}

	if result.ExpirationTime.Before(time.Now()) {
		msg := "Short URL expired."

		// TODO: Format messages and errors.

		w.Write([]byte(msg))
		w.WriteHeader(http.StatusGone)
	}

	result.UsedCount++
	// TODO: Update the DB

	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)

}

func (h URLHandler) Detail(w http.ResponseWriter, r *http.Request) {

}
