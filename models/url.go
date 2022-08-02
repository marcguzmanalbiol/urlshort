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
