package cache

import (
	"context"
	"fmt"

	"github.com/PacemakerX/ledger-core/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("cache.NewRedisClient: invalid redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache.NewRedisClient: ping failed: %w", err)
	}

	return client, nil
}
