package cache

import (
	"context"
	"time"
)

const redisRuntimeKeyPrefix = redisKeyPrefix + ":runtime"

func (c *RedisClient) AddRuntimeCache(key string, value string, expire time.Duration) error {
	defer c.Expire(context.Background(), redisRuntimeKeyPrefix, expire)
	return c.HSet(context.Background(), redisRuntimeKeyPrefix, key, value).Err()
}

func (c *RedisClient) GetRuntimeCache(key string) (string, error) {
	return c.HGet(context.Background(), redisRuntimeKeyPrefix, key).Result()
}

func (c *RedisClient) ClearRuntimeCache() error {
	return c.Del(context.Background(), redisRuntimeKeyPrefix).Err()
}
