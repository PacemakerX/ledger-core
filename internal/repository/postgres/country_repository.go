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

type countryRepository struct {
	pool *pgxpool.Pool
}

func NewCountryRepository(pool *pgxpool.Pool) *countryRepository {
	return &countryRepository{pool: pool}
}

func (r *countryRepository) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	query := `SELECT id, name, iso_code, dial_code, currency_id, created_at
			  FROM countries
			  WHERE iso_code = $1`

	var country models.Country
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&country.ID,
		&country.Name,
		&country.ISOCode,
		&country.DialCode,
		&country.CurrencyID,
		&country.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("countryRepository.GetByCode: %w", domainerrors.ErrDatabase)
	}
	return &country, nil
}
