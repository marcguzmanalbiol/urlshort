package server

import (
	"log"
	"net/http"
	"urlshort/api"
	"urlshort/db"

	"github.com/gorilla/mux"
)

func Run(port string) {

	log.Println("[server] Server is running on port", port)
	r := mux.NewRouter()

	sqlURLRepo := db.SQLRepo{
		DB: db.GetSQLPool(),
	}

	cacheRepo := db.CacheURLRepo{
		DB:          sqlURLRepo,
		RedisClient: db.GetRedisPool(),
	}

	urlHandler := api.URLHandler{
		Port:    "8080",
		URLRepo: cacheRepo,
	}

	r.HandleFunc("/", urlHandler.Home).Methods("GET")
	r.HandleFunc("/create", urlHandler.CreateShortURL).Methods("POST")
	r.HandleFunc("/{short_url}", urlHandler.Redirect).Methods("GET")
	r.HandleFunc("/detail", urlHandler.Detail).Methods("GET")

	err := http.ListenAndServe(port, r)

	if err != nil {
		log.Fatalf(
			"An error has occurred when trying to Run the server: %v",
			err.Error(),
		)
	}

}
