# ledger-core

A production-grade double-entry accounting ledger built in Go and PostgreSQL.
Handles concurrent transfers with full ACID guarantees, idempotency, and real-time observability.

---
```bash
.
├── cmd
│   └── server
│       └── main.go
├── config
│   └── config.go
├── docs
│   ├── adr
│   │   ├── 000-template.md
│   │   ├── 001-why-go.md
│   │   ├── 002-why-postgreSQL.md
│   │   ├── 003-why-pgx.md
│   │   ├── 004-why-chi.md
│   │   ├── 005-why-zap.md
│   │   ├── 006-why-uuid.md
│   │   ├── 007-why-golang-migrate-over-gorm-automigrate.md
│   │   ├── 008-why-double-entry-accounting.md
│   │   └── 009-original-transaction-id.md
│   ├── images
│   │   ├── grafana.jpg
│   │   └── load_test.png
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── db
│   │   └── postgres.go
│   ├── errors
│   │   ├── errors.go
│   │   └── response.go
│   ├── handler
│   │   ├── account.go
│   │   ├── customer.go
│   │   ├── health.go
│   │   ├── refund.go
│   │   ├── transfer.go
│   │   └── validate.go
│   ├── metrics
│   │   └── metrics.go
│   ├── middleware
│   │   └── metrics_middleware.go
│   ├── models
│   │   ├── account.go
│   │   ├── account_limit.go
│   │   ├── account_type.go
│   │   ├── audit_log.go
│   │   ├── country.go
│   │   ├── currency.go
│   │   ├── customer.go
│   │   ├── exchange_rate.go
│   │   ├── idempotency_key.go
│   │   ├── journal_entry.go
│   │   └── transaction.go
│   ├── repository
│   │   ├── postgres
│   │   │   ├── account_limit_repository.go
│   │   │   ├── account_repository.go
│   │   │   ├── account_type_repository.go
│   │   │   ├── country_repository.go
│   │   │   ├── currency_repository.go
│   │   │   ├── customer_repository.go
│   │   │   ├── idempotency_repository.go
│   │   │   ├── journal_entry_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   └── tx_manager.go
│   │   └── interfaces.go
│   └── service
│       ├── account.go
│       ├── customer.go
│       ├── refund.go
│       └── transfer.go
├── migrations
│   ├── 000001_create_currencies.down.sql
│   ├── 000001_create_currencies.up.sql
│   ├── 000002_create_exchange_rates.down.sql
│   ├── 000002_create_exchange_rates.up.sql
│   ├── 000003_create_account_types.down.sql
│   ├── 000003_create_account_types.up.sql
│   ├── 000004_create_countries.down.sql
│   ├── 000004_create_countries.up.sql
│   ├── 000005_create_customers.down.sql
│   ├── 000005_create_customers.up.sql
│   ├── 000006_create_accounts.down.sql
│   ├── 000006_create_accounts.up.sql
│   ├── 000007_create_transactions.down.sql
│   ├── 000007_create_transactions.up.sql
│   ├── 000008_create_journal_entries.down.sql
│   ├── 000008_create_journal_entries.up.sql
│   ├── 000009_create_idempotency_keys.down.sql
│   ├── 000009_create_idempotency_keys.up.sql
│   ├── 000010_create_audit_logs.down.sql
│   ├── 000010_create_audit_logs.up.sql
│   ├── 000011_create_account_limits.down.sql
│   ├── 000011_create_account_limits.up.sql
│   ├── 000012_create_indexes.down.sql
│   ├── 000012_create_indexes.up.sql
│   ├── 000013_seed_platform_accounts.up.sql
│   ├── 000014_fix_idempotency_keys.down.sql
│   ├── 000014_fix_idempotency_keys.up.sql
│   ├── 000015_add_account_ids_to_transactions.down.sql
│   ├── 000015_add_account_ids_to_transactions.up.sql
│   ├── 000016_add_amount_to_transactions.down.sql
│   ├── 000016_add_amount_to_transactions.up.sql
│   ├── 000017_add_original_transaction_id.down.sql
│   └── 000017_add_original_transaction_id.up.sql
├── scripts
│   └── loadtest
│       └── k6.js
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── prometheus.yml
└── README.md

20 directories, 96 files
```

---

## How a Transfer Works

Every transfer executes these 14 steps atomically:

1. Check idempotency key — return cached response if duplicate request
2. Fetch sender account — verify it exists
3. Fetch receiver account — verify it exists
4. Check both accounts are active
5. Fetch sender customer — verify KYC = `verified`
6. Fetch receiver customer — verify KYC = `verified`
7. Check sender account limits — DAILY, MONTHLY, YEARLY, TRANSACTION
8. Verify sufficient balance (derived from journal entries via `SUM`, never stored)
9. `BEGIN` transaction
10. `defer tx.Rollback` — automatic rollback on any failure
11. Create idempotency key (`PENDING`) inside the transaction
12. `SELECT FOR UPDATE` both accounts (lower UUID first — deadlock prevention)
13. Create transaction record (`PENDING`)
14. `CreateBatch` — 4 journal entries:

| Account        | Entry  | Effect            |
| -------------- | ------ | ----------------- |
| Sender         | CREDIT | money leaving     |
| Platform Float | DEBIT  | platform receives |
| Platform Float | CREDIT | platform releases |
| Receiver       | DEBIT  | money arriving    |

15. Verify `SUM(debits) - SUM(credits) = 0` — mathematically enforced
16. Update transaction → `COMPLETED`
17. Update limit usage
18. Set idempotency response → `COMPLETED`
19. `COMMIT`

If any step fails, the entire transaction rolls back automatically.

---

## Key Design Decisions

**Why `SELECT FOR UPDATE` with lock ordering?**
Concurrent transfers between the same accounts cause lost updates without row-level locks.
Always locking the lower UUID first across all code paths prevents deadlocks.

**Why derive balance from journal entries?**
Storing balance as a column creates a TOCTOU race condition under concurrency.
Deriving it from immutable journal entries means the ledger is always the source of truth.

```sql
SELECT COALESCE(
    SUM(CASE WHEN entry_type = 'DEBIT' THEN amount ELSE -amount END), 0
)
FROM journal_entries WHERE account_id = $1
```

**Why idempotency keys inside the database transaction?**
If the transaction rolls back, the idempotency key rolls back with it.
This prevents a failed transfer from being treated as already-processed on retry.

**Why UUIDv7 over v4?**
UUIDv7 is time-ordered — sequential inserts cause less B-tree fragmentation in PostgreSQL indexes.

**Why `BIGINT` for money?**
Floating point arithmetic is non-deterministic. `BIGINT` in smallest currency unit (paise for INR, cents for USD) is exact.

**Why append-only journal entries?**
Financial ledgers must be auditable. No `UPDATE` or `DELETE` ever touches `journal_entries`.
Corrections are made via new entries (refunds, adjustments), never by editing history.

---

## Database Schema

14 migrations, all append-only:

```bash
001 currencies               — INR, USD, SGD (50 seeded)
002 exchange_rates           — NUMERIC(20,8), daily rates
003 account_types            — asset/liability/equity/revenue/expense
004 countries                — iso_code, dial_code, FK to currencies
005 customers                — UUID PK, KYC status, is_active
006 accounts                 — UUIDv7, NO balance column
007 transactions             — TRANSFER/REFUND/ADJUSTMENT, PENDING/COMPLETED/FAILED
008 journal_entries          — immutable, append-only, BIGINT amounts
009 idempotency_keys         — exactly-once semantics
010 audit_logs               — compliance trail
011 account_limits           — DAILY/MONTHLY/YEARLY/TRANSACTION limits per account
012 indexes                  — composite indexes for query optimization
013 seed_platform_accounts  — platform float/cash/revenue accounts + test data
014 fix_idempotency_keys    — corrected column types (append-only migration pattern)
```

---

## Performance

Load tested with k6 on a single development machine (Go app + PostgreSQL + Prometheus + Grafana running locally).
![k6 Load Test](docs/images/load_test.png)
| Scenario              | VUs | TPS | p50   | p95   | p99   | Errors |
| --------------------- | --- | --- | ----- | ----- | ----- | ------ |
| Baseline              | 1   | 95  | 10ms  | 12ms  | 14ms  | 0%     |
| Realistic concurrency | 20  | 332 | 53ms  | 107ms | 147ms | 0%     |
| Stress test           | 200 | 526 | 143ms | 459ms | 579ms | 0%     |

**94,660 transfers completed under stress test with zero errors or data corruption.**

Latency increases under high concurrency because `SELECT FOR UPDATE` serializes writes
to the same account pair by design. This is correct behavior for a financial system —
concurrent writes to the same account must be ordered. In production, load is distributed
across millions of account pairs, eliminating this bottleneck.

---

## Observability
![Grafana Dashboard](docs/images/grafana.jpg)
- **Prometheus** — scrapes `/metrics` every 15 seconds
- **Grafana** — 7-panel dashboard: request rate, p50/p95/p99 latency, transfer count, error rate, total requests
- **Structured logging** — zap with request IDs on every log line
- **Health endpoint** — `GET /health`

---

## Tech Stack

| Layer            | Technology             |
| ---------------- | ---------------------- |
| Language         | Go                     |
| Database         | PostgreSQL             |
| Router           | Chi v5                 |
| Connection Pool  | pgx/v5 + pgxpool       |
| Migrations       | golang-migrate         |
| Logging          | zap                    |
| Metrics          | Prometheus             |
| Dashboards       | Grafana                |
| Load Testing     | k6                     |
| Containerization | Docker, Docker Compose |

---

## API

### POST /api/v1/transfers

```json
{
  "from_account_id": "uuid",
  "to_account_id": "uuid",
  "amount": 1000,
  "currency": "INR",
  "idempotency_key": "unique-key-per-request"
}
```

**Response 200:**

```json
{
  "transaction_id": "uuid",
  "status": "COMPLETED",
  "created_at": "2026-04-02T00:00:00Z"
}
```

**Error responses:**
| Status | Reason |
|--------|--------|
| 400 | Invalid request body |
| 404 | Account not found |
| 403 | KYC not verified |
| 422 | Insufficient balance, account inactive, or limit exceeded |
| 500 | Internal server error |

---

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)

### Run Locally

```bash
# Clone the repo
git clone https://github.com/PacemakerX/ledger-core.git
cd ledger-core

# Copy environment variables
cp .env.example .env
cp .env.docker.example .env.docker

# Start PostgreSQL
docker-compose up postgres -d

# Run database migrations
make migrate-up

# Start the server
go run cmd/server/main.go
```

Server runs at `http://localhost:8080`

### Run Full Stack (with Prometheus + Grafana)

```bash
docker-compose up -d
```

| Service    | URL                   |
| ---------- | --------------------- |
| API        | http://localhost:8080 |
| Prometheus | http://localhost:9090 |
| Grafana    | http://localhost:3000 |

### Run Load Tests

```bash
k6 run load_test.js
```

---

## Architecture Decision Records

8 ADRs documented in `docs/adr/`:

```bash
001 — Why Go
002 — Why PostgreSQL (MVCC, ACID, WAL, SELECT FOR UPDATE)
003 — Why pgx over database/sql
004 — Why Chi over Gin/Echo
005 — Why Zap (structured logging, zero allocation)
006 — Why UUIDv7 (time-ordered, less index fragmentation)
007 — Why golang-migrate over GORM AutoMigrate
008 — Why Double-Entry Accounting
```

---

## License

MIT — see [LICENSE](LICENSE) for details.