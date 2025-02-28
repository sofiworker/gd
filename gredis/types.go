package gredis

type Option struct {
}

type Redis interface {
	Get(key string) (string, error)
	Set(key, val string) error
	HGet(key, field string) (string, error)
	HSet(key string, values ...interface{}) (int, error)
}
