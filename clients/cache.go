package clients

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type CacheService interface {
	AddToList(key string, values ...interface{}) error
	GetFromList(key string, offset int, limit int) ([]string, error)
	AddKey(key string, values any) error
	AddKeyTtl(key string, values any, ttl time.Duration) error
	GetKey(key string) (string, error)
	Rename(key string, newKey string) error
	KeyExists(key string) bool
	Count(key string) (cnt int64, err error)
	Del(key string) (bool, error)
}

func GetCacheService(host string, pwd string, slaveHost string, appKey string) CacheService {
	c := CacheServiceImpl{appKey: appKey}
	c.init(host, pwd, slaveHost)
	return &c
}

func (c *CacheServiceImpl) init(host string, pwd string, slaveHost string) CacheService {
	c.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	c.readOnlyClient = redis.NewClient(&redis.Options{
		Addr:     slaveHost,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	c.ttlNewVideo = time.Hour * 2160
	return c
}

type CacheServiceImpl struct {
	client         *redis.Client
	readOnlyClient *redis.Client
	ttlNewVideo    time.Duration
	appKey         string
}

func (c *CacheServiceImpl) Del(key string) (bool, error) {
	result, err := c.client.Del(context.Background(), key).Result()
	if err != nil || result == 0 {
		return false, err
	}
	return true, nil
}

func (c *CacheServiceImpl) Count(key string) (cnt int64, err error) {
	lLen := c.client.LLen(context.Background(), key)
	return lLen.Result()
}

func (c *CacheServiceImpl) KeyExists(key string) bool {
	val := c.client.Exists(context.Background(), key).Val()
	return val > 0
}

func (c *CacheServiceImpl) Rename(key string, newKey string) error {
	return c.client.Rename(context.Background(), key, newKey).Err()
}

func (c *CacheServiceImpl) AddToList(key string, values ...interface{}) error {
	err := c.client.RPush(context.Background(), key, values).Err()
	if err != nil {
		logrus.Warn(err)
		return err
	}
	return nil
}

func (c *CacheServiceImpl) GetFromList(key string, offset int, limit int) ([]string, error) {
	if !c.KeyExists(key) {
		return nil, fmt.Errorf("key not exists")
	}

	data, err := c.readOnlyClient.LRange(context.Background(), key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}
	return data, nil
}

func (c *CacheServiceImpl) AddKey(key string, values any) error {
	return c.AddKeyTtl(key, values, -1)
}

func (c *CacheServiceImpl) AddKeyTtl(key string, values any, ttl time.Duration) error {
	err := c.client.Set(context.Background(), key, values, ttl).Err()
	if err != nil {
		logrus.Warn(err)
	}
	return err
}

func (c *CacheServiceImpl) GetKey(key string) (string, error) {
	data, err := c.readOnlyClient.Get(context.Background(), key).Result()
	if err != nil {
		logrus.Debug(err)
		return "", err
	}
	return data, nil
}
