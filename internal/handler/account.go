package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/service"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req service.CreateAccountRequest) (*service.CreateAccountResponse, error)
}

type accountHandler struct {
	service AccountService
}

func NewAccountHandler(service AccountService) *accountHandler {
	return &accountHandler{service: service}
}

// HandleCreateAccount opens a new account for a verified customer
// @Summary      Create account
// @Description  Opens a new account for a KYC-verified customer. One account per customer per currency per type.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body service.CreateAccountRequest true "Account details"
// @Success      201 {object} service.CreateAccountResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts [post]
func (h *accountHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	requestID := chimiddleware.GetReqID(r.Context())

	var req service.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_INVALID_REQUEST",
			"invalid request body")
		return
	}

	// Validation
	if req.CustomerID == uuid.Nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"customer_id is required")
		return
	}
	if req.CurrencyCode == "" {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"currency_code is required")
		return
	}
	if req.AccountType == "" {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"account_type is required")
		return
	}

	response, err := h.service.CreateAccount(r.Context(), req)
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "customer, currency or account type not found")
		case errors.Is(err, domainerrors.ErrKYCNotVerified):
			domainerrors.WriteError(w, requestID, http.StatusForbidden,
				domainerrors.CodeKYCNotVerified, "customer KYC is not verified")
		case errors.Is(err, domainerrors.ErrAlreadyExists):
			domainerrors.WriteError(w, requestID, http.StatusConflict,
				domainerrors.CodeIdempotencyConflict, "account already exists for this customer, currency and type")
		default:
			sentry.CaptureException(err)
			domainerrors.WriteError(w, requestID, http.StatusInternalServerError,
				domainerrors.CodeInternalError, "an unexpected error occurred")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
