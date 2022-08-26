package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,  // default DB
		PoolSize: 10, // 连接池
	})
	result := rdb.Ping(context.Background())
	fmt.Println("redis ping: ", result.Val())
	if result.Val() != "PONG" {
		return nil
	}
	return rdb
}

func RedisClose() {
	if err := Redis.Close(); err != nil {
		return
	}
}
