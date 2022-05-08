package db

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var MyRedis myRedis = myRedis{}

type myRedis struct {
	Client *redis.Client
	Cache  *cache.Cache
}

func (r *myRedis) InitRedis() {

	r.Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISADDRESS"),
		Password: "",
		DB:       0,
	})
	r.Cache = cache.New(&cache.Options{
		Redis: r.Client,
	})
}

func (r *myRedis) SetCache(key string, ctx context.Context, value interface{}, ttl time.Duration) error {
	if err := r.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}

	return nil
}

func (r *myRedis) DeleteCache(key string, ctx context.Context) error {
	if err := r.Cache.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}

func (r *myRedis) GetCache(key string, ctx context.Context, value interface{}) error {
	err := r.Cache.Get(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}
