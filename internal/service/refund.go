package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PacemakerX/ledger-core/config"
	"github.com/PacemakerX/ledger-core/internal/cache"
	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
)

type RefundRequest struct {
	TransactionID  uuid.UUID `json:"transactionID"`
	Amount         int64     `json:"amount"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type RefundResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type refundService struct {
	txManager        repository.TxManager
	account          repository.AccountRepository
	transaction      repository.TransactionRepository
	journal          repository.JournalEntryRepository
	idempotency      repository.IdempotencyRepository
	customer         repository.CustomerRepository
	accountLimit     repository.AccountLimitRepository
	auditLog         repository.AuditLogRepository
	idempotencyCache cache.IdempotencyCache
	cfg              *config.Config
}

func NewRefundService(
	txManager repository.TxManager,
	account repository.AccountRepository,
	transaction repository.TransactionRepository,
	journal repository.JournalEntryRepository,
	idempotency repository.IdempotencyRepository,
	customer repository.CustomerRepository,
	accountLimit repository.AccountLimitRepository,
	auditLog repository.AuditLogRepository,
	idempotencyCache cache.IdempotencyCache,
	cfg *config.Config,
) *refundService {
	return &refundService{
		txManager:        txManager,
		account:          account,
		transaction:      transaction,
		journal:          journal,
		idempotency:      idempotency,
		customer:         customer,
		accountLimit:     accountLimit,
		auditLog:         auditLog,
		idempotencyCache: idempotencyCache,
		cfg:              cfg,
	}
}

func (s *refundService) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {

	// Step 1 — check idempotency key (Redis first, Postgres fallback)
	existing, err := s.idempotencyCache.Get(ctx, req.IdempotencyKey)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		// Redis down — fall through to Postgres silently
		existing, err = s.idempotency.Get(ctx, req.IdempotencyKey)
		if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
			return nil, fmt.Errorf("refund: checking idempotency: %w", err)
		}
	}

	if existing != nil {
		txID, err := uuid.Parse(existing.ResponseBody)
		if err != nil {
			txID = uuid.Nil
		}
		return &RefundResponse{
			TransactionID: txID,
			Status:        existing.ResponseStatus,
			CreatedAt:     existing.CreatedAt,
		}, nil
	}

	// Step 2 — fetch and validate original transaction
	original, err := s.transaction.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("refund: fetching original transaction: %w", err)
	}

	// Step 3 — must be COMPLETED
	if original.Status != "COMPLETED" {
		return nil, fmt.Errorf("refund: original transaction is not completed")
	}

	// Step 4 — must be TRANSFER type
	if original.Type != "TRANSFER" {
		return nil, fmt.Errorf("refund: only TRANSFER transactions can be refunded")
	}

	// Step 5 — must be within 90 days
	if original.CompletedAt == nil {
		return nil, fmt.Errorf("refund: original transaction has no completion time")
	}
	if time.Now().After(original.CompletedAt.AddDate(0, 0, 90)) {
		return nil, fmt.Errorf("refund: refund window expired (90 days)")
	}

	// Step 6 — validate refund amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("refund: amount must be greater than zero")
	}
	if req.Amount > original.Amount {
		return nil, fmt.Errorf("refund: amount exceeds original transaction amount")
	}

	// Step 7 — check total already refunded (partial refund guard)
	totalRefunded, err := s.transaction.GetTotalRefunded(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("refund: checking total refunded: %w", err)
	}
	if req.Amount+totalRefunded > original.Amount {
		return nil, fmt.Errorf("refund: total refunded would exceed original amount (already refunded: %d)", totalRefunded)
	}

	// Step 8 — fetch receiver account (original receiver is now the one giving money back)
	// In original transfer: fromAccount sent, toAccount received
	// In refund: toAccount gives back, fromAccount receives
	toAccount, err := s.account.GetByID(ctx, *original.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("refund: fetching receiver account: %w", err)
	}

	// Step 9 — check receiver (refund sender) has sufficient balance
	receiverBalance, err := s.journal.GetBalance(ctx, *original.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("refund: checking receiver balance: %w", err)
	}
	if req.Amount > receiverBalance {
		return nil, domainerrors.ErrInsufficientBalance
	}

	// Step 10 — fetch account limits for receiver
	receiverLimits, err := s.accountLimit.GetByAccountID(ctx, *original.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("refund: fetching receiver limits: %w", err)
	}
	for _, limit := range receiverLimits {
		switch limit.LimitType {
		case "DAILY":
			if req.Amount+limit.CurrentUsage > limit.MaxAmount {
				return nil, domainerrors.ErrDailyLimitExceeded
			}
		case "MONTHLY":
			if req.Amount+limit.CurrentUsage > limit.MaxAmount {
				return nil, domainerrors.ErrMonthlyLimitExceeded
			}
		case "YEARLY":
			if req.Amount+limit.CurrentUsage > limit.MaxAmount {
				return nil, domainerrors.ErrYearlyLimitExceeded
			}
		case "TRANSACTION":
			if req.Amount > limit.MaxAmount {
				return nil, domainerrors.ErrTransactionLimitExceeded
			}
		}
	}

	// Step 11 — BEGIN transaction
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("refund: beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Step 12 — create idempotency key inside tx
	err = s.idempotency.Create(ctx, tx, &models.IdempotencyKey{
		Key:            req.IdempotencyKey,
		RequestHash:    req.IdempotencyKey,
		ResponseStatus: "PENDING",
		ResponseBody:   "",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return nil, fmt.Errorf("refund: creating idempotency key: %w", err)
	}

	// Step 13 — SELECT FOR UPDATE (lower UUID first — deadlock prevention)
	first, second := *original.FromAccountID, *original.ToAccountID
	if second.String() < first.String() {
		first, second = second, first
	}
	_, err = s.account.GetByIDForUpdate(ctx, tx, first)
	if err != nil {
		return nil, fmt.Errorf("refund: locking first account: %w", err)
	}
	_, err = s.account.GetByIDForUpdate(ctx, tx, second)
	if err != nil {
		return nil, fmt.Errorf("refund: locking second account: %w", err)
	}

	// Step 14 — create REFUND transaction record
	refundTx := &models.Transaction{
		ID:                    uuid.New(),
		IdempotencyKey:        req.IdempotencyKey,
		Type:                  "REFUND",
		Status:                "PENDING",
		InitiatedBy:           *original.ToAccountID,
		FromAccountID:         original.ToAccountID,
		ToAccountID:           original.FromAccountID,
		Amount:                req.Amount,
		CurrencyID:            toAccount.CurrencyID,
		OriginalTransactionID: &req.TransactionID,
	}
	createdTx, err := s.transaction.Create(ctx, tx, refundTx)
	if err != nil {
		return nil, fmt.Errorf("refund: creating transaction: %w", err)
	}

	// Step 15 — create 4 journal entries (reversed from original transfer)
	entries := []models.JournalEntry{
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     *original.ToAccountID, // original receiver — money leaving
			EntryType:     "CREDIT",
			Amount:        req.Amount,
			CurrencyID:    toAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     s.cfg.PlatformFloatAccountID, // platform receives
			EntryType:     "DEBIT",
			Amount:        req.Amount,
			CurrencyID:    toAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     s.cfg.PlatformFloatAccountID, // platform releases
			EntryType:     "CREDIT",
			Amount:        req.Amount,
			CurrencyID:    toAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     *original.FromAccountID, // original sender — money arriving
			EntryType:     "DEBIT",
			Amount:        req.Amount,
			CurrencyID:    toAccount.CurrencyID,
		},
	}
	err = s.journal.CreateBatch(ctx, tx, entries)
	if err != nil {
		return nil, fmt.Errorf("refund: creating journal entries: %w", err)
	}

	// Step 16 — verify SUM(debits) - SUM(credits) = 0
	balance, err := s.journal.VerifyBalance(ctx, tx, createdTx.ID)
	if err != nil {
		return nil, fmt.Errorf("refund: balance verification: %w", err)
	}
	if balance != 0 {
		return nil, domainerrors.ErrBalanceVerificationFailed
	}

	// Step 17 — update transaction COMPLETED
	err = s.transaction.UpdateStatus(ctx, tx, createdTx.ID, "COMPLETED")
	if err != nil {
		return nil, fmt.Errorf("refund: updating transaction status: %w", err)
	}

	// Step 18 — update limit usage
	for _, limit := range receiverLimits {
		err = s.accountLimit.UpdateUsage(ctx, tx, limit.ID, req.Amount)
		if err != nil {
			return nil, fmt.Errorf("refund: updating limit usage: %w", err)
		}
	}

	// Step 19 — set idempotency response
	err = s.idempotency.SetResponse(ctx, tx, req.IdempotencyKey, "COMPLETED", createdTx.ID.String())
	if err != nil {
		return nil, fmt.Errorf("refund: setting idempotency response: %w", err)
	}

	// Step 20 — COMMIT
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("refund: committing transaction: %w", err)
	}

	// Populate Redis cache after successful commit
	s.idempotencyCache.Set(ctx, &models.IdempotencyKey{
		Key:            req.IdempotencyKey,
		ResponseStatus: "COMPLETED",
		ResponseBody:   createdTx.ID.String(),
		CreatedAt:      createdTx.CreatedAt,
	})

	// Fire-and-forget audit log — after commit, outside transaction
	s.auditLog.Create(ctx, &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "transaction",
		EntityID:   createdTx.ID,
		Action:     "REFUND_COMPLETED",
		ActorID:    *original.FromAccountID,
		ActorType:  "customer",
		CreatedAt:  time.Now(),
	})

	return &RefundResponse{
		TransactionID: createdTx.ID,
		Status:        "COMPLETED",
		CreatedAt:     createdTx.CreatedAt,
	}, nil
}
