# ledger-core

A production-grade double-entry accounting ledger system built in Go with PostgreSQL. Handles concurrent transfers, idempotency, and refunds with full observability.

---

## Overview

ledger-core is a financial ledger backend implementing double-entry bookkeeping principles. Every transaction creates balanced debit/credit journal entries, ensuring the ledger is always consistent and auditable.

Built to demonstrate production-grade fintech backend engineering:
- **Correctness** — ACID transactions, SELECT FOR UPDATE concurrency control
- **Reliability** — Idempotency keys ensuring exactly-once semantics
- **Observability** — Prometheus metrics, Grafana dashboards, structured logging
- **Performance** — Target 1500 TPS, p99 latency under 50ms


### Core Design Decisions

- **Immutable Ledger** — Append-only journal entries. No UPDATE, only INSERT.
- **Double-Entry Accounting** — Every transaction creates balanced DEBIT + CREDIT entries.
- **Idempotency** — All write operations accept an idempotency key to prevent duplicates.
- **Concurrency** — SELECT FOR UPDATE with consistent lock ordering to prevent deadlocks.
- **Refunds** — Time-windowed (90 days), implemented as reverse journal entries.

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go |
| Database | PostgreSQL |
| Containerization | Docker, Docker Compose |
| Metrics | Prometheus |
| Dashboards | Grafana |
| Load Testing | k6 |
| Migrations | golang-migrate |

---

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose

### Run Locally

```bash
# Clone the repo
git clone https://github.com/PacemakerX/ledger-core.git
cd ledger-core

# Copy environment variables
cp .env.example .env

# Start PostgreSQL, Prometheus, Grafana
docker-compose up -d

# Run database migrations
make migrate-up

# Start the server
make run
```

Server runs at `http://localhost:8080`

## Performance

| Metric | Target | Achieved |
|--------|--------|----------|
| Throughput | 1500 TPS | TBD |
| p99 Latency | < 50ms | TBD |
| Uptime | 99.9% | TBD |

> Results will be updated after k6 load testing is complete.

## License

MIT — see [LICENSE](LICENSE) for details.