package errors

import "errors"

// Sentinel errors for the ledger-core domain.
// Repository adapters wrap raw pgx errors into these.
// Service layer branches on these to determine HTTP response codes.
// Handler layer maps these to HTTP status codes via errors.Is().

var (
	// ErrNotFound is returned when a requested entity does not exist.
	// Maps to HTTP 404.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is returned when a unique constraint is violated.
	// Primary use: idempotency key already exists (duplicate request).
	// Maps to HTTP 409.
	ErrAlreadyExists = errors.New("already exists")

	// ErrInsufficientBalance is returned when a debit would make
	// an account balance go negative.
	// Maps to HTTP 422.
	ErrInsufficientBalance = errors.New("insufficient balance")

	// ErrDailyLimitExceeded is returned when a transfer would breach
	// the account's daily debit limit.
	// Maps to HTTP 422.
	ErrDailyLimitExceeded = errors.New("daily limit exceeded")

	// ErrMonthlyLimitExceeded is returned when a transfer would breach
	// the account's monthly debit limit.
	// Maps to HTTP 422.
	ErrMonthlyLimitExceeded = errors.New("monthly limit exceeded")

		// ErrMonthlyLimitExceeded is returned when a transfer would breach
	// the account's monthly debit limit.
	// Maps to HTTP 422.
	ErrYearlyLimitExceeded = errors.New("yearly limit exceeded")
	// ErrTransactionLimitExceeded is returned when a single transfer
	// amount exceeds the per-transaction limit on the account.
	// Maps to HTTP 422.
	ErrTransactionLimitExceeded = errors.New("transaction limit exceeded")

	// ErrKYCNotVerified is returned when a customer's KYC status
	// is not VERIFIED. Transfers are blocked until KYC clears.
	// Maps to HTTP 403.
	ErrKYCNotVerified = errors.New("KYC not verified")

	// ErrAccountInactive is returned when either the source or
	// destination account is not active.
	// Maps to HTTP 422.
	ErrAccountInactive = errors.New("account is inactive")

	// ErrIdempotencyConflict is returned when the same idempotency key
	// is reused with a different request hash — same key, different payload.
	// This is not a replay (which is fine) — it is a misuse.
	// Maps to HTTP 409.
	ErrIdempotencyConflict = errors.New("idempotency key conflict")

	// ErrBalanceVerificationFailed is returned when the double-entry
	// verification check fails — SUM(debits) - SUM(credits) != 0.
	// This should never happen in production. If it does, it signals
	// a bug in journal entry creation. Transaction is rolled back.
	// Maps to HTTP 500.
	ErrBalanceVerificationFailed = errors.New("balance verification failed: debits do not equal credits")

	// ErrDatabase is returned for unexpected postgresql errors that do
	// not map to a specific domain error. Wraps the underlying pgx error.
	// Maps to HTTP 500.
	ErrDatabase = errors.New("database error")
)