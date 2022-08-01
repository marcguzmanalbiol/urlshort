package db

import (
	"database/sql"
	"log"
	"urlshort/models"
)

type SQLRepo struct {
	DB *sql.DB
}

func (repo SQLRepo) Save(shortURL string, mapping models.URLMap) error {
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

func (repo SQLRepo) Get(shortURL string) (models.URLMap, error) {

	var mapping models.URLMap

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
		return models.URLMap{}, err
	}

	return mapping, err

}

func (repo SQLRepo) Update(shortURL string, mapping models.URLMap) error {
	_, err := repo.DB.Exec(
		`
		UPDATE url_mapping
		SET used_count = $1
		WHERE short_url = $2
		`, mapping.UsedCount, shortURL,
	)

	return err
}
