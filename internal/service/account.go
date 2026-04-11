package service

import (
	"context"
	"fmt"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	CustomerID   uuid.UUID `json:"customer_id"`
	CurrencyCode string    `json:"currency_code"`
	AccountType  string    `json:"account_type"`
}

type CreateAccountResponse struct {
	AccountID     uuid.UUID `json:"account_id"`
	AccountNumber string    `json:"account_number"`
	CustomerID    uuid.UUID `json:"customer_id"`
	CurrencyCode  string    `json:"currency_code"`
	AccountType   string    `json:"account_type"`
	Status        string    `json:"status"`
}

type accountService struct {
	customer    repository.CustomerRepository
	account     repository.AccountRepository
	currency    repository.CurrencyRepository
	accountType repository.AccountTypeRepository
	auditLog    repository.AuditLogRepository
}

func NewAccountService(
	customer repository.CustomerRepository,
	account repository.AccountRepository,
	currency repository.CurrencyRepository,
	accountType repository.AccountTypeRepository,
	auditLog repository.AuditLogRepository,
) *accountService {
	return &accountService{
		customer:    customer,
		account:     account,
		currency:    currency,
		accountType: accountType,
		auditLog:    auditLog,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req CreateAccountRequest) (*CreateAccountResponse, error) {

	// Step 1 — customer must exist
	customer, err := s.customer.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("accountService.CreateAccount: fetching customer: %w", err)
	}

	// Step 2 — customer must be KYC verified
	if customer.KycStatus != "verified" {
		return nil, domainerrors.ErrKYCNotVerified
	}

	// Step 3 — look up currency
	currency, err := s.currency.GetByCode(ctx, req.CurrencyCode)
	if err != nil {
		return nil, fmt.Errorf("accountService.CreateAccount: invalid currency: %w", err)
	}

	// Step 4 — look up account type
	accountType, err := s.accountType.GetByName(ctx, req.AccountType)
	if err != nil {
		return nil, fmt.Errorf("accountService.CreateAccount: invalid account type: %w", err)
	}

	// Step 5 — create account
	account := &models.Account{
		CustomerID: req.CustomerID,
		CurrencyID: currency.ID,
		TypeID:     accountType.ID,
		CountryID:  customer.CountryID,
		IsActive:   true,
	}

	created, err := s.account.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("accountService.CreateAccount: %w", err)
	}

	s.auditLog.Create(ctx, &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "account",
		EntityID:   created.ID,
		Action:     "ACCOUNT_CREATED",
		ActorID:    req.CustomerID,
		ActorType:  "customer",
	})
	return &CreateAccountResponse{
		AccountID:     created.ID,
		AccountNumber: created.AccountNumber,
		CustomerID:    req.CustomerID,
		CurrencyCode:  req.CurrencyCode,
		AccountType:   req.AccountType,
		Status:        "active",
	}, nil
}
