package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)


type AuditLog struct{
	ID  			uuid.UUID 			`db:"id"`
	EntityType 		string				`db:"entity_type"`
	EntityID 		uuid.UUID			`db:"entity_id"`
	Action  		string 				`db:"action"`
	ActorID  		uuid.UUID 			`db:"actor_id"`
	ActorType  		string 				`db:"actor_type"`
	OldValue 		*json.RawMessage    `db:"old_value"`
	NewValue 		*json.RawMessage	`db:"new_value"`
	IPAddress 		*string				`db:"ip_address"`
	CreatedAt 		time.Time			`db:"created_at"`
}