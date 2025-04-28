/**
 * Copyright 2019 gd Author. All rights reserved.
 * Author: Chuck1024
 */

package redisdb

import (
	"context"
	"errors"
	"fmt"
	log "github.com/chuck1024/gd/dlog"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
	"net"
	"strings"
	"sync"
	"time"
)

type RedisPoolClientV2 struct {
	RedisConfig   *RedisConfig `inject:"redisConfig" canNil:"true"`
	RedisConf     *ini.File    `inject:"redisConf" canNil:"true"`
	RedisConfPath string       `inject:"redisConfPath" canNil:"true"`
	PoolName      string       `inject:"poolName" canNil:"true"`
	redisRing     *redis.Ring
	startOnce     sync.Once
	closeOnce     sync.Once
}

type RedisClient struct {
	*redis.Ring
}

func (p *RedisPoolClientV2) Start() error {
	var err error
	p.startOnce.Do(func() {
		if p.RedisConfig != nil {
			err = p.newRedisPools(p.RedisConfig)
		} else if p.RedisConf != nil {
			err = p.initRedis(p.RedisConf, p.PoolName)
		} else {
			if p.RedisConfPath == "" {
				p.RedisConfPath = defaultConf
			}

			err = p.initObjForRedisDb(p.RedisConfPath)
		}
	})
	return err
}

func (p *RedisPoolClientV2) Close() {
	p.closeOnce.Do(func() {
		var e error
		if p.redisRing != nil {
			e = p.redisRing.Close()
		}
		if e == nil {
			log.Info("redis pool close ok")
		}
	})
}

func (p *RedisPoolClientV2) newRedisPools(cfg *RedisConfig) error {
	if len(cfg.Addrs) <= 0 {
		return errors.New("servers empty")
	}
	maxActive := cfg.MaxActive
	if maxActive <= 0 {
		maxActive = DefaultMaxActive
	}
	maxIdle := cfg.MaxIdle
	if maxIdle <= 0 {
		maxIdle = DefaultMaxIdle
	}
	idleTimeout := time.Duration(cfg.IdleTimeoutSec) * time.Second
	if idleTimeout <= 0 {
		idleTimeout = time.Duration(DefaultIdleTimeout) * time.Second
	}
	retry := cfg.Retry
	if retry <= 0 {
		retry = DefaultRetryTimes
	}
	connTimeout := time.Duration(cfg.ConnTimeoutMs) * time.Millisecond
	if connTimeout <= 0 {
		connTimeout = DefaultConnTimeout * time.Millisecond
	}
	readTimeout := time.Duration(cfg.ReadTimeoutMs) * time.Millisecond
	if readTimeout <= 0 {
		readTimeout = DefaultReadTimeout * time.Millisecond
	}
	writeTimeout := time.Duration(cfg.WriteTimeoutMs) * time.Millisecond
	if writeTimeout <= 0 {
		writeTimeout = DefaultWriteTimeout * time.Millisecond
	}

	servers := make(map[string]string)
	for idx, addr := range cfg.Addrs {
		servers[fmt.Sprintf("%d", idx)] = addr
	}
	rdb := redis.NewRing(&redis.RingOptions{
		DB:              cfg.DbNumber,
		Username:        cfg.Username,
		Password:        cfg.Password,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		DialTimeout:     connTimeout,
		MaxRetries:      retry,
		MaxIdleConns:    maxIdle,
		PoolSize:        maxActive,
		ConnMaxIdleTime: idleTimeout,
		Addrs:           servers,
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   connTimeout,
				KeepAlive: 30 * time.Second,
			}
			return netDialer.DialContext(ctx, network, addr)
			//if opt.TLSConfig == nil {
			//}
			//return tls.DialWithDialer(netDialer, network, addr, opt.TLSConfig)
		},
	})
	p.redisRing = rdb
	return nil
}

func (p *RedisPoolClientV2) initObjForRedisDb(redisConfPath string) error {
	redisConfRealPath := redisConfPath
	if redisConfRealPath == "" {
		return errors.New("redisConf not set in g_cfg")
	}

	if !strings.HasSuffix(redisConfRealPath, ".ini") {
		return errors.New("redisConf not an ini file")
	}

	redisConf, err := ini.Load(redisConfRealPath)
	if err != nil {
		return err
	}

	if err = p.initRedis(redisConf, p.PoolName); err != nil {
		return err
	}
	return nil
}

func (p *RedisPoolClientV2) initRedis(f *ini.File, pn string) error {
	r := f.Section(fmt.Sprintf("%s.%s", "Redis", pn))
	addr := r.Key("addr").String()
	password := r.Key("password").String()
	username := r.Key("username").String()
	maxActive, _ := r.Key("maxActive").Int()
	maxIdle, _ := r.Key("maxIdle").Int()
	retry, _ := r.Key("retry").Int()
	idleTimeout, _ := r.Key("idleTimeout").Int()
	connTimeout, _ := r.Key("connTimeout").Int64()
	readTimeout, _ := r.Key("readTimeout").Int64()
	writeTimeout, _ := r.Key("writeTimeout").Int64()
	dbNumber, _ := r.Key("dbNumber").Int()

	addrs := strings.Split(addr, ",")
	err := p.newRedisPools(&RedisConfig{
		Addrs:          addrs,
		MaxActive:      maxActive,
		MaxIdle:        maxIdle,
		Retry:          retry,
		IdleTimeoutSec: idleTimeout,
		ConnTimeoutMs:  connTimeout,
		ReadTimeoutMs:  readTimeout,
		WriteTimeoutMs: writeTimeout,
		Password:       password,
		DbNumber:       dbNumber,
		Username:       username,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *RedisPoolClientV2) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	cmd := p.redisRing.Do(context.Background(), commandName, args)
	return cmd.Result()
}

func (p *RedisPoolClientV2) GetClient() *RedisClient {
	return &RedisClient{p.redisRing}
}

func (p *RedisPoolClientV2) Get(key string) (ret string, errRet error) {
	cmd := p.redisRing.Get(context.Background(), key)
	return cmd.Result()
}

func (p *RedisPoolClientV2) Set(key, value string) (err error) {
	cmd := p.redisRing.Set(context.Background(), key, value, 0)
	return cmd.Err()
}

func (p *RedisPoolClientV2) Del(key ...string) (err error) {
	cmd := p.redisRing.Del(context.Background(), key...)
	return cmd.Err()
}

func (p *RedisPoolClientV2) HGetAll(key string) (ret map[string]string, err error) {
	cmd := p.redisRing.HGetAll(context.Background(), key)
	return cmd.Result()
}

func (p *RedisPoolClientV2) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	cmd := p.redisRing.HScan(context.Background(), key, cursor, match, count)
	return cmd.Result()
}

func (p *RedisPoolClientV2) HGet(key string, field string) (string, error) {
	cmd := p.redisRing.HGet(context.Background(), key, field)
	return cmd.Result()
}

func (p *RedisPoolClientV2) HDel(key, field string) (err error) {
	cmd := p.redisRing.HDel(context.Background(), key, field)
	return cmd.Err()
}

func (p *RedisPoolClientV2) HSet(key string, field string, value string) (err error) {
	cmd := p.redisRing.HSet(context.Background(), key, field, value)
	return cmd.Err()
}

func (p *RedisPoolClientV2) HMGet(key string, fields []string) ([]string, error) {
	cmd := p.redisRing.HMGet(context.Background(), key, fields...)
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(fields))
	for i, v := range result {
		ret[i] = fmt.Sprintf("%v", v)
	}
	return ret, nil
}

func (p *RedisPoolClientV2) IncrBy(key string, val int64) (int64, error) {
	cmd := p.redisRing.IncrBy(context.Background(), key, val)
	return cmd.Result()
}

func (p *RedisPoolClientV2) HIncrBy(key, field string, val int64) (int64, error) {
	cmd := p.redisRing.HIncrBy(context.Background(), key, field, val)
	return cmd.Result()
}

func (p *RedisPoolClientV2) Incr(key string) (int, error) {
	cmd := p.redisRing.Incr(context.Background(), key)
	return int(cmd.Val()), cmd.Err()
}

func (p *RedisPoolClientV2) Expire(key string, t time.Duration) (int, error) {
	cmd := p.redisRing.Expire(context.Background(), key, t)
	_, err := cmd.Result()
	return 0, err
}

func (p *RedisPoolClientV2) SetNX(key, value string, expire int) (interface{}, error) {
	cmd := p.redisRing.SetNX(context.Background(), key, value, time.Duration(expire))
	return cmd.Result()
}

func (p *RedisPoolClientV2) MGet(keys []string) (ret []interface{}, errRet error) {
	cmd := p.redisRing.MGet(context.Background(), keys...)
	return cmd.Result()
}

func (p *RedisPoolClientV2) SetEx(key string, value string, expire time.Duration) (err error) {
	cmd := p.redisRing.SetEx(context.Background(), key, value, expire)
	return cmd.Err()
}

func (p *RedisPoolClientV2) SAdd(key string, vals []string) error {
	var args []interface{}
	for _, v := range vals {
		args = append(args, v)
	}
	cmd := p.redisRing.SAdd(context.Background(), key, args...)
	return cmd.Err()
}

func (p *RedisPoolClientV2) ZAdd(key string, score int64, val string) error {
	cmd := p.redisRing.ZAdd(context.Background(), key, redis.Z{Score: float64(score), Member: val})
	return cmd.Err()
}

func (p *RedisPoolClientV2) ZRemByScore(key string, start string, end string) error {
	cmd := p.redisRing.ZRemRangeByScore(context.Background(), key, start, end)
	return cmd.Err()
}

func (p *RedisPoolClientV2) ZRange(key string, start int64, end int64) ([]string, error) {
	cmd := p.redisRing.ZRange(context.Background(), key, start, end)
	return cmd.Result()
}

func (p *RedisPoolClientV2) Exists(key string) (bool, error) {
	cmd := p.redisRing.Exists(context.Background(), key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if cmd.Val() == 0 {
		return false, nil
	}
	return true, nil
}

func (p *RedisPoolClientV2) SPop(key string) (string, error) {
	cmd := p.redisRing.SPop(context.Background(), key)
	return cmd.Result()
}

func (p *RedisPoolClientV2) LIndex(key string, index int64) (string, error) {
	cmd := p.redisRing.LIndex(context.Background(), key, index)
	return cmd.Result()
}

func (p *RedisPoolClientV2) LPop(key string) (string, error) {
	cmd := p.redisRing.LPop(context.Background(), key)
	return cmd.Result()
}

func (p *RedisPoolClientV2) RPush(key, val string) (int64, error) {
	cmd := p.redisRing.RPush(context.Background(), key, val)
	return cmd.Result()
}

func (p *RedisPoolClientV2) LPush(key string, val string) (int64, error) {
	cmd := p.redisRing.LPush(context.Background(), key, val)
	return cmd.Result()
}

func (p *RedisPoolClientV2) HLen(key string) (int64, error) {
	cmd := p.redisRing.HLen(context.Background(), key)
	return cmd.Result()
}
