package postgres

import (
	"context"
	"errors"
	"fmt"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type idempotencyRepository struct{
	pool *pgxpool.Pool
}

func NewIdempotencyRepository(pool *pgxpool.Pool) repository.IdempotencyRepository{
	return &idempotencyRepository{ pool: pool}
}


func (r *idempotencyRepository)  Get(ctx context.Context, key string) (*models.IdempotencyKey, error){

	query:=`SELECT key, request_hash, response_status,
	response_body, expires_at, created_at
	FROM idempotency_keys
	WHERE key = $1`
	
	var idempotencyKey models.IdempotencyKey

	err:=r.pool.QueryRow(ctx,query,key).Scan(
		&idempotencyKey.Key,
		&idempotencyKey.RequestHash,
		&idempotencyKey.ResponseStatus,
		&idempotencyKey.ResponseBody,
		&idempotencyKey.ExpiresAt,
		&idempotencyKey.CreatedAt)

	if(err!=nil){

		if errors.Is(err,pgx.ErrNoRows){
			return nil,domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("IdempotencyRepository.GetByKey: %w", 
		domainerrors.ErrDatabase)
	}
	return &idempotencyKey,nil
}

func (r *idempotencyRepository)  Create(ctx context.Context, idempotencyKey *models.IdempotencyKey) error{

	query:=`INSERT INTO idempotency_keys (key, request_hash, response_status, response_body,expires_at)
	VALUES($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query,
		idempotencyKey.Key,
		idempotencyKey.RequestHash,
		idempotencyKey.ResponseStatus,
		idempotencyKey.ResponseBody,
		idempotencyKey.ExpiresAt,
	)

	if(err!=nil){
		return fmt.Errorf("idempotencyRepository.Create %w",domainerrors.ErrDatabase) 
	}
	return nil
}
func (r *idempotencyRepository)  SetResponse(ctx context.Context, tx repository.Tx, key string, responseStatus string, responseBody string) error{
	
	query:=`UPDATE idempotency_keys
			SET response_status = $2,
			response_body = $3
			WHERE key = $1`


	pgxTx := tx.(pgx.Tx)

	_,err:=pgxTx.Exec(ctx,query,key,responseStatus,responseBody)

	if(err!=nil){

		return fmt.Errorf("idempotencyRepository.SetResponse: %w",domainerrors.ErrDatabase)
	}
	
	return nil
}