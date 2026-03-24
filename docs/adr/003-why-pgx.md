# ADR-003: Why PGX

## Status

Accepted

## Date

2026-03-24

## Context

We are selecting a database driver, a layer that enables our backend server to communicate with PostgreSQL.

In Go, database interaction is typically abstracted via the standard `database/sql` package, which requires a driver to implement the actual communication with the database.

- We evaluated the following options:
	- datbase/sql - its an extension to the standard `database/sql` package from the standard library.
	- github.com/llib/pq - pure Go Postgres driver for database/sql. For a long time, this driver was the standard by default. Currently , it has lost its relevance and is  not developed by the author.
		- pq - only works with `database/sql`
		- no standalone mode 
		- no advanced feature
		- `database/sql → pq → PostgreSQL`
	- github.com/jack/pgx  - PostgreSQL driver and toolkit for Go.
		- Works in two modes 
			- Compatibility mode `database/sql → pq → PostgreSQL `
			- Native mode ( powerful  `pgx → PostgreSQL`
		- When we use **pgx natively**, don’t just get one package we you get a **toolkit**.
			- pgx  -> low level  driver
			- pgxpool -> connection pool built on top of pgx 
			- `pgx = 1 connection (low-level)   pgxpool = many connections (production-ready)`

### Core parts:

- `pgx` → low-level driver (single connection)
- `pgxpool` → connection pool built on top of pgx

### Note 
Connection pool = reuse DB connections instead of recreating them per request

## Decision

We decide to use the pgx driver as our driver

## Consequences

### Positive

- px has human-readable erros, while lib/pq throws panics. if you don't catch a panic, the program will crash.
- With pgx we have an option to configure every connection independently explicit control over max-connection, idle connections and lifecycle.
- Access to advanced postgreSQL features like `BATCH Queries` 
- Uses PostgreSQL binary protocol instead of text Numbers sent as actual bytes, not string "12345" 
	- Faster serialization, less parsing overhead Lower CPU usage at high throughput
- pgxpool.Pool is safe for concurrent use
	- Multiple goroutines can safely share one pool
	-  Critical for handling concurrent transfers

### Negative

- Higher complexity vs `database/sql`
	- We have to deal with context everywhere
	- mangy tools expect `database/sql`
- No built-in struct mapping like `sqlx`
- Risk of misconfiguration
### Neutral

- We lose abstraction `database/sql`
- More Postgres-specific code
	- Using:
	    - COPY
	    - batch queries
	    - pgx types

## Alternatives Considered
| Option                | Reason Rejected                                                          |
| --------------------- | ------------------------------------------------------------------------ |
| database/sql + lib/pq | Generic interface, loses PostgreSQL specific features, lib/pq deprecated |
| GORM                  | Too much magic, hides SQL queries, can't debug what queries it generates |
| sqlx                  | Better than GORM but still wraps database/sql, lose pgx native benefits  |
| sqlc                  | Good option, generates type-safe code from SQL, revisit in v2            |
