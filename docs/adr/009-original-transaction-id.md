# ADR 009 — Why originalTransactionID on Refund Transactions

## Status

Accepted

## Date

2026-03-25

## Context

ledger-core supports partial refunds — a merchant can refund ₹600 of a ₹1000
transfer in one request, and ₹400 in a second request. This requires the system
to track how much has already been refunded against any given transfer.

The naive approach is to query journal entries and sum CREDIT entries on the
original sender's account. This is fragile — it conflates refunds with other
credits (top-ups, adjustments) and requires filtering by transaction type
across a join, which gets expensive at scale.

## Decision

Every REFUND transaction stores an `original_transaction_id` — a foreign key
pointing to the TRANSFER transaction it is reversing.

```sql
ALTER TABLE transactions
ADD COLUMN original_transaction_id UUID REFERENCES transactions(id);
```

To calculate total refunded against a transfer:

```sql
SELECT COALESCE(SUM(amount), 0)
FROM transactions
WHERE original_transaction_id = $1
AND type = 'REFUND'
AND status = 'COMPLETED'
```

The partial refund guard then becomes:
if req.Amount + totalAlreadyRefunded > original.Amount → reject

## Consequences

**Positive:**

- Single query to calculate total refunded — no joins, no journal entry scanning
- Refund transactions are self-documenting — the link to the original is explicit
- Supports arbitrary partial refunds up to the original amount
- Prevents over-refunding even across multiple partial refund requests
- Audit trail is complete — every refund points to its origin

**Negative:**

- Adds a nullable column to transactions — ADJUSTMENT and TRANSFER rows
  will always have NULL here, which is acceptable
- Circular FK risk is avoided because a REFUND can only point to a TRANSFER,
  never to another REFUND — enforced at the service layer

## Alternatives Considered

**Query journal entries directly:**
Sum CREDIT entries on the sender's account filtered by transaction type.
Rejected — fragile, expensive, conflates refunds with other credit operations.

**Store refunded amount on the original transaction:**
Add a `refunded_amount` column to transactions and update it on each refund.
Rejected — violates append-only principle. Updating the original transaction
record after the fact breaks auditability. The original transfer should be
immutable once COMPLETED.

**Allow only full refunds:**
No partial refund tracking needed — either fully refunded or not.
Rejected — real payment systems (Stripe, Razorpay) support partial refunds.
Limiting to full refunds would make ledger-core less realistic.

## References

- Stripe refunds API — partial refunds via `amount` parameter
- Razorpay refund documentation — supports multiple partial refunds per payment
- ADR 008 — Why Double-Entry Accounting
