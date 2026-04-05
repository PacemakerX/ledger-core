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

type currencyRepository struct {
	pool *pgxpool.Pool
}

func NewCurrencyRepository(pool *pgxpool.Pool) *currencyRepository {
	return &currencyRepository{pool: pool}
}

func (r *currencyRepository) GetByCode(ctx context.Context, code string) (*models.Currency, error) {
	query := `SELECT id, code, name, symbol, is_active, created_at, updated_at
			  FROM currencies
			  WHERE code = $1`

	var currency models.Currency
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.IsActive,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("currencyRepository.GetByCode: %w", domainerrors.ErrDatabase)
	}
	return &currency, nil
}
