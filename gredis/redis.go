package gredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNoAddresses = fmt.Errorf("no addresses")
)

type RedisClient struct {
	*redis.Client
	globalCtx context.Context
}

func WithGlobalCtx(ctx context.Context) OptionFunc {
	return func(o *Option) {
		o.ctx = ctx
	}
}

func WithAddr(addrs ...string) OptionFunc {
	return func(o *Option) {
		o.addrs = addrs
	}
}

func WithAuth(name, password string) OptionFunc {
	return func(o *Option) {
		o.name, o.password = name, password
	}
}

func WithDBNum(idx int) OptionFunc {
	return func(o *Option) {
		o.dbNumber = idx
	}
}

func NewRedisClient(opts ...OptionFunc) (Redis, error) {
	o := &Option{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(o)
	}
	if len(o.addrs) == 0 {
		return nil, ErrNoAddresses
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     o.addrs[0],
		Username: o.name,
		Password: o.password,
		DB:       o.dbNumber,
	})
	return &RedisClient{
		Client:    rdb,
		globalCtx: context.Background(),
	}, nil
}

func (c *RedisClient) Ping() {
	c.Client.Ping(c.globalCtx)
}

func (c *RedisClient) Get(key string) (string, error) {
	strCmd := c.Client.Get(c.globalCtx, key)
	return strCmd.String(), strCmd.Err()
}

func (c *RedisClient) Set(key, val string) error {
	statusCmd := c.Client.Set(c.globalCtx, key, val, 0)
	return statusCmd.Err()
}

func (c *RedisClient) Del(key, val string) error {
	intCmd := c.Client.Del(c.globalCtx, key, val)
	return intCmd.Err()
}

func (c *RedisClient) HDel(key, field string) (string, error) {
	strCmd := c.Client.HDel(c.globalCtx, key, field)
	return strCmd.String(), strCmd.Err()
}

func (c *RedisClient) HExists(key string, field string) (bool, error) {
	boolCmd := c.Client.HExists(c.globalCtx, key, field)
	return boolCmd.Val(), boolCmd.Err()
}

func (c *RedisClient) HGet(key string, field string) (string, error) {
	strCmd := c.Client.HGet(c.globalCtx, key, field)
	return strCmd.Val(), strCmd.Err()
}

func (c *RedisClient) HGetAll(key string) (map[string]string, error) {
	cmd := c.Client.HGetAll(c.globalCtx, key)
	return cmd.Val(), cmd.Err()
}

func (c *RedisClient) HIncrBy(key, field string, incr int) (int, error) {
	intCmd := c.Client.HIncrBy(c.globalCtx, key, field, int64(incr))
	return int(intCmd.Val()), intCmd.Err()
}

func (c *RedisClient) HIncrByFloat(key, field string, incr float64) (int, error) {
	intCmd := c.Client.HIncrByFloat(c.globalCtx, key, field, incr)
	return int(intCmd.Val()), intCmd.Err()
}

func (c *RedisClient) HKeys(key string) ([]string, error) {
	strCmd := c.Client.HKeys(c.globalCtx, key)
	return strCmd.Val(), strCmd.Err()
}

func (c *RedisClient) HLen(key string) (int, error) {
	intCmd := c.Client.HLen(c.globalCtx, key)
	return int(intCmd.Val()), intCmd.Err()
}

func (c *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	sliceCmd := c.Client.HMGet(c.globalCtx, key, fields...)
	return sliceCmd.Val(), sliceCmd.Err()
}

// HSet
//
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (c *RedisClient) HSet(key string, fields ...interface{}) (int, error) {
	intCmd := c.Client.HSet(c.globalCtx, key, fields...)
	return int(intCmd.Val()), intCmd.Err()
}

func (c *RedisClient) HMSet(key string, values ...interface{}) (bool, error) {
	boolCmd := c.Client.HMSet(c.globalCtx, key, values...)
	return boolCmd.Val(), boolCmd.Err()
}

func (c *RedisClient) HSetNX(key, field string, value interface{}) (bool, error) {
	boolCmd := c.Client.HSetNX(c.globalCtx, key, field, value)
	return boolCmd.Val(), boolCmd.Err()
}

func (c *RedisClient) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	scanCmd := c.Client.HScan(c.globalCtx, key, cursor, match, count)
	err := scanCmd.Err()
	if err != nil {
		return nil, 0, err
	}
	page, cursor := scanCmd.Val()
	return page, cursor, nil
}

func (c *RedisClient) HVals(key string) ([]string, error) {
	cmd := c.Client.HVals(c.globalCtx, key)
	return cmd.Val(), cmd.Err()
}

func (c *RedisClient) HRandField(key string, count int) ([]string, error) {
	sliceCmd := c.Client.HRandField(c.globalCtx, key, count)
	return sliceCmd.Val(), sliceCmd.Err()
}

func (c *RedisClient) HRandFieldWithValues(key string, count int) ([]KeyValue, error) {
	kvCmd := c.Client.HRandFieldWithValues(c.globalCtx, key, count)
	if kvCmd.Err() != nil {
		return nil, kvCmd.Err()
	}
	keyValues := kvCmd.Val()
	ret := make([]KeyValue, len(keyValues))
	for i, kv := range keyValues {
		ret[i] = KeyValue{
			Key:   kv.Key,
			Value: kv.Value,
		}
	}
	return ret, nil
}

func (c *RedisClient) LPop(key string) (string, error) {
	strCmd := c.Client.LPop(c.globalCtx, key)
	return strCmd.Val(), strCmd.Err()
}

func (c *RedisClient) Execute(args ...interface{}) (interface{}, error) {
	cmd := c.Client.Do(c.globalCtx, args...)
	return cmd.Result()
}

func (c *RedisClient) Pub(chanName, msg string) (interface{}, error) {
	cmd := c.Client.Publish(c.globalCtx, chanName, msg)
	return cmd.Result()
}

func (c *RedisClient) Sub(msgChan chan<- string, chanName string) {
	pubsub := c.Client.Subscribe(c.globalCtx, chanName)
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		messages := make([]string, 0)
		if msg.Payload != "" {
			messages = append(messages, msg.Payload)
		}
		if len(msg.PayloadSlice) > 0 {
			messages = append(messages, msg.PayloadSlice...)
		}
		for _, m := range messages {
			msgChan <- m
		}
	}
}

func (c *RedisClient) Tx(fn func(Redis) error) []error {
	//cmders, err := c.Client.TxPipelined(c.globalCtx, func(pipeliner redis.Pipeliner) error {
	//	pipeliner
	//	return fn(&pipeliner)
	//})
	//if err != nil {
	//	return []error{err}
	//}
	//ret := make([]error, 0)
	//for _, cmder := range cmders {
	//	err := cmder.Err()
	//	if err != nil {
	//		ret = append(ret, err)
	//	}
	//}
	//return ret
	return nil
}

func (c *RedisClient) LuaScript(content string, keys []string, values []interface{}) (int, error) {
	script := redis.NewScript(content)
	cmd := script.Run(c.globalCtx, c.Client, keys, values...)
	return cmd.Int()
}
