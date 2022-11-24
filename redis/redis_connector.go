package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConnector interface {
	GetClient() *redis.Client
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Incr(key string) *redis.IntCmd
	Decr(key string) *redis.IntCmd
	Del(key string) *redis.IntCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Publish(channel string, payload []byte) *redis.IntCmd
	Scan(pattern string) *redis.ScanCmd

	FlushAll() *redis.StatusCmd
}

type redisConnector struct {
	sync.Mutex
	Client *redis.Client
}

func NewRedisConnector(config *ConnectorConfig) *redisConnector {
	switch config.Mode {
	case Standalone:
		if config.RedisConfig != nil {
			rdb := redis.NewClient(&redis.Options{
				Addr:       config.RedisConfig.Addr,
				Username:   config.RedisConfig.Username,
				Password:   config.RedisConfig.Password,
				DB:         0,
				MaxRetries: 3,
			})
			return &redisConnector{
				Client: rdb,
			}
		}
	}
	return nil
}

func (r *redisConnector) GetClient() *redis.Client {
	return r.Client
}

func (r *redisConnector) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctx := context.Background()
	return r.Client.SetNX(ctx, key, value, expiration)
}

func (r *redisConnector) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration)
}

func (r *redisConnector) Get(key string) *redis.StringCmd {
	ctx := context.Background()
	return r.Client.Get(ctx, key)
}

func (r *redisConnector) Incr(key string) *redis.IntCmd {
	ctx := context.Background()
	return r.Client.Incr(ctx, key)
}

func (r *redisConnector) Decr(key string) *redis.IntCmd {
	ctx := context.Background()
	return r.Client.Decr(ctx, key)
}
func (r *redisConnector) Del(key string) *redis.IntCmd {
	ctx := context.Background()
	return r.Client.Del(ctx, key)
}

func (r *redisConnector) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	ctx := context.Background()
	return r.Client.Expire(ctx, key, expiration)
}

func (r *redisConnector) Publish(channel string, payload []byte) *redis.IntCmd {
	ctx := context.Background()
	return r.Client.Publish(ctx, channel, payload)
}

func (r *redisConnector) Scan(pattern string) *redis.ScanCmd {
	ctx := context.Background()
	return r.Client.Scan(ctx, 0, pattern, 0)
}

func (r *redisConnector) FlushAll() *redis.StatusCmd {
	ctx := context.Background()
	return r.Client.FlushAll(ctx)
}
