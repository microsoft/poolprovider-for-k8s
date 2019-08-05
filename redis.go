package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Storage interface {
	Ping() (string, error)
	Get(string) (string, error)
	Set(string, string) error
}

type RedisStorage struct {
	connectionPool *redis.Pool
}

func NewRedisStorage(server string) Storage {
	return &RedisStorage{
		connectionPool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

// Ping will "ping" the storage backend.
// The function can use to check the connection from
// the app to the storage backend.
func (r *RedisStorage) Ping() (string, error) {
	conn := r.connectionPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("PING"))
	return res, err
}

func (r *RedisStorage) Get(key string) (string, error) {
	conn := r.connectionPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("GET", key))
	return res, err
}

func (r *RedisStorage) Set(key string, value string) error {
	conn := r.connectionPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", key, value))
	return err
}