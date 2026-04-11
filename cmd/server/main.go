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
	_ "github.com/PacemakerX/ledger-core/docs"
	"github.com/PacemakerX/ledger-core/internal/db"
	"github.com/PacemakerX/ledger-core/internal/handler"
	"github.com/PacemakerX/ledger-core/internal/middleware"
	"github.com/PacemakerX/ledger-core/internal/repository/postgres"
	"github.com/PacemakerX/ledger-core/internal/service"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title           ledger-core API
// @version         1.0
// @description     Production-grade double-entry accounting ledger in Go and PostgreSQL
// @host            localhost:8080
// @BasePath        /api/v1
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

	if cfg.App.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.App.SentryDSN,
			Environment:      cfg.App.Env,
			Release:          cfg.App.Version,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			logger.Fatal("failed to initialize sentry", zap.Error(err))
		}
		defer sentry.Flush(2 * time.Second)
		logger.Info("sentry initialized")
	}

	// Repositories
	accountRepo := postgres.NewAccountRepository(pool)
	transactionRepo := postgres.NewTransactionRepository(pool)
	journalRepo := postgres.NewJournalEntryRepository(pool)
	idempotencyRepo := postgres.NewIdempotencyRepository(pool)
	customerRepo := postgres.NewCustomerRepository(pool)
	accountLimitRepo := postgres.NewAccountLimitRepository(pool)
	txManager := postgres.NewTxManager(pool)
	countryRepo := postgres.NewCountryRepository(pool)
	currencyRepo := postgres.NewCurrencyRepository(pool)
	accountTypeRepo := postgres.NewAccountTypeRepository(pool)
	auditLogRepo := postgres.NewAuditLogRepository(pool)

	// Services
	transferSvc := service.NewTransferService(
		txManager,
		accountRepo,
		transactionRepo,
		journalRepo,
		idempotencyRepo,
		customerRepo,
		accountLimitRepo,
		auditLogRepo,
		cfg,
	)
	refundSvc := service.NewRefundService(
		txManager,
		accountRepo,
		transactionRepo,
		journalRepo,
		idempotencyRepo,
		customerRepo,
		accountLimitRepo,
		auditLogRepo,
		cfg,
	)
	customerSvc := service.NewCustomerService(customerRepo, countryRepo,auditLogRepo)
	accountSvc := service.NewAccountService(customerRepo, accountRepo, currencyRepo, accountTypeRepo)
	transactionSvc := service.NewTransactionService(accountRepo, transactionRepo)
	statementSvc := service.NewStatementService(accountRepo, transactionRepo, journalRepo)

	// Handlers
	transferHandler := handler.NewTransferHandler(transferSvc)
	refundHandler := handler.NewRefundHandler(refundSvc)
	healthHandler := handler.NewHealthHandler(
		logger,
		"ledger-core",
		cfg.App.Env,
		cfg.App.Version,
		pool,
	)
	customerHandler := handler.NewCustomerHandler(customerSvc)
	accountHandler := handler.NewAccountHandler(accountSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)
	statementHandler := handler.NewStatementHandler(statementSvc)
	//  Setup Router ───────────────────────────────────────────
	r := chi.NewRouter()

	// Core middleware
	r.Use(chimiddleware.RequestID)                 // Adds X-Request-Id to every request
	r.Use(chimiddleware.RealIP)                    // Uses X-Forwarded-For header
	r.Use(chimiddleware.Recoverer)                 // Recovers from panics gracefully
	r.Use(chimiddleware.Timeout(60 * time.Second)) // Request timeout
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.RequestLogger(logger))

	//  Routes ─────────────────────────────────────────────────
	r.Get("/health", healthHandler.ServeHTTP)

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API v1 group — all ledger routes will go here
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(httprate.LimitByIP(cfg.App.RateLimitRequests, time.Duration(cfg.App.RateLimitWindow)*time.Second))

		// Transfer
		// @Summary Create transfer
		// @Description Creates a new funds transfer
		// @Tags transfers
		// @Accept json
		// @Produce json
		// @Success 201 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 409 {object} map[string]string
		// @Router /api/v1/transfers [post]
		r.Post("/transfers", transferHandler.HandleTransfer)

		// Refund
		// @Summary Create refund
		// @Description Creates a refund for a previous transfer
		// @Tags refunds
		// @Accept json
		// @Produce json
		// @Success 201 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 409 {object} map[string]string
		// @Router /api/v1/refunds [post]
		r.Post("/refunds", refundHandler.HandleRefund)

		// Customer routes
		// @Summary Create customer
		// @Description Creates a new customer
		// @Tags customers
		// @Accept json
		// @Produce json
		// @Success 201 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Router /api/v1/customers [post]
		r.Post("/customers", customerHandler.HandleCreateCustomer)

		// @Summary Update customer KYC
		// @Description Updates customer KYC information
		// @Tags customers
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /api/v1/customers/{id}/kyc [patch]
		r.Patch("/customers/{id}/kyc", customerHandler.HandleUpdateKYC)

		// @Summary Create account
		// @Description Creates a new account
		// @Tags accounts
		// @Accept json
		// @Produce json
		// @Success 201 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Router /api/v1/accounts [post]
		r.Post("/accounts", accountHandler.HandleCreateAccount)

		r.Get("/accounts/{id}/transactions", transactionHandler.HandleGetTransactionHistory)

		r.Get("/accounts/{id}/statement", statementHandler.HandleGetStatement)
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
