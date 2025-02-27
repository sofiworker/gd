package gredis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	globalCtx context.Context
}

type Option struct {
}

type Redis interface {
	Get(key string) (string, error)
	Set(key, val string) error
	HGet(key, field string) (string, error)
	HSet(key string, values ...interface{}) (int, error)
}

func NewRedisClient(option Option) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &RedisClient{
		Client:    rdb,
		globalCtx: context.Background(),
	}
}

func (c *RedisClient) Get(key string) (string, error) {
	strCmd := c.Client.Get(c.globalCtx, key)
	return strCmd.String(), strCmd.Err()
}

func (c *RedisClient) Set(key, val string) error {
	statusCmd := c.Client.Set(c.globalCtx, key, val, 0)
	return statusCmd.Err()
}

func (c *RedisClient) HGet(key, field string) (string, error) {
	strCmd := c.Client.HGet(c.globalCtx, key, field)
	return strCmd.String(), strCmd.Err()
}

// HSet
//
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (c *RedisClient) HSet(key string, values ...interface{}) (int, error) {
	intCmd := c.Client.HSet(c.globalCtx, key, values)
	return int(intCmd.Val()), intCmd.Err()
}
