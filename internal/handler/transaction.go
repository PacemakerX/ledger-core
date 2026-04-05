package handler

import (
	"context"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/service"
)

type TransactionService interface {
	GetHistory(ctx context.Context, req service.TransactionHistoryRequest) (*service.TransactionHistoryResponse, error)
}

type transactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *transactionHandler {
	return &transactionHandler{service: service}
}

// HandleGetTransactionHistory returns paginated transaction history for an account
// @Summary      Get transaction history
// @Description  Returns paginated transaction history for an account using cursor-based pagination
// @Tags         accounts
// @Produce      json
// @Param        id      path      string  true   "Account ID"
// @Param        limit   query     int     false  "Number of results (default 20, max 100)"
// @Param        cursor  query     string  false  "Cursor from previous page (last transaction ID)"
// @Success      200     {object}  service.TransactionHistoryResponse
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /accounts/{id}/transactions [get]
func (h *transactionHandler) HandleGetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	requestID := chimiddleware.GetReqID(r.Context())

	// Extract account ID from URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"invalid account_id")
		return
	}

	// Parse limit query param
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed <= 0 {
			domainerrors.WriteError(w, requestID,
				http.StatusBadRequest,
				"LEDGER_400_VALIDATION_ERROR",
				"limit must be a positive integer")
			return
		}
		limit = parsed
	}

	// Parse cursor query param
	var cursor *uuid.UUID
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		parsed, err := uuid.Parse(cursorStr)
		if err != nil {
			domainerrors.WriteError(w, requestID,
				http.StatusBadRequest,
				"LEDGER_400_VALIDATION_ERROR",
				"invalid cursor")
			return
		}
		cursor = &parsed
	}

	req := service.TransactionHistoryRequest{
		AccountID: accountID,
		Limit:     limit,
		Cursor:    cursor,
	}

	response, err := h.service.GetHistory(r.Context(), req)
	if err != nil {
		switch {
		case isNotFound(err):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "account not found")
		default:
			sentry.CaptureException(err)
			domainerrors.WriteError(w, requestID, http.StatusInternalServerError,
				domainerrors.CodeInternalError, "an unexpected error occurred")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
