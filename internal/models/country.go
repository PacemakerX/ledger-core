package models

import "time"

type Country struct {
	ID 					int				`db:"id"`
	Name    			string			`db:"name"`
	ISOCode  			string			`db:"iso_code"`
	DialCode			string 			`db:"dial_code"`
	CurrencyID  		int				`db:"currency_id"`
	IsActive 			bool			`db:"is_active"`
	CreatedAt 			time.Time		`db:"created_at"`
}