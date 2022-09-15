package redis

import (
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	once sync.Once
)

func Connect() {
	once.Do(func() {
		Client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ENDPOINT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
}

func Disconnect() {
	Client.Close()
}