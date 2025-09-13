package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	client  *redis.Client
	maxLen  int
	address string
	key     string
}

func New(address, key string, maxLen int) *RedisProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})
	return &RedisProvider{
		client:  rdb,
		address: address,
		key:     key,
		maxLen:  maxLen,
	}
}

func (provider *RedisProvider) Add(ctx context.Context, repeatNum int, m string) error {
	script := redis.NewScript(`
		local k = KEYS[1]
		local payload = ARGV[1]
		local maxlen = tonumber(ARGV[2])
		local repeatNum = tonumber(ARGV[3])
		redis.call('LPUSH', k, payload)
		redis.call('LTRIM', k, 0, repeatNum - 1)
		return redis.call('LLEN', k)
`)
	script.Run(ctx, provider.client, []string{provider.key}, m, provider.maxLen, repeatNum)
	return nil
}
