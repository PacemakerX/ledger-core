package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PacemakerX/ledger-core/config"
	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
)

type TransferRequest struct {
	FromAccountID  uuid.UUID `json:"from_account_id"`
	ToAccountID    uuid.UUID `json:"to_account_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	IdempotencyKey string    `json:"idempotency_key"`
	InitiatedBy    uuid.UUID `json:"initiated_by"`
	Description    *string   `json:"description"`
}

type TransferResponse struct {
	TransactionID uuid.UUID
	Status        string
	CreatedAt     time.Time
}

type transferService struct {
	txManager    repository.TxManager
	account      repository.AccountRepository
	transaction  repository.TransactionRepository
	journal      repository.JournalEntryRepository
	idempotency  repository.IdempotencyRepository
	customer     repository.CustomerRepository
	accountLimit repository.AccountLimitRepository
	cfg          *config.Config
}

func NewTransferService(
	txManager repository.TxManager,
	account repository.AccountRepository,
	transaction repository.TransactionRepository,
	journal repository.JournalEntryRepository,
	idempotency repository.IdempotencyRepository,
	customer repository.CustomerRepository,
	accountLimit repository.AccountLimitRepository,
	cfg *config.Config,
) *transferService {
	return &transferService{
		txManager:    txManager,
		account:      account,
		transaction:  transaction,
		journal:      journal,
		idempotency:  idempotency,
		customer:     customer,
		accountLimit: accountLimit,
		cfg:          cfg,
	}
}

/*
	Steps for a transaction

1.  Check idempotency key
2.  Validate accounts exist
3.  Check KYC status
4.  Check account limits
5.  Check sufficient balance
6.  BEGIN transaction
7.  SELECT FOR UPDATE (lower UUID first)
8.  Create transaction record (PENDING)
9.  Create 4 journal entries
10. Verify SUM = 0
11. Update transaction (COMPLETED)
12. Update limit usage
13. Set idempotency response
14. COMMIT
*/
func (s *transferService) Transfer(ctx context.Context, req TransferRequest) (*TransferResponse, error) {

	// Step 1 — check idempotency key
	existing, err := s.idempotency.Get(ctx, req.IdempotencyKey)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, fmt.Errorf("transfer: checking idempotency: %w", err)
	}

	// if key exists, this is a duplicate request — return cached response
	if existing != nil {
		txID, err := uuid.Parse(existing.ResponseBody)
		if err != nil {
			txID = uuid.Nil
		}
		return &TransferResponse{
			TransactionID: txID,
			Status:        existing.ResponseStatus,
			CreatedAt:     existing.CreatedAt,
		}, nil
	}

	// Verifying weather the sender account exists or not
	fromAccount, err := s.account.GetByID(ctx, req.FromAccountID)

	if err != nil {
		return nil, fmt.Errorf("tranfer: getting from account %w", err)

	}

	// Verifying weather the receiver account exists or not
	toAccount, err := s.account.GetByID(ctx, req.ToAccountID)

	if err != nil {
		return nil, fmt.Errorf("tranfer: getting to account %w", err)

	}

	// Verifying weather the sender account is Active or not
	if !fromAccount.IsActive {
		return nil, domainerrors.ErrAccountInactive
	}

	// Verifying weather the receiver account is Active or not
	if !toAccount.IsActive {
		return nil, domainerrors.ErrAccountInactive
	}

	// Verifying weather the sender customer exists or not
	fromCustomer, err := s.customer.GetByID(ctx, fromAccount.CustomerID)

	if err != nil {
		return nil, fmt.Errorf("transfer: from customer %w", err)
	}

	if fromCustomer.KycStatus != "verified" {

		return nil, domainerrors.ErrKYCNotVerified
	}

	// Verifying weather the receiver customer exists or not
	toCustomer, err := s.customer.GetByID(ctx, toAccount.CustomerID)

	if err != nil {
		return nil, fmt.Errorf("transfer: to customer %w", err)
	}

	if toCustomer.KycStatus != "verified" {

		return nil, domainerrors.ErrKYCNotVerified
	}

	// Checking for the daily,monthly, yearly and transaction limits
	fromAccountLimit, err := s.accountLimit.GetByAccountID(ctx, fromAccount.ID)

	if err != nil {
		return nil, fmt.Errorf("transfer: from account_limit %w", err)
	}

	for _, limit := range fromAccountLimit {

		switch limit.LimitType {

		case "DAILY":
			{

				if req.Amount+limit.CurrentUsage > limit.MaxAmount {
					return nil, domainerrors.ErrDailyLimitExceeded
				}
			}
		case "MONTHLY":
			{

				if req.Amount+limit.CurrentUsage > limit.MaxAmount {
					return nil, domainerrors.ErrMonthlyLimitExceeded
				}
			}
		case "YEARLY":
			{

				if req.Amount+limit.CurrentUsage > limit.MaxAmount {
					return nil, domainerrors.ErrYearlyLimitExceeded
				}
			}
		case "TRANSACTION":
			{

				if req.Amount > limit.MaxAmount {
					return nil, domainerrors.ErrTransactionLimitExceeded
				}
			}

		}
	}

	// We don't need the limit here

	// toAccountLimit,err:=s.accountLimit.GetByAccountID(ctx,toAccount.ID)

	// if err!=nil{
	//     return nil,fmt.Errorf("transfer: to account_limit %w",err)
	// }

	// for _,limit := range toAccountLimit{

	//     switch limit.LimitType{

	//     case "DAILY":{

	//         if (req.Amount + limit.CurrentUsage > limit.MaxAmount){
	//             return nil,domainerrors.ErrDailyLimitExceeded
	//         }
	//     }
	//     case "MONTHLY":{

	//         if (req.Amount + limit.CurrentUsage > limit.MaxAmount){
	//             return nil,domainerrors.ErrMonthlyLimitExceeded
	//         }
	//     }
	//     case "YEARLY":{

	//         if (req.Amount + limit.CurrentUsage > limit.MaxAmount){
	//             return nil,domainerrors.ErrYearlyLimitExceeded
	//         }
	//     }
	//     case "TRANSACTION":{

	//        if req.Amount > limit.MaxAmount {
	//             return nil, domainerrors.ErrTransactionLimitExceeded
	//         }
	//     }

	//     }
	// }

	// Checking the sender balance
	fromBalance, err := s.journal.GetBalance(ctx, req.FromAccountID)

	if err != nil {
		return nil, fmt.Errorf("transfer: checking balance %w", err)
	}

	if req.Amount > fromBalance {
		return nil, domainerrors.ErrInsufficientBalance
	}

	tx, err := s.txManager.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("transfer: beginning transaction %w", err)
	}

	/*
	 * The defer tx.Rollback(ctx) is important — if anything fails after this point, the transaction automatically rolls back.
	 * If we successfully commit, Rollback becomes a no-op.
	 */
	defer tx.Rollback(ctx)

	err = s.idempotency.Create(ctx, tx, &models.IdempotencyKey{
		Key:            req.IdempotencyKey,
		RequestHash:    req.IdempotencyKey,
		ResponseStatus: "PENDING",
		ResponseBody:   "",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return nil, fmt.Errorf("transfer: creating idempotency key: %w", err)
	}

	first, second := req.FromAccountID, req.ToAccountID

	if req.ToAccountID.String() < req.FromAccountID.String() {
		first, second = req.ToAccountID, req.FromAccountID

	}

	_, err = s.account.GetByIDForUpdate(ctx, tx, first)
	if err != nil {
		return nil, fmt.Errorf("transfer: locking first account: %w", err)
	}
	_, err = s.account.GetByIDForUpdate(ctx, tx, second)

	if err != nil {
		return nil, fmt.Errorf("transfer: locking second account: %w", err)
	}
	if req.InitiatedBy == uuid.Nil {
		req.InitiatedBy = fromAccount.CustomerID
	}
	transaction := &models.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		Type:           "TRANSFER",
		Status:         "PENDING",
		InitiatedBy:    req.InitiatedBy,
		FromAccountID:  &req.FromAccountID,
		ToAccountID:    &req.ToAccountID,
		Amount:         req.Amount,
		CurrencyID:     fromAccount.CurrencyID,
	}

	createdTx, err := s.transaction.Create(ctx, tx, transaction)

	if err != nil {
		return nil, fmt.Errorf("transfer: creating transaction: %w", err)
	}

	// Entrying into the journal
	entries := []models.JournalEntry{
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     req.FromAccountID,
			EntryType:     "CREDIT",
			Amount:        req.Amount,
			CurrencyID:    fromAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     s.cfg.PlatformFloatAccountID,
			EntryType:     "DEBIT",
			Amount:        req.Amount,
			CurrencyID:    fromAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     s.cfg.PlatformFloatAccountID,
			EntryType:     "CREDIT",
			Amount:        req.Amount,
			CurrencyID:    fromAccount.CurrencyID,
		},
		{
			ID:            uuid.New(),
			TransactionID: createdTx.ID,
			AccountID:     req.ToAccountID,
			EntryType:     "DEBIT",
			Amount:        req.Amount,
			CurrencyID:    fromAccount.CurrencyID,
		},
	}

	err = s.journal.CreateBatch(ctx, tx, entries)

	if err != nil {
		return nil, fmt.Errorf("transfer: creating journal entries: %w", err)
	}

	/* Balance Verification
	* We need to verify the balance of one transaction no matter how many parties are involved
	 */
	balance, err := s.journal.VerifyBalance(ctx, tx, createdTx.ID)
	if err != nil {
		return nil, fmt.Errorf("transfer: balance verification: %w", err)
	}
	if balance != 0 {
		return nil, domainerrors.ErrBalanceVerificationFailed
	}

	// updation of transaction status status

	err = s.transaction.UpdateStatus(ctx, tx, createdTx.ID, "COMPLETED")

	if err != nil {
		return nil, fmt.Errorf("transfer: updating transaction status: %w", err)
	}

	// Updating Limit Usage

	for _, limit := range fromAccountLimit {
		err = s.accountLimit.UpdateUsage(ctx, tx, limit.ID, req.Amount)
		if err != nil {
			return nil, fmt.Errorf("transfer: updating limit usage: %w", err)
		}
	}

	// Updating the status of idempotency key
	err = s.idempotency.SetResponse(ctx, tx, req.IdempotencyKey, "COMPLETED", createdTx.ID.String())

	if err != nil {

		return nil, fmt.Errorf("transfer: status update of idempotency key %w", err)
	}

	// Commiting
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("transfer: committing transaction: %w", err)
	}

	return &TransferResponse{
		TransactionID: createdTx.ID,
		Status:        "COMPLETED",
		CreatedAt:     createdTx.CreatedAt,
	}, nil
}
