package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "urlshort"
	DB_PORT     = "5432"
	DB_USER     = "urlshort"
	DB_PASSWORD = "urlshort"
	DB_NAME     = "urlshort"
)

var DB *sql.DB = InitDB()

func InitDB() *sql.DB {

	log.Println("[postgres] Creating a postgres database connection.")

	initString := fmt.Sprintf(
		`postgresql://%s:%s@localhost:%s/%s?sslmode=disable`,
		DB_USER, DB_PASSWORD, DB_PORT, DB_HOST,
	)

	db, err := sql.Open("postgres", initString)
	if err != nil {
		panic(err)
	}

	log.Println("[postgres] Successfully connected to the postgres database.")
	return db

}

func GetPool() *sql.DB {
	return DB
}
