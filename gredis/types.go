package gredis

import (
	"context"
	"time"
)

type Option struct {
	ctx            context.Context
	addrs          []string
	name, password string
	dbNumber       int
}

type OptionFunc func(*Option)

type Redis interface {
	Get(key string) (string, error)
	Set(key, val string) error
	HGet(key, field string) (string, error)
	HSet(key string, values ...interface{}) (int, error)
	Tx(fn func(Redis) error) []error
}

type KeyValue struct {
	Key   string
	Value string
}

type TxPipeline interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (string, error)
	Incr(key string)
}
