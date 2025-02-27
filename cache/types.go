package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, time.Time, error)
	Put(ctx context.Context, key string, val interface{}, d time.Duration) error
	Delete(ctx context.Context, key string) error
	String() string
}

type Options struct {
}
