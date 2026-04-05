package postgres

import (
	"context"
	"errors"
	"fmt"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountTypeRepository struct {
	pool *pgxpool.Pool
}

func NewAccountTypeRepository(pool *pgxpool.Pool) *accountTypeRepository {
	return &accountTypeRepository{pool: pool}
}

func (r *accountTypeRepository) GetByName(ctx context.Context, name string) (*models.AccountType, error) {
	query := `SELECT id, name, normal_balance, description, created_at
			  FROM account_types
			  WHERE name = $1`

	var accountType models.AccountType
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&accountType.ID,
		&accountType.Name,
		&accountType.NormalBalance,
		&accountType.Description,
		&accountType.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("accountTypeRepository.GetByName: %w", domainerrors.ErrDatabase)
	}
	return &accountType, nil
}
