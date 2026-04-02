package postgres

import (
	"context"
	"errors"
	"fmt"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type customerRepository struct{
	pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) repository.CustomerRepository{
	return &customerRepository{pool: pool}
}

func (r *customerRepository) GetByID(ctx context.Context,id uuid.UUID)(*models.Customer,error){

	query:=`SELECT id, first_name, middle_name, last_name, aadhar_number, country_id, phone_number, email, kyc_status, is_active,
	created_at, updated_at
	FROM customers
	WHERE id = $1
	`

	var Customer models.Customer

	err:= r.pool.QueryRow(ctx,query,id).Scan(
		&Customer.ID,
		&Customer.FirstName,
		&Customer.MiddleName,
		&Customer.LastName,
		&Customer.AadharNumber,
		&Customer.CountryID,
		&Customer.PhoneNumber,
		&Customer.Email,
		&Customer.KycStatus,
		&Customer.IsActive,
		&Customer.CreatedAt,
		&Customer.UpdatedAt,
	)	

	if(err!=nil){

		if errors.Is(err,pgx.ErrNoRows){
			return nil,domainerrors.ErrNotFound
		}
		return nil, fmt.Errorf("customerRepository.GetByID: %w", 
		domainerrors.ErrDatabase)
	}

	return &Customer,nil
}