package handler

import (
	"fmt"

	"github.com/PacemakerX/ledger-core/internal/service"
	"github.com/google/uuid"
)

func validateTransferRequest(req service.TransferRequest) error {
	if req.FromAccountID == uuid.Nil {
		return fmt.Errorf("from_account_id is required")
	}
	if req.ToAccountID == uuid.Nil {
		return fmt.Errorf("to_account_id is required")
	}
	if req.FromAccountID == req.ToAccountID {
		return fmt.Errorf("from_account_id and to_account_id must be different")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if req.IdempotencyKey == "" {
		return fmt.Errorf("idempotency_key is required")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	const MinTransferAmount = int64(100) // 1 rupee minimum
	const PlatformFee = int64(5)         // 5 paise

	if req.Amount < MinTransferAmount {
		return fmt.Errorf("amount must be at least %d paise", MinTransferAmount)
	}
	return nil
}
