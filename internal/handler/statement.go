package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	domainerrors "github.com/PacemakerX/ledger-core/internal/errors"
)

type StatementService interface {
	GenerateStatement(ctx context.Context, accountID uuid.UUID, w io.Writer) error
}

type statementHandler struct {
	service StatementService
}

func NewStatementHandler(service StatementService) *statementHandler {
	return &statementHandler{service: service}
}

// HandleGetStatement generates a PDF statement for an account
// @Summary      Get account statement
// @Description  Generates a PDF statement for the last 100 transactions
// @Tags         accounts
// @Produce      application/pdf
// @Param        id   path      string  true  "Account ID"
// @Success      200  {file}    binary
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /accounts/{id}/statement [get]
func (h *statementHandler) HandleGetStatement(w http.ResponseWriter, r *http.Request) {
	requestID := chimiddleware.GetReqID(r.Context())

	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		domainerrors.WriteError(w, requestID,
			http.StatusBadRequest,
			"LEDGER_400_VALIDATION_ERROR",
			"invalid account_id")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition",
		"attachment; filename=statement-"+accountID.String()+".pdf")

	err = h.service.GenerateStatement(r.Context(), accountID, w)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrNotFound):
			domainerrors.WriteError(w, requestID, http.StatusNotFound,
				domainerrors.CodeNotFound, "account not found")
		default:
			domainerrors.WriteError(w, requestID, http.StatusInternalServerError,
				domainerrors.CodeInternalError, "failed to generate statement")
		}
		return
	}
}
