package models

import (
	"time"

	"github.com/google/uuid"
)


type Account struct{
	ID  				uuid.UUID 		`db:"id"`
	AccountNumber 		string 			`db:"account_number"`
	CustomerID  		uuid.UUID 		`db:"customer_id"`
	CurrencyID 			int 		    `db:"currency_id"`
	TypeID 				int   			`db:"type_id"`
	CountryID 			int   			`db:"country_id"`
	IsActive 			bool  			`db:"is_active"`
	DailyDebitLimit 	*int64 			`db:"daily_debit_limit"`
	DailyCreditLimit 	*int64 			`db:"daily_credit_limit"`
	CreatedAt 			time.Time  		`db:"created_at"`
	UpdatedAt 			time.Time 		`db:"updated_at"`
}