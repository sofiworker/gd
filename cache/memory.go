package cache

import "sync"

type memCache struct {
	opts Options

	sync.RWMutex
}
