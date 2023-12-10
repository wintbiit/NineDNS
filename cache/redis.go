package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wintbiit/ninedns/utils"
)

const redisKeyPrefix = "ninedns"

type RedisClient struct {
	*redis.Client
	API

	Domain string
	TTL    uint16
}

func NewClient(domain string) (*RedisClient, error) {
	client := &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     utils.C.Redis.Addr,
			Password: utils.C.Redis.Password,
			DB:       utils.C.Redis.DB,
		}),
		Domain: domain,
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
