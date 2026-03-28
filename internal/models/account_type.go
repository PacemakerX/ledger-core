package models

import "time"

type AccountType struct {
	ID 					int  		`db:"id"`
	Name 				string  	`db:"name"`
	NormalBalance 		string 		`db:"normal_balance"`
	description    		string 		`db:"description"`
	CreatedAt   		time.Time  	`db:"created_at"`
}