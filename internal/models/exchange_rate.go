package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	
	ID 					int  				`db:"id"`
	FromCurrencyID 		int					`db:"from_currency"`
	ToCurrencyID 		int					`db:"to_currency"`	
	Rate				decimal.Decimal		`db:"rate"`
	Source 				string   			`db:"source"`
	EffectiveDate   	time.Time			`db:"effective_date"`
	CreatedAt 	 		time.Time			`db:"created_at"`	
}