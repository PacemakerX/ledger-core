# ADR-002: Why PostgreSQL

## Status

Accepted

## Date

2026-03-24

## Context

We need a database which can store the records of concurrent users and their transaction details

## Decision

We have decided to go with PostgreSQL

## Consequences

### Positive

- PostgreSQL is ACID compliant
  - Atomicity : Transaction is either fully succeeds or fails, no hanging -> `BEGIN / COMMIT / ROLLBACK ` and `WAL`
  - Consistency : Data always remain valid - > `CONSTRAINS:  PRIMARY KEYS/FOREGIN KEYS/UNIQUE / CHECK`
  - Isolation : Concurrent transaction doesn't corrupt state -> `MVCC + Locks  + Isolation Levels`
  - Durability : Committed data is never lost -> `WAls and fsync`
- PostgreSQL is is known for High Concurrency with MVCC ( Multi-Version Concurrency Control)
  - The Key insight is PostgreSQL keeps multiple version of the same row and shows each transactions the "right" version based on when it started.
  - Each transaction gets a snapshot when it begins, We only see data that was committed before your snapshot, Multiple version of the same row exit in the table simultaneously
- MVCC allows transactions to read consistent snapshots, so Transaction B can still see the old balance while Transaction A is in progress.
  - However, MVCC alone does not prevent lost updates.
  - To ensure correctness, mechanisms like `SELECT FOR UPDATE` or atomic updates are required, which lock the row and ensure that concurrent transactions are executed safely in sequence.
- PostgreSQL supports parallel queries and native partitioning
- Full-text Search , JSONB for unstructured data
  - Relational data - > for records
  - JSON -> meta data
- Indexing Flexibility
  - B-Tree -> default
  - GIN -> JSON/search
- Write-Ahead-Logs
  - Every change is logged before being committed.
  - We get crash recovery and point-in-time recovery.
- Visibility and Control
  - Use pg_stat_statement to find slow queries.
  - Tune Memory,Write-Ahead-Logs.
  - Leverage EXPLAIN ANALYZE for granular insight into execution plans.

### Negative

- In scenarios of write-heavy operations dominate the workload, PostgeSQL can shows slower performance compared to other database systems like MySQL or MariaDB. While PostgreSQL excels in read-heavy environments and complex query processing, the additional overhead from ACID compliance and its robust transaction handling can sow down write performance in certain cases.
- Disk + WAL = inherent latency
  - Every commit -> WAL flush ( fsync)
  - Ensures durability but adds delay
- Vaccum Overhead
  - MVCC creates multiple version of same row but this comes at an overhead of clearing the older versions ( vaccummed)

### Neutral

- Schema changes require migrations
  (intentional — enforces discipline)
- Single database simplifies architecture
  at current scale
- PostgreSQL is open source, no licensing cost

## Alternatives Considered

| Option      | Reason Rejected                                                                  |
| ----------- | -------------------------------------------------------------------------------- |
| MongoDB     | No ACID across documents, eventual consistency unacceptable for financial ledger |
| MySQL       | Weaker default isolation levels, less sophisticated MVCC                         |
| CockroachDB | Distributed PostgreSQL — overkill for current single-node scale                  |
| Redis       | In-memory, not durable enough for financial records                              |
