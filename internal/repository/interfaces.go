package repository

import (
	"context"
	"time"

	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/google/uuid"
)

// Tx abstracts a database transaction.
// pgx.Tx satisfies this interface automatically via structural typing —
// no wrapper needed. Service layer stays free of pgx imports.
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// TxManager abstracts BEGIN. The postgres adapter implements this.
// Service layer calls Begin, passes the Tx into repository methods,
// then calls Commit or Rollback — never touches pgxpool directly.
type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}

// AccountRepository defines read and write operations on accounts.
// Methods that mutate state inside a transaction accept a Tx parameter.
// Read-only methods that run outside a transaction do not.
type AccountRepository interface {
	// GetByID fetches an account by primary key. Used for validation
	// and reads that do not require locking.
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)

	// GetByIDForUpdate fetches an account inside a transaction with
	// SELECT FOR UPDATE. Prevents lost updates during concurrent transfers.
	// Always lock accounts in ascending UUID order to prevent deadlocks.
	GetByIDForUpdate(ctx context.Context, tx Tx, id uuid.UUID) (*models.Account, error)

	// Create inserts a new account record.
	Create(ctx context.Context, account *models.Account) (*models.Account, error)

	// UpdateBalance sets the balance on an account inside a transaction.
	// Balance is not stored as a column — it is derived from journal entries.
	// This method updates daily_debit_limit / daily_credit_limit usage tracking
	// or any denormalised balance cache if one is added later.
	// NOTE: actual balance is always computed via SUM of journal entries.
	UpdateBalance(ctx context.Context, tx Tx, id uuid.UUID, newBalance int64) error

	// GetDailySpend returns total amount debited from an account on a given date.
	// Used to enforce daily debit limits before committing a transfer.
	GetDailySpend(ctx context.Context, accountID uuid.UUID, date time.Time) (int64, error)
}

// TransactionRepository defines operations on the transactions table.
// Transactions are created inside a db transaction and updated on completion.
type TransactionRepository interface {
	// Create inserts a new transaction record with PENDING status.
	// Must be called inside a db transaction (tx).
	Create(ctx context.Context, tx Tx, transaction *models.Transaction) (*models.Transaction, error)

	// UpdateStatus transitions a transaction to COMPLETED or FAILED.
	// Must be called inside the same db transaction.
	UpdateStatus(ctx context.Context, tx Tx, id uuid.UUID, status string) error

	// GetByID fetches a transaction by primary key. Used for idempotency
	// lookups and status checks. Runs outside a transaction.
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)

	GetTotalRefunded(ctx context.Context, originalTransactionID uuid.UUID) (int64, error)
}

// JournalEntryRepository defines operations on the journal_entries table.
// Journal entries are immutable — no Update or Delete methods exist by design.
type JournalEntryRepository interface {
	// CreateBatch inserts multiple journal entries atomically inside a transaction.
	// A transfer always creates a minimum of 4 entries (2 per account touched).
	// All entries must be inserted together — partial inserts are never valid.
	CreateBatch(ctx context.Context, tx Tx, entries []models.JournalEntry) error

	// VerifyBalance checks that SUM(debits) - SUM(credits) = 0 for a given
	// transaction. Called after CreateBatch and before Commit as a safeguard.
	// Returns the net sum — caller must assert it equals zero.
	VerifyBalance(ctx context.Context, tx Tx, transactionID uuid.UUID) (int64, error)

	GetBalance(ctx context.Context, accountID uuid.UUID) (int64, error)
}

// IdempotencyRepository defines operations on the idempotency_keys table.
// Ensures exactly-once semantics for transfer requests.
type IdempotencyRepository interface {
	// Get fetches an existing idempotency record by key string.
	// Returns nil, nil if the key does not exist (not an error).
	Get(ctx context.Context, key string) (*models.IdempotencyKey, error)

	// Create inserts a new idempotency key at the start of a transfer,
	// before the db transaction begins. Acts as a reservation.
	Create(ctx context.Context, tx Tx, idempotencyKey *models.IdempotencyKey) error

	// SetResponse stores the final response body and status on the key
	// after the transfer completes. Must be called inside the db transaction
	// so the response is only persisted if the transfer committed.
	SetResponse(ctx context.Context, tx Tx, key string, responseStatus string, responseBody string) error
}

// CustomerRepository defines read operations on the customers table.
// Customers are never mutated by the transfer flow — only read for KYC checks.
type CustomerRepository interface {
	// GetByID fetches a customer by primary key.
	// Used to validate KYC status before allowing a transfer.
	GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error)

	Create(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	UpdateKYC(ctx context.Context, id uuid.UUID, status string) error
}

// AccountLimitRepository defines operations on the account_limits table.
// Limits are checked before a transfer and updated after it commits.
type AccountLimitRepository interface {
	// GetByAccountID returns all limit records for a given account.
	// The service layer filters by LimitType (DAILY, MONTHLY, TRANSACTION).
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.AccountLimit, error)

	// UpdateUsage increments current_usage for a limit record inside a transaction.
	// Called after journal entries are written, before Commit.
	UpdateUsage(ctx context.Context, tx Tx, limitID uuid.UUID, amount int64) error
}

type CountryRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Country, error)
}

type CurrencyRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Currency, error)
}

type AccountTypeRepository interface {
	GetByName(ctx context.Context, name string) (*models.AccountType, error)
}
