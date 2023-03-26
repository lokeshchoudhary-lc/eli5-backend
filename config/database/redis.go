package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	Client *redis.Client
}

var Redis RedisInstance

var ctx = context.Background()

func ConnectRedis() {
	uri := os.Getenv("REDIS_URI")
	password := os.Getenv("REDIS_PASSWORD")
	username := os.Getenv("REDIS_USERNAME")

	if uri := os.Getenv("REDIS_URI"); uri == "" {
		log.Fatal("You must set your 'REDIS_URI' environmental variable. ")
	}

	if username := os.Getenv("REDIS_USERNAME"); username == "" {
		log.Fatal("You must set your 'REDIS_USERNAME' environmental variable. ")
	}

	if password := os.Getenv("REDIS_PASSWORD"); password == "" {
		log.Fatal("You must set your 'REDIS_PASSWORD' environmental variable. ")
	}

	// opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	// if err != nil {
	// 	panic(err)
	// }

	rdb := redis.NewClient(&redis.Options{
		Addr:     uri,
		Username: username,
		Password: password,
		DB:       0, // use default DB
	})

	Redis = RedisInstance{
		Client: rdb,
	}

	val, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged to redis. ", val)
}
