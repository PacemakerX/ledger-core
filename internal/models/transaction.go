package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)


type Transaction struct{
	ID  				uuid.UUID  			`db:"id"`
	IdempotencyKey		string 				`db:"idempotency_key"`
	Type 				string 				`db:"type"`
	Status 				string 				`db:"status"`
	InitiatedBy 		uuid.UUID 			`db:"initiated_by"`
	Metadata  			*json.RawMessage 	`db:"metadata"`
	CreatedAt 			time.Time  			`db:"created_at"`
	CompletedAt 		*time.Time 			`db:"completed_at"`
}