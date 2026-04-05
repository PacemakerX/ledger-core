package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
)

type TransactionHistoryRequest struct {
	AccountID uuid.UUID
	Limit     int
	Cursor    *uuid.UUID
}

type TransactionItem struct {
	ID                    uuid.UUID  `json:"id"`
	Type                  string     `json:"type"`
	Status                string     `json:"status"`
	Amount                int64      `json:"amount"`
	FromAccountID         *uuid.UUID `json:"from_account_id"`
	ToAccountID           *uuid.UUID `json:"to_account_id"`
	OriginalTransactionID *uuid.UUID `json:"original_transaction_id,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	CompletedAt           *time.Time `json:"completed_at"`
}

type TransactionHistoryResponse struct {
	Transactions []TransactionItem `json:"transactions"`
	NextCursor   *uuid.UUID        `json:"next_cursor"`
	HasMore      bool              `json:"has_more"`
	Count        int               `json:"count"`
}

type transactionService struct {
	account     repository.AccountRepository
	transaction repository.TransactionRepository
}

func NewTransactionService(
	account repository.AccountRepository,
	transaction repository.TransactionRepository,
) *transactionService {
	return &transactionService{
		account:     account,
		transaction: transaction,
	}
}

func (s *transactionService) GetHistory(ctx context.Context, req TransactionHistoryRequest) (*TransactionHistoryResponse, error) {

	// Default limit
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	// Step 1 — verify account exists
	_, err := s.account.GetByID(ctx, req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("transactionService.GetHistory: %w", err)
	}

	// Step 2 — fetch limit+1 rows
	transactions, err := s.transaction.GetByAccountID(ctx, req.AccountID, req.Limit, req.Cursor)
	if err != nil {
		return nil, fmt.Errorf("transactionService.GetHistory: %w", err)
	}

	// Step 3 — determine if there are more pages
	hasMore := len(transactions) > req.Limit
	if hasMore {
		transactions = transactions[:req.Limit]
	}

	// Step 4 — build response items
	items := make([]TransactionItem, 0, len(transactions))
	for _, t := range transactions {
		items = append(items, toTransactionItem(t))
	}

	// Step 5 — set next cursor
	var nextCursor *uuid.UUID
	if hasMore && len(items) > 0 {
		lastID := items[len(items)-1].ID
		nextCursor = &lastID
	}

	return &TransactionHistoryResponse{
		Transactions: items,
		NextCursor:   nextCursor,
		HasMore:      hasMore,
		Count:        len(items),
	}, nil
}

func toTransactionItem(t models.Transaction) TransactionItem {
	return TransactionItem{
		ID:                    t.ID,
		Type:                  t.Type,
		Status:                t.Status,
		Amount:                t.Amount,
		FromAccountID:         t.FromAccountID,
		ToAccountID:           t.ToAccountID,
		OriginalTransactionID: t.OriginalTransactionID,
		CreatedAt:             t.CreatedAt,
		CompletedAt:           t.CompletedAt,
	}
}
