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

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) repository.TransactionRepository {
	return &transactionRepository{pool: pool}
}

func (r *transactionRepository) Create(ctx context.Context, tx repository.Tx, transaction *models.Transaction) (*models.Transaction, error) {

	pgxTx := tx.(pgx.Tx)

	query := `INSERT INTO transactions(idempotency_key, type, status, initiated_by, metadata, from_account_id, to_account_id, amount, currency_id, original_transaction_id)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, idempotency_key, type, status, initiated_by, metadata, from_account_id, to_account_id, amount, currency_id, original_transaction_id, created_at, completed_at`

	err := pgxTx.QueryRow(ctx, query,
		transaction.IdempotencyKey,
		transaction.Type,
		transaction.Status,
		transaction.InitiatedBy,
		transaction.Metadata,
		transaction.FromAccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.CurrencyID,
		transaction.OriginalTransactionID,
	).Scan(
		&transaction.ID,
		&transaction.IdempotencyKey,
		&transaction.Type,
		&transaction.Status,
		&transaction.InitiatedBy,
		&transaction.Metadata,
		&transaction.FromAccountID,
		&transaction.ToAccountID,
		&transaction.Amount,
		&transaction.CurrencyID,
		&transaction.OriginalTransactionID,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("transactionRepository.Create: %w", domainerrors.ErrDatabase)
	}
	return transaction, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, tx repository.Tx, id uuid.UUID, status string) error {

	pgxTx := tx.(pgx.Tx)

	query := `UPDATE transactions
			SET status= $3,
			completed_at= $2
			WHERE id = $1`

	_, err := pgxTx.Exec(ctx, query, id, time.Now(), status)

	if err != nil {
		return fmt.Errorf("transactionRepository.UpdateStatus: %w", domainerrors.ErrDatabase)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	query := `SELECT id, idempotency_key, type, status, initiated_by, metadata, 
              from_account_id, to_account_id, amount, currency_id, original_transaction_id, created_at, completed_at 
              FROM transactions
              WHERE id = $1`

	var transaction models.Transaction
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.IdempotencyKey,
		&transaction.Type,
		&transaction.Status,
		&transaction.InitiatedBy,
		&transaction.Metadata,
		&transaction.FromAccountID,
		&transaction.ToAccountID,
		&transaction.Amount,
		&transaction.CurrencyID,
		&transaction.OriginalTransactionID,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		fmt.Printf("GetByID scan error: %v\n", err)  // add this
		return nil, fmt.Errorf("transactionRepository.GetByID: %w", domainerrors.ErrDatabase)
	}

	return &transaction, nil
}

func (r *transactionRepository) GetTotalRefunded(ctx context.Context, originalTransactionID uuid.UUID) (int64, error) {
	query := `SELECT COALESCE(SUM(amount), 0)
              FROM transactions
              WHERE original_transaction_id = $1
              AND type = 'REFUND'
              AND status = 'COMPLETED'`

	var total int64
	err := r.pool.QueryRow(ctx, query, originalTransactionID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("transactionRepository.GetTotalRefunded: %w", domainerrors.ErrDatabase)
	}
	return total, nil
}
