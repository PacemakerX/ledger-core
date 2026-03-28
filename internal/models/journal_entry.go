package models

import (
	"time"

	"github.com/google/uuid"
)


type JournalEntry struct {
	ID 					uuid.UUID 		`db:"id"`
	TransactionID 		uuid.UUID 		`db:"transaction_id"`
	AccountID  			uuid.UUID 		`db:"account_id"`
	EntryType 			string 			`db:"entry_type"`
	Amount 				int64 			`db:"amount"`
	CurrencyID 			int 			`db:"currency_id"`
	ExchangeRateID 		*int 			`db:"exchange_rate_id"`
	Description 		*string 		`db:"description"`
	CreatedAt 			time.Time  		`db:"created_at"`
}