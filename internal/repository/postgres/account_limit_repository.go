package postgres

import (
	"context"
	"fmt"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type accountLimitRepository struct {
	pool *pgxpool.Pool
}

func NewAccountLimitRepository(pool *pgxpool.Pool) repository.AccountLimitRepository {
	return &accountLimitRepository{pool: pool}
}


func (r *accountLimitRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.AccountLimit, error) {
		
	
	query:=`SELECT id, account_id, limit_type, max_amount, current_usage, reset_at, created_at, updated_at
			FROM account_limits
			WHERE account_id = $1`
	rows, err := r.pool.Query(ctx, query, accountID)
	
	if err != nil {
		return nil, fmt.Errorf("accountLimitRepository.GetByAccountID: %w", domainerrors.ErrDatabase)
	}
	defer rows.Close()

	var limits []models.AccountLimit
	for rows.Next() {
		var limit models.AccountLimit
		err := rows.Scan(
			&limit.ID,
			&limit.AccountID,
			&limit.LimitType,
			&limit.MaxAmount,
			&limit.CurrentUsage,
			&limit.ResetAt,
			&limit.CreatedAt,
			&limit.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("accountLimitRepository.GetByAccountID: %w", domainerrors.ErrDatabase)
		}
		limits = append(limits, limit)

	}
	if err := rows.Err(); err != nil {
   		return nil, fmt.Errorf("accountLimitRepository.GetByAccountID: %w", domainerrors.ErrDatabase)
	}
	return limits, nil
}

func (r *accountLimitRepository) UpdateUsage(ctx context.Context, tx repository.Tx, limitID uuid.UUID, amount int64) error {
    
	pgxTx:=tx.(pgx.Tx)

	query:=`UPDATE account_limits
			SET current_usage= current_usage + $2
			WHERE id = $1`
	
	_,err:= pgxTx.Exec(ctx,query,limitID,amount)
	
	if(err!=nil){

	return fmt.Errorf("accountLimitRepository.UpdateUsage: %w", domainerrors.ErrDatabase) 

	}
	
	return nil
}