package models

import (
	"time"

	"github.com/google/uuid"
)


type Customer struct {
	ID 					uuid.UUID 		`db:"id"`
	FirstName 			string 			`db:"first_name"`
	MiddleName 			*string 		`db:"middle_name"`
	LastName 			string 			`db:"last_name"`
	AadharNumber 		*string 			`db:"aadhar_number"`
	CountryID     		int 			`db:"country_id"`
	PhoneNumber 		string 			`db:"phone_number"`
	Email 				string 			`db:"email"`
	KycStatus 			string 			`db:"kyc_status"`
	IsActive 			bool 			`db:"is_active"`
	CreatedAt 			time.Time  		`db:"created_at"`
	UpdatedAt 			time.Time 		`db:"updated_at"`
}