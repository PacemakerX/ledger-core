package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/service"
)

type TransferService interface {
	Transfer(ctx context.Context, req service.TransferRequest) (*service.TransferResponse, error)
}
type transferHandler struct {
	service TransferService
}

func NewTransferHandler(service TransferService) *transferHandler {

	return &transferHandler{service: service}
}

func (h *transferHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	// 1. decode JSON body into service.TransferRequest
	// 2. call h.service.Transfer
	// 3. handle error → map to HTTP status
	// 4. encode response as JSON
	var req service.TransferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.Transfer(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domainerrors.ErrInsufficientBalance):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, domainerrors.ErrKYCNotVerified):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, domainerrors.ErrAccountInactive):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, domainerrors.ErrDailyLimitExceeded),
			errors.Is(err, domainerrors.ErrMonthlyLimitExceeded),
			errors.Is(err, domainerrors.ErrTransactionLimitExceeded):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			fmt.Printf("transfer error: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
