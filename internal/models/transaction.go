package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                    uuid.UUID        `db:"id"`
	FromAccountID         *uuid.UUID       `db:"from_account_id"`
	ToAccountID           *uuid.UUID       `db:"to_account_id"`
	IdempotencyKey        string           `db:"idempotency_key"`
	Amount                int64            `db:"amount"`
	CurrencyID            int              `db:"currency_id"`
	Type                  string           `db:"type"`
	Status                string           `db:"status"`
	InitiatedBy           uuid.UUID        `db:"initiated_by"`
	Metadata              *json.RawMessage `db:"metadata"`
	OriginalTransactionID *uuid.UUID       `db:"original_transaction_id"`
	CreatedAt             time.Time        `db:"created_at"`
	CompletedAt           *time.Time       `db:"completed_at"`
}
