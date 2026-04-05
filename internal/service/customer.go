package service

import (
	"context"
	"fmt"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
)

type CreateCustomerRequest struct {
	FirstName    string  `json:"first_name"`
	MiddleName   *string `json:"middle_name"`
	LastName     string  `json:"last_name"`
	AadharNumber *string `json:"aadhar_number"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	CountryCode  string  `json:"country_code"`
}

type CreateCustomerResponse struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Status     string    `json:"status"`
	KycStatus  string    `json:"kyc_status"`
}

type UpdateKYCRequest struct {
	Status string `json:"status"`
}

type customerService struct {
	customer repository.CustomerRepository
	country  repository.CountryRepository
}

func NewCustomerService(
	customer repository.CustomerRepository,
	country repository.CountryRepository,
) *customerService {
	return &customerService{
		customer: customer,
		country:  country,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*CreateCustomerResponse, error) {

	// Step 1 — look up country by code
	country, err := s.country.GetByCode(ctx, req.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("customerService.CreateCustomer: invalid country code: %w", err)
	}

	// Step 2 — build customer model
	customer := &models.Customer{
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		AadharNumber: req.AadharNumber,
		Email:        req.Email,
		CountryID:    country.ID,
		PhoneNumber:  req.PhoneNumber,
		KycStatus:    "unverified",
	}

	// Step 3 — create customer
	created, err := s.customer.Create(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("customerService.CreateCustomer: %w", err)
	}

	return &CreateCustomerResponse{
		CustomerID: created.ID,
		Status:     "created",
		KycStatus:  created.KycStatus,
	}, nil
}

func (s *customerService) UpdateKYC(ctx context.Context, id uuid.UUID, req UpdateKYCRequest) error {

	// Validate status
	validStatuses := map[string]bool{
		"pending":  true,
		"verified": true,
		"rejected": true,
	}

	if !validStatuses[req.Status] {
		return fmt.Errorf("customerService.UpdateKYC: invalid status %q — must be pending, verified or rejected", req.Status)
	}

	// Fetch customer — verify it exists
	_, err := s.customer.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customerService.UpdateKYC: %w", err)
	}

	// Update KYC status
	err = s.customer.UpdateKYC(ctx, id, req.Status)
	if err != nil {
		return fmt.Errorf("customerService.UpdateKYC: %w", err)
	}

	return nil
}
