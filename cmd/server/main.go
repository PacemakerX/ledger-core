package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PacemakerX/ledger-core/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/PacemakerX/ledger-core/internal/db"
	"go.uber.org/zap"
)

func main() {

	//  Load Config ────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	

	
	//  Initialize Logger ──────────────────────────────────────
	var logger *zap.Logger
	if cfg.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()


	pool, err := db.NewPool(context.Background(), &cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()
	logger.Info("database connected",
		zap.String("host", cfg.Database.Host),
		zap.String("db", cfg.Database.Name),
	)

	//  Setup Router ───────────────────────────────────────────
	r := chi.NewRouter()

	// Core middleware
	r.Use(middleware.RequestID)   // Adds X-Request-Id to every request
	r.Use(middleware.RealIP)      // Uses X-Forwarded-For header
	r.Use(middleware.Recoverer)   // Recovers from panics gracefully
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	//  Routes ─────────────────────────────────────────────────
	r.Get("/health", healthCheck(logger))

	// API v1 group — all ledger routes will go here
	r.Route("/api/v1", func(r chi.Router) {
		
		// coming soon
	})



	//  Start Server ───────────────────────────────────────────
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}



	// Start server in goroutine so it doesn't block shutdown logic
	go func() {
		logger.Info("server starting",
			zap.String("port", cfg.App.Port),
			zap.String("env", cfg.App.Env),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()




	//  Graceful Shutdown ──────────────────────────────────────

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down gracefully...")

	// Give active requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped")
}


// ── Handlers ──────────────────────────────────────────────────────

type healthResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

func healthCheck(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("health check called",
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"ledger-core","version":"0.1.0"}`))
	}
}
