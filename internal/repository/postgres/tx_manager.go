package postgres

import (
	"context"

	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManagerRepository struct{
	pool *pgxpool.Pool;
}

func NewTxManager(pool *pgxpool.Pool) repository.TxManager{
	return &txManagerRepository{pool: pool}
}

func (r *txManagerRepository)Begin(ctx context.Context) (repository.Tx, error){

	return r.pool.Begin(ctx)
}