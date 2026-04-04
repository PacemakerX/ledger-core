package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/service"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type RefundService interface {
	Refund(ctx context.Context, req service.RefundRequest) (*service.RefundResponse, error)
}
type refundHandler struct {
	service RefundService
}

func NewRefundHandler(service RefundService) *refundHandler {

	return &refundHandler{service: service}
}

func (h *refundHandler) HandleRefund(w http.ResponseWriter, r *http.Request) {
	// 1. decode JSON body into service.TransferRequest
	// 2. call h.service.Transfer
	// 3. handle error → map to HTTP status
	// 4. encode response as JSON
	var req service.RefundRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateRefundRequest(req); err != nil {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			err.Error())
		return
	}

	response, err := h.service.Refund(r.Context(), req)

	if err != nil {
		fmt.Printf("refund error: %v\n", err)
		requestID := chimiddleware.GetReqID(r.Context())
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "requested resource was not found")
		case errors.Is(err, domainerrors.ErrInsufficientBalance):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeInsufficientBalance, "account does not have sufficient balance")
		case errors.Is(err, domainerrors.ErrKYCNotVerified):
			domainerrors.WriteError(w, requestID, http.StatusForbidden,
				domainerrors.CodeKYCNotVerified, "customer KYC verification is not complete")
		case errors.Is(err, domainerrors.ErrAccountInactive):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeAccountInactive, "one or both accounts are inactive")
		case errors.Is(err, domainerrors.ErrDailyLimitExceeded):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeDailyLimitExceeded, "daily transfer limit exceeded")
		case errors.Is(err, domainerrors.ErrMonthlyLimitExceeded):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeMonthlyLimitExceeded, "monthly transfer limit exceeded")
		case errors.Is(err, domainerrors.ErrYearlyLimitExceeded):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeYearlyLimitExceeded, "yearly transfer limit exceeded")
		case errors.Is(err, domainerrors.ErrTransactionLimitExceeded):
			domainerrors.WriteError(w, requestID, http.StatusUnprocessableEntity,
				domainerrors.CodeTransactionLimitExceeded, "amount exceeds per-transaction limit")
		case errors.Is(err, domainerrors.ErrIdempotencyConflict):
			domainerrors.WriteError(w, requestID, http.StatusConflict,
				domainerrors.CodeIdempotencyConflict, "idempotency key reused with different payload")
		default:
			domainerrors.WriteError(w, requestID, http.StatusInternalServerError,
				domainerrors.CodeInternalError, "an unexpected error occurred")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
