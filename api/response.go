package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func newErrorMessage(m string, code int) ErrorMessage {
	return ErrorMessage{
		Message:   m,
		ErrorCode: code,
	}
}

type SuccessfullyCreated struct {
	Message  string `json:"message"`
	ShortURL string `json:"ShortURL"`
}

type Detail struct {
	OriginalURL    string `json:"OriginalURL"`
	ShortURL       string `json:"ShortURL"`
	ExpirationDate string `json:"ExpirationDate"`
	UsedCount      int    `json:"UsedCount"`
}

func (e ErrorMessage) send(w http.ResponseWriter) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.ErrorCode)

	jsonErr, err := json.Marshal(e)

	if err != nil {
		log.Println("[server] An error occurred parsing an error message. ")
		return err
	}

	_, err = w.Write(jsonErr)
	if err != nil {
		log.Println("[server] An error occurred sending an error response. ")
	}

	return err
}

func (e ErrorMessage) Error() string {
	return e.Message
}
