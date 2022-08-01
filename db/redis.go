package db

import (
	"encoding/json"
	"log"
	"time"
	"urlshort/models"

	"github.com/go-redis/redis"
)

const (
	redisTimeout = 6 * time.Hour
)

type CacheURLRepo struct {
	DB          models.URLRepository
	RedisClient *redis.Client
}

func (s CacheURLRepo) Save(shortURL string, mapping models.URLMap) error {
	err := s.DB.Save(shortURL, mapping)
	if err != nil {
		return err
	}

	return nil
}

func (s CacheURLRepo) Get(shortURL string) (models.URLMap, error) {
	var mapping models.URLMap

	result, err := s.RedisClient.Get(shortURL).Result()
	json.Unmarshal([]byte(result), &mapping)

	if err == redis.Nil {
		log.Printf("[redis] Redis cache miss: %v", err)

		mapping, err = s.DB.Get(shortURL)
		if err != nil {
			return mapping, err
		}

		register, err := json.Marshal(mapping)
		if err != nil {
			log.Printf("[redis] There was an error tryin to Marshal the URL mapping: %v", err)
		}

		err = s.RedisClient.Set(shortURL, register, redisTimeout).Err()

		return mapping, err

	} else if err != nil {

		return mapping, err
	}

	log.Println("[redis] Redis cache hit")

	return mapping, nil
}

func (s CacheURLRepo) Update(shortURL string, mapping models.URLMap) error {
	err := s.DB.Update(shortURL, mapping)
	if err != nil {
		log.Println("[redis] Error trying to update a cached register", err)
	}

	register, err := json.Marshal(mapping)
	if err != nil {
		log.Printf("[redis] There was an error tryin to Marshal the URL mapping: %v", err)
	}

	err = s.RedisClient.Set(shortURL, register, 0).Err()

	return err
}
