package models

import (
	"time"

	"github.com/google/uuid"
)


type AccountLimit struct { 
	ID 				uuid.UUID		`db:"id"`
	AccountID 		uuid.UUID		`db:"account_id"`
	LimitType 		string			`db:"limit_type"`
	MaxAmount 		int64			`db:"max_amount"`
	CurrentUsage 	int64			`db:"current_usage"`
	ResetAt 		*time.Time      `db:"reset_at"`
	CreatedAt 		time.Time		`db:"created_at"`
	UpdatedAt 		time.Time		`db:"updated_at"`
}