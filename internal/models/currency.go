package models

import "time"


type Currency struct {
	ID				int			`db:"id"`
	Code    		string		`db:"code"`
	Name 			string		`db:"name"`
	Symbol  		string		`db:"symbol"`
	IsActive 		bool		`db:"is_active"`
	CreatedAt 		time.Time   `db:"created_at"`
	UpdatedAt 		time.Time   `db:"updated_at"`
}