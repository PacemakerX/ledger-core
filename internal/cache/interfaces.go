package cache

import (
	"context"

	"github.com/PacemakerX/ledger-core/internal/models"
)

type IdempotencyCache interface {
	Get(ctx context.Context, key string) (*models.IdempotencyKey, error)
	Set(ctx context.Context, key *models.IdempotencyKey) error
}
