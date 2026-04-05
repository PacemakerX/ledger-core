package models

import "time"

type AccountType struct {
	ID 					int  		`db:"id"`
	Name 				string  	`db:"name"`
	NormalBalance 		string 		`db:"normal_balance"`
	Description    		string 		`db:"description"`
	CreatedAt   		time.Time  	`db:"created_at"`
}