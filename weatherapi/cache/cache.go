package cache

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func StartRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

}

func GetRedis() *redis.Client {
	return rdb
}

func CloseRedis() {
	if err := rdb.Close(); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}
}
