package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"urlshort/models"
	"urlshort/utils"

	"github.com/gorilla/mux"
)

type URLHandler struct {
	Port    string
	URLRepo models.URLRepository
}

func (h URLHandler) Home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Welcome! You are about to shorten a URL. \n"))
}

func (h URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := URLShortenRequest{}

	errResp := decodeJSONBody(w, r, &req)
	if errResp != (ErrorMessage{}) {
		errResp.send(w)
		return
	}

	err := req.Validate()
	if err != nil {
		msg := fmt.Sprintf("Not a valid request: %v", err)
		newErrorMessage(msg, http.StatusBadRequest).send(w)
		return
	}

	mapping := models.URLMap{
		OriginalURL:    req.OriginalURL,
		UsedCount:      0,
		ExpirationTime: time.Now().Add(time.Hour),
	}

	shortURL := utils.GenerateShortLink(req.OriginalURL)

	err = h.URLRepo.Save(shortURL, mapping)
	if err != nil {
		msg := err.Error()
		newErrorMessage(msg, http.StatusInternalServerError).send(w)
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

	result, err := h.URLRepo.Get(shortURL)
	if err != nil {
		return
	}

	w.Header().Set("Cache-Control", "no-cache")

	if (models.URLMap{}) == result {
		newErrorMessage("Short URL not found", http.StatusNotFound).send(w)
		return
	}

	if result.ExpirationTime.Before(time.Now()) {
		newErrorMessage("Short URL expired", http.StatusNotFound).send(w)
		return
	}

	result.UsedCount++
	err = h.URLRepo.Update(shortURL, result)
	if err != nil {
		log.Println(
			"[server] An error occurred when trying to update a register:",
			err,
		)
	}

	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)

}

func (h URLHandler) Detail(w http.ResponseWriter, r *http.Request) {

}
