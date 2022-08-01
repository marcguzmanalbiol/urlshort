package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
	DB_USER     = "urlshort"
	DB_PASSWORD = "urlshort"
	DB_NAME     = "urlshort"
)

const (
	REDIS_HOST     = "localhost"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
)

var postgresDB *sql.DB = InitPostgres()
var redisDB *redis.Client = InitRedis()

func InitPostgres() *sql.DB {

	log.Println("[postgres] Creating a postgres database connection.")

	initString := fmt.Sprintf(
		`postgresql://%s:%s@%s:%s/%s?sslmode=disable`,
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME,
	)

	db, err := sql.Open("postgres", initString)
	if err != nil {
		panic(err)
	}

	log.Println("[postgres] Successfully connected to the postgres database.")
	return db

}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("[redis] An error has occurred when trying to connect to Redis database: %v", err)
	}

	log.Println("[redis] Successfully connected to the redis database.")

	return client
}

func GetSQLPool() *sql.DB {
	return postgresDB
}

func GetRedisPool() *redis.Client {
	return redisDB
}
