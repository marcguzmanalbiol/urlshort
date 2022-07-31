package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type MapURL struct {
	OriginalURL    string    `json:"original_url"`
	UsedCount      int       `json:"used_count"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func (s MapURL) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *MapURL) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

type URLDatabase interface {
	Save(shortURL string, mapURL MapURL) error
	Get(shortURL string) (MapURL, error)
	Update(shortURL string, mapURL MapURL) error
}

type DBRepo struct {
	DB *sql.DB
}

func (repo DBRepo) Save(shortURL string, mapping MapURL) error {
	_, err := repo.DB.Exec(
		`
		INSERT INTO url_mapping(short_url, original_url, used_count, expiration_time)
		VALUES ($1, $2, $3, $4)
		`,
		shortURL,
		mapping.OriginalURL,
		mapping.UsedCount,
		mapping.ExpirationTime,
	)

	return err
}

func (repo DBRepo) Get(shortURL string) (MapURL, error) {

	var mapping MapURL

	err := repo.DB.QueryRow(
		`
		SELECT original_url, used_count, expiration_time 
		FROM url_mapping 
		WHERE short_url = $1
		`,
		shortURL,
	).Scan(&mapping.OriginalURL, &mapping.UsedCount, &mapping.ExpirationTime)

	if err != nil {
		log.Println("[postgres] An error occurred when trying to get from the db", err)
		return MapURL{}, err
	}

	return mapping, err

}
