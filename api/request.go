package api

import (
	"errors"
	"urlshort/utils"
)

type URLShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

func (r URLShortenRequest) Validate() error {
	if r.OriginalURL == "" {
		return errors.New("original URL not provided")
	}

	switch {
	case r.OriginalURL == "":
		return errors.New("original URL not provided")
	case !utils.CheckIfURL(r.OriginalURL):
		return errors.New("original URL provided is not a valid url")
	}

	return nil
}
