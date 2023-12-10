package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/wintbiit/ninedns/model"
	"github.com/wintbiit/ninedns/utils"
)

type API interface {
	FindRecords(name, qType, identify string) ([]model.Record, error)
	AddRecord(identify string, record *model.Record) error
}

func (c *RedisClient) FindRecords(name, qType, identify string) ([]model.Record, error) {
	key := fmt.Sprintf("%s:%s:%s:%s:%s", redisKeyPrefix, c.Domain, name, qType, identify)
	values, err := c.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	records := make([]model.Record, len(values))
	for i, v := range values {
		if err := utils.UnmarshalFromString(v, &records[i]); err != nil {
			return nil, err
		}
	}

	return records, nil
}

func (c *RedisClient) AddRecord(identify string, record *model.Record) error {
	key := fmt.Sprintf("%s:%s:%s:%s:%s", redisKeyPrefix, c.Domain, record.Host, record.Type, identify)
	defer c.Expire(context.Background(), key, time.Duration(c.TTL)*time.Second)

	value, err := utils.MarshalToString(record)
	if err != nil {
		return err
	}
	return c.SAdd(context.Background(), key, value).Err()
}
