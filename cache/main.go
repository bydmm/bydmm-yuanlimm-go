package cache

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PW"),
		DB:       int(db),
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	RedisClient = client
}

// GetInt64 按照int64取缓存
func GetInt64(key string) (int64, error) {
	cache, err := RedisClient.Get(key).Result()
	if err == nil {
		intCache, err := strconv.ParseInt(cache, 10, 64)
		return intCache, err
	}
	return 0, err
}

// Fetch 有值返回，没值设置
func Fetch(key string, expiration time.Duration, callback func() string) (string, error) {
	cache, err := RedisClient.Get(key).Result()
	if err == nil {
		return cache, nil
	}
	cache = callback()
	err = RedisClient.Set(key, cache, expiration).Err()
	return cache, err
}
