
# ADR-007: Why Golang-Migrate over GORM AutoMigrate

## Status

Accepted 

## Date

2026-03-25
## Context

A migration is a file or set of instruction that defines how our database changes over schema changes over time. 

Migration results in 
1.  Deterministic schema evolution
2. Safe Deployments 
3. Auditability 
	1. We get a history of schema changes
	2. We know who changed what
	3. and Why it changed ( through ARD's)
4. Team Collaboration
	1. When something is added to the databse
	2. We just pull the latest migration files and run the migrations

We have two tools for migrations 
	1.  Golang-migrate 
	2. GORM's auto migrate

## Decision

We decide to go with Golang-migrate.


### golang-migrate
golang-migrate is a database migratio tool for Go that lets you apply versioned, explicit, and controlled schema changes to our databas.
It is widely used in production systems because it treats migrations as first-class artifcacts

> **We write SQL migrations -> tool applies them in order -> tracks what's already applied** 
>   schema_migrations = checkpoint of DB schema


1. golang-migrate maintains a table like 

```SQL
schema migrations 
------------------
version   BIGINT
dirty     BOOLEAN 

version -> last applied migration
dirty -> indicates failed/partial migration
```

2. It's Deterministic and reproducible
3. Rollback support 
4. Explicit Control ( we write the exact SQL queries)

### GORM AutoMigrate
Gorm AutoMigrate is a feature of the GORM ORM that automatically updates our database schema based on our Go struct models

> **You define models → GORM figures out schema → applies changes automatically**

**What it does not do**
	 - It does not drop unused columns
	 - It does not rename columns
	 - It cannot handle complex schema changes
	 - Provide proper rollback support
	 - Track migration history 

Positives
- Fast development time 
	- We can instantly go from models to schema -> Best for MVP's 
	- Less boilerplate, no need to write explicit SQL
	- Schema stays insync with the code
## Consequences

### Positive

- We have explicit control over schema ( we write the exact SQL queries which needs to be executed ).
- We can track each migration and rollback in case of failure or partial migrations
- The dirty flag is a safety mechanism
  - If a migration fails halfway through dirty = true prevents any further migrations 
  - Forces explicit investigation before proceeding Protects against cascading schema corruption.
### Negative

- More boilerplate than GORM AutoMigrate
- Every schema change requires writing  both up and down migration files

- Developer discipline required
- Team must remember to create migration  files for every schema change Easy to forget in fast iteration

- Cannot auto-generate from structs
 -  Unlike GORM, we write raw SQL manually Acceptable tradeoff for explicit control

### Neutral

- Migration files are plain SQL
  - Readable by anyone regardless of Go knowledge DBAs can review without understanding Go code

- golang-migrate is database agnostic 
	- Works with PostgreSQL, MySQL, SQLite Not locked into one database

## Alternatives Considered

| Option           | Reason Rejected                                                                                                                                                                 |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GORM AutoMigrate | No explicit control over schema  + No versioning + No rollback support  + Can silently alter schema + Dosen't drop unused columns  + Dosen't rename columns  + Weak auditabilty |
