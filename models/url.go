package models

import (
	"time"
)

type URLMap struct {
	OriginalURL    string    `json:"original_url"`
	UsedCount      int       `json:"used_count"`
	ExpirationTime time.Time `json:"expiration_time"`
}
type URLRepository interface {
	Save(string, URLMap) error
	Get(string) (URLMap, error)
	Update(string, URLMap) error
}

var implementation URLRepository

func SetURLRepository(repo URLRepository) {
	implementation = repo
}

func Save(shortURL string, mapping URLMap) error {
	return implementation.Save(shortURL, mapping)
}

func Get(shortURL string) (URLMap, error) {
	return implementation.Get(shortURL)
}

func Update(shortURL string, mapping URLMap) error {
	return implementation.Save(shortURL, mapping)
}
