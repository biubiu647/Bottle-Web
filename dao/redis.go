package dao

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	FOREVER = time.Duration(-1)
)

type Cache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var Field = []string{"msg", "img_path", "bottle_owner"}

func NewRedisCache(db int, host string, defaultExpiration time.Duration) *Cache {
	pool := &redis.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: time.Duration(100) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			//if _, err = conn.Do("AUTH", "XXXXXX"); err != nil {
			//	conn.Close()
			//	return nil, err
			//}
			return conn, nil
		},
	}
	return &Cache{pool: pool, defaultExpiration: defaultExpiration}
}

func (c Cache) Expire(name string, newSecondLiveTime int64) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", name, newSecondLiveTime)
	return err
}

func (c Cache) Delete(key ...interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("DEL", key...))
	return ok, err
}

func (c Cache) HDel(name, key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("HDEL", name, key))
	return ok, err
}

func (c Cache) HMGet(name string, fields ...string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	args := []interface{}{name}
	for _, field := range fields {
		args = append(args, field)
	}
	value, err := redis.Values(conn.Do("HMGET", args...))
	return value, err
}

func (c Cache) HGet(name, field string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bytes(conn.Do("GET", name, field))
	return v, err
}
func (c Cache) HSet(key, field string, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	return err
}

func (c Cache) Smembers(name string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bytes(conn.Do("smembers", name))
	return v, err
}

// 获取集合中元素的个数
func (c Cache) ScardInt64s(name string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int64(conn.Do("SCARD", name))
	return v, err
}
func (c Cache) Setmembers(name string, v ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	args := []interface{}{name}
	for _, vv := range v {
		args = append(args, vv)
	}
	_, err := conn.Do("SADD", args)
	return err
}

func (c Cache) Conn() redis.Conn {
	conn := c.pool.Get()
	return conn
}
