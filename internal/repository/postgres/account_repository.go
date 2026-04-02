package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
)

type accountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) repository.AccountRepository {
    return &accountRepository{pool: pool}
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {

	query := `
		SELECT id, account_number, customer_id, currency_id, type_id,
		       country_id, is_active, daily_debit_limit, daily_credit_limit,
		       created_at, updated_at
		FROM accounts
		WHERE id = $1`

	var account models.Account
		err := r.pool.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.AccountNumber,
		&account.CustomerID,
		&account.CurrencyID,
		&account.TypeID,
		&account.CountryID,
		&account.IsActive,
		&account.DailyDebitLimit,
		&account.DailyCreditLimit,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("accountRepository.GetByID: %w", domainerrors.ErrDatabase)
	}

	return &account, nil
}

func (r *accountRepository) GetByIDForUpdate(ctx context.Context, tx repository.Tx, id uuid.UUID) (*models.Account, error) {

	query := `
		SELECT id, account_number, customer_id, currency_id, type_id, 
			country_id, is_active, daily_debit_limit, daily_credit_limit,
			created_at, updated_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE`

	var account models.Account
	// 	Type assertion
	// In Go, when you have an interface, you only see the methods that interface declares. The concrete type underneath is hidden.
	// Type assertion is you saying — "I know what's actually hiding under this interface, let me reach in and get it."

		pgxTx := tx.(pgx.Tx) // this is type assertions  
		err := pgxTx.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.AccountNumber,
		&account.CustomerID,
		&account.CurrencyID,
		&account.TypeID,
		&account.CountryID,
		&account.IsActive,
		&account.DailyDebitLimit,
		&account.DailyCreditLimit,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("accountRepository.GetByID: %w", domainerrors.ErrDatabase)
	}

	return &account, nil
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
    
	query:=`
		INSERT INTO accounts ( account_number, customer_id, currency_id, type_id,country_id, is_active, daily_debit_limit, daily_credit_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8 )
		RETURNING id, account_number, customer_id, currency_id, type_id,country_id, is_active, daily_debit_limit, daily_credit_limit, created_at, updated_at`

		err := r.pool.QueryRow(ctx, query,
			account.AccountNumber,
			account.CustomerID,
			account.CurrencyID,
			account.TypeID,
			account.CountryID,
			account.IsActive,
			account.DailyDebitLimit,
			account.DailyCreditLimit,
		).Scan(
			&account.ID,
			&account.AccountNumber,
			&account.CustomerID,
			&account.CurrencyID,
			&account.TypeID,
			&account.CountryID,
			&account.IsActive,
			&account.DailyDebitLimit,
			&account.DailyCreditLimit,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("accountRepository.Create: %w", domainerrors.ErrDatabase)
	}
    return account, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, tx repository.Tx, id uuid.UUID, newBalance int64) error {
 	// TODO: revisit when balance caching strategy is decided
    return nil
}

func (r *accountRepository) GetDailySpend(ctx context.Context, accountID uuid.UUID, date time.Time) (int64, error) {
    
	query:=`SELECT COALESCE(SUM(amount), 0)
			FROM journal_entries
			WHERE account_id = $1
			AND entry_type = 'DEBIT'
			AND created_at >= $2
			AND created_at < $3`

	var dailySpend int64

	err := r.pool.QueryRow(ctx, query,
		accountID,
		date,
		date.AddDate(0, 0, 1),
	).Scan(&dailySpend)

	if err != nil {
		return 0, fmt.Errorf("accountRepository.GetDailySpend: %w", domainerrors.ErrDatabase)
	}

	return dailySpend, nil
}