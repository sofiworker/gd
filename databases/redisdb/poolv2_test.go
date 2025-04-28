package redisdb

import (
	"github.com/chuck1024/gd"
	"testing"
)

func TestClient(t *testing.T) {
	c := &RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	}

	o := &RedisPoolClientV2{
		RedisConfig: c,
	}
	err := o.Start()
	if err != nil {
		gd.Debug("err:%s", err)
	}

	o.Set("test", "ok")
	v, err := o.Get("test")
	if err != nil {
		gd.Debug("err:%s", err)
	}
	gd.Debug("%s", v)
}
