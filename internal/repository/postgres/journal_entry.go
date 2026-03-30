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

type journalEntryRepository struct{
	pool *pgxpool.Pool	
}

func NewJournalEntryRepository( pool *pgxpool.Pool) repository.JournalEntryRepository{
	return &journalEntryRepository{pool:pool}

}

func ( r *journalEntryRepository) CreateBatch(ctx context.Context, tx repository.Tx, entries []models.JournalEntry) error{


	batch := &pgx.Batch{}
	for _, entry := range entries {
		batch.Queue(`INSERT INTO journal_entries (transaction_id, account_id, entry_type, amount, currency_id, exchange_rate_id, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			entry.TransactionID,
			entry.AccountID,
			entry.EntryType,
			entry.Amount,
			entry.CurrencyID,
			entry.ExchangeRateID,
			entry.Description,
		)
	}

	pgxTx:=tx.(pgx.Tx)

	results:=pgxTx.SendBatch(ctx,batch)
	defer results.Close();

	for range entries {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("journalEntryRepository.CreateBatch: %w", domainerrors.ErrDatabase)
		}
	}

	return nil
}

func (r *journalEntryRepository) VerifyBalance(ctx context.Context, tx repository.Tx,transactionsID uuid.UUID) (int64,error){

	pgxTx:=tx.(pgx.Tx)
	query:=`SELECT COALESCE(
			SUM(CASE WHEN entry_type = 'DEBIT' THEN amount ELSE 0 END) -
			SUM(CASE WHEN entry_type = 'CREDIT' THEN amount ELSE 0 END)
		, 0)
		FROM journal_entries
		WHERE transaction_id = $1`

	var balance int64

	err := pgxTx.QueryRow(ctx, query,
		transactionsID,
	).Scan(&balance)

	if err != nil {
		return 0, fmt.Errorf("journalEntryRepository.VerifyBalance: %w", domainerrors.ErrDatabase)
	}

	return balance, nil
}