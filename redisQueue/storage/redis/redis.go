package storage

import (
	"context"
	"encoding/json"
	"log/slog"
	"redisQueue/domains/models"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	client  *redis.Client
	maxLen  int
	address string
	key     string
}

func New(address, key string, maxLen int) *RedisProvider {
	slog.Info(address, key, maxLen)
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
	slog.Info(m)

	script := redis.NewScript(`
		local k = KEYS[1]
		local payload = ARGV[1]
		local maxlen = tonumber(ARGV[2])
		local repeatNum = tonumber(ARGV[3])
		local args = {k}
		for i = 1, repeatNum do
		  args[#args+1] = payload
		end
		redis.call('LPUSH', unpack(args))
		redis.call('LTRIM', k, 0, maxlen - 1)
		return redis.call('LLEN', k)
`)
	script.Run(ctx, provider.client, []string{provider.key}, m, provider.maxLen, repeatNum)
	return nil
}
func (provider *RedisProvider) Get(ctx context.Context) ([]models.Msg, error) {
	raw, err := provider.client.LRange(ctx, provider.key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	out := make([]models.Msg, 0, len(raw))
	for _, v := range raw {
		var msg models.Msg
		if err := json.Unmarshal([]byte(v), &msg); err == nil {
			out = append(out, msg)
		}
	}

	return out, nil
}
