package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/redis/go-redis/v9"
)

type idempotencyCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewIdempotencyCache(client *redis.Client) *idempotencyCache {
	return &idempotencyCache{
		client: client,
		ttl:    24 * time.Hour,
	}
}

func key(k string) string {
	return fmt.Sprintf("idempotency:%s", k)
}

func (c *idempotencyCache) Get(ctx context.Context, k string) (*models.IdempotencyKey, error) {
	val, err := c.client.Get(ctx, key(k)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("idempotencyCache.Get: %w", domainerrors.ErrDatabase)
	}

	var idempotencyKey models.IdempotencyKey
	if err := json.Unmarshal([]byte(val), &idempotencyKey); err != nil {
		return nil, fmt.Errorf("idempotencyCache.Get: unmarshal: %w", err)
	}

	return &idempotencyKey, nil
}

func (c *idempotencyCache) Set(ctx context.Context, k *models.IdempotencyKey) error {
	val, err := json.Marshal(k)
	if err != nil {
		return fmt.Errorf("idempotencyCache.Set: marshal: %w", err)
	}

	if err := c.client.Set(ctx, key(k.Key), val, c.ttl).Err(); err != nil {
		return fmt.Errorf("idempotencyCache.Set: %w", domainerrors.ErrDatabase)
	}

	return nil
}