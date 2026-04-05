package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
	"github.com/PacemakerX/ledger-core/internal/service"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, req service.CreateCustomerRequest) (*service.CreateCustomerResponse, error)
	UpdateKYC(ctx context.Context, id uuid.UUID, req service.UpdateKYCRequest) error
}

type customerHandler struct {
	service CustomerService
}

func NewCustomerHandler(service CustomerService) *customerHandler {
	return &customerHandler{service: service}
}

// HandleCreateCustomer registers a new customer
// @Summary      Create customer
// @Description  Registers a new customer with unverified KYC status
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        request body service.CreateCustomerRequest true "Customer details"
// @Success      201 {object} service.CreateCustomerResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /customers [post]
func (h *customerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req service.CreateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_INVALID_REQUEST",
			"invalid request body")
		return
	}

	// Basic validation
	if req.FirstName == "" || req.LastName == "" {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"first_name and last_name are required")
		return
	}
	if req.Email == "" {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"email is required")
		return
	}
	if req.PhoneNumber == "" {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"phone_number is required")
		return
	}
	if req.CountryCode == "" {
		domainerrors.WriteError(w, chimiddleware.GetReqID(r.Context()),
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"country_code is required")
		return
	}

	response, err := h.service.CreateCustomer(r.Context(), req)
	if err != nil {
		requestID := chimiddleware.GetReqID(r.Context())
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "country not found")
		case errors.Is(err, domainerrors.ErrAlreadyExists):
			domainerrors.WriteError(w, requestID, http.StatusConflict,
				domainerrors.CodeIdempotencyConflict, "customer already exists")
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

func (h *customerHandler) HandleUpdateKYC(w http.ResponseWriter, r *http.Request) {
	requestID := chimiddleware.GetReqID(r.Context())

	// Extract customer ID from URL
	customerIDStr := chi.URLParam(r, "id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"invalid customer_id")
		return
	}

	var req service.UpdateKYCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_INVALID_REQUEST",
			"invalid request body")
		return
	}

	err = h.service.UpdateKYC(r.Context(), customerID, req)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "customer not found")
		default:
			sentry.CaptureException(err)
			domainerrors.WriteError(w, requestID, http.StatusInternalServerError,
				domainerrors.CodeInternalError, "an unexpected error occurred")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "kyc updated successfully",
	})
}
