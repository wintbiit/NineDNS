package cache

import (
	"context"
	"time"

	"github.com/wintbiit/ninedns/model"

	"github.com/redis/go-redis/v9"
	"github.com/wintbiit/ninedns/utils"
)

const redisKeyPrefix = "ninedns"

type RedisClient struct {
	*redis.Client
	API

	Domain string
	TTL    uint32
}

type API interface {
	FindRecords(name, qType, identify string) ([]model.Record, error)
	AddRecord(identify string, record *model.Record) error
	AddRuntimeCache(key string, value string, expire time.Duration) error
	GetRuntimeCache(key string) (string, error)
	Close() error
}

func NewClient(domain string, ttl uint32) (*RedisClient, error) {
	client := &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     utils.C.Redis.Addr,
			Password: utils.C.Redis.Password,
			DB:       utils.C.Redis.DB,
		}),
		Domain: domain,
		TTL:    ttl,
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	client.Set(context.Background(), redisKeyPrefix+":lastrun", time.Now().String(), 0)
	client.SAdd(context.Background(), redisKeyPrefix+":domains", domain)

	return client, nil
}
