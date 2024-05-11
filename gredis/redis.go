package gredis

import "github.com/redis/go-redis/v9"

type GRedis struct {
	client *redis.Client
}

func New() interface{} {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, 
	})
	g := &GRedis{
		client: rdb,
	}
	return g
}
