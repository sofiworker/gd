package kit

import (
	"github.com/sourcegraph/conc/pool"
)

func NewPoolWithSize(size int) *pool.Pool {
	p := pool.New().WithMaxGoroutines(size)
	return p
}
