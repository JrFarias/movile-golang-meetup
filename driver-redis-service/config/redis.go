package config

import "github.com/go-redis/redis"

// RedisConnect used to connect with redis
func RedisConnect(addr string, pass string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
	return client
}
