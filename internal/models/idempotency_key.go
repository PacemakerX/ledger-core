package models

import "time"


type IdempotencyKey struct{
	Key  			string 			`db:"id"`
	RequestHash 	string			`db:"request_hash"`	
	ResponseStatus  string			`db:"response_status"`
	ResponseBody    string			`db:"response_body"`
	ExpiresAt       time.Time		`db:"expires_at"`
	CreatedAt 		time.Time		`db:"created_at"`
}