# ADR-006: Why UUID for transaction ID  

## Status

Accepted

## Date

2026-03-24

## Context

We are debating about how to uniquely identify each row in a database, the debase is between UUID ( Universal Unique Identifiers) Vs Auto-Incrementing IDs
A UUID is a 128-bit unique identifier used for transaction IDs, customer accounts, API request and database recors. Fintech application rely on UUIDs to ensure each transaction is unique and traceable.


If two loan approvals occur simultaneously , a race condition could assign the same transaction ID to two different user, leading to financial mismatches. UUIDs eliminate this risk by ensuring every transaction has a globally unique identifier.

### Idempotency 
- Idempotency means performing the same operation multiple times has the same effect as doing it once.
	- Payment Processing - Prevent duplicate charges due to network  failures.
	- Loan Disbursements - Ensures a customer doesn't receive funds twice
### UUID Variants

- **UUIDv1**
    - Time + MAC address based
    - Leaks machine identity → not suitable
- **UUIDv4**
    - Fully random
    - Strong uniqueness
    - Causes index fragmentation due to randomness
- **UUIDv7**
    - Time-ordered (timestamp + randomness)
    - Maintains uniqueness while improving index locality
- **ULID**
    - Lexicographically sortable alternative

### **Index Fragmentation**:  
Random identifiers result in fragmented indexes, meaning data is stored non-sequentially.For databases liek MySQL and PostgreSQL, this fragmentation impacts performance by increasing page splits and slowing down access .


## Decision

We decide to go forward with UUIDv7


## Consequences

### Positive

- Globally Unique across systems.
- hard to guess.
- Supports distributed inserts without collisions.
- Improved write performance than UUIDv7 
- 
### Negative

- Large index sizez ( 16 bytes).
- Slight timestamp information leakage.
- Requires library support .

### Neutral

- Slightly more complex debugging 
- Human readability is reduced.

## Alternatives Considered

| Option         | Reason Rejected                                                     |
| -------------- | ------------------------------------------------------------------- |
| Auto-Increment | Not suitable for distributed/ multi-writer system. Information leak |
| UUIDv1         | Leaks MAC address and system meta data                              |
| UUIDv4         | Suffers with index fragmentation                                    |
| ULID           | Good alternative, UUIDv7 aligns better with industry standard.      |
