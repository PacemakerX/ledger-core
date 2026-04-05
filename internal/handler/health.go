package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type HealthHandler struct {
	Logger    *zap.Logger
	Service   string
	Env       string
	Version   string
	StartTime time.Time
	Pool      *pgxpool.Pool // add this
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Env       string            `json:"env"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks,omitempty"`
}

func NewHealthHandler(logger *zap.Logger, service, env, version string, pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		Logger:    logger,
		Service:   service,
		Env:       env,
		Version:   version,
		StartTime: time.Now(),
		Pool:      pool,
	}
}

// ServeHTTP returns service and database health
// @Summary      Health check
// @Description  Returns service status, uptime, version, and database connectivity
// @Tags         health
// @Produce      json
// @Success      200 {object} handler.HealthResponse
// @Failure      503 {object} handler.HealthResponse
// @Router       /health [get]
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqID := chimiddleware.GetReqID(r.Context())

	h.Logger.Info("health check called",
		zap.String("request_id", reqID),
	)

	// Basic liveness check (always OK if server is running)
	checks := map[string]string{
		"server": "ok",
	}
	status := "ok"

	// Check database FIRST
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.Pool.Ping(ctx); err != nil {
		checks["database"] = "unavailable"
		status = "degraded"
	} else {
		checks["database"] = "ok"
	}

	// THEN build response
	uptime := time.Since(h.StartTime).String()
	resp := HealthResponse{
		Status:    status,
		Service:   h.Service,
		Env:       h.Env,
		Version:   h.Version,
		Timestamp: time.Now().UTC(),
		Uptime:    uptime,
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")

	if status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error("failed to encode health response",
			zap.String("request_id", reqID),
			zap.Error(err),
		)
	}
}
