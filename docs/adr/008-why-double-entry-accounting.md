# ADR-008: Why Double Entry Accounting

## Status

Accepted 

## Date

2026-03-25

## Context
Double entry is the bookkeeping principle which states that every financial transaction has equal and opposite effects as both an asset and a liability, and therefore it must be recorded as such in two separate place. `Assets = Liabilities + Equity`

IN accounting, a credit as an entry that increases a liability account or decreases an asset account. A debit is the opposi8te. it is an entry that increases an asset account or decreases a liability account.

IN the double-entry accounting system, transactions are recorded terms of debits and credit. Since a debit in one account offsets a credit in another, the sum of all debits must equal the sum of all credits.

Five types of accounts that all business transaction s can be classified:
1. **Assets** (what you own)
2. **Liabilities** (what you owe)
3. **Equity** (owner’s value)
4. **Revenue** (money coming in)
5. **Expenses** (money going out)
## Decision

We decide to go with double-entry accounting ledger.

Single-Entry: Records what happened
	`balance = balance - 500`		``

Double-Entry: Records what happened + verifies it must be true

```
Debit  : Expense 500
Credit : Bank 500
```

Now we have 

```
SUN ( DEBITS ) = SUM ( CREDITS )
```

```
Transfer ₹500 from Account A to Account B:

DEBIT  Account A  ₹500  (asset decreases)
CREDIT Account B  ₹500  (asset increases)

Sum of debits = Sum of credits = ₹500 ✓
Ledger is balanced.
```
## Consequences

### Positive

- Partial Write / system crash
	- Scenario Debit recorded , Credit failed
	- Result Debit != Credit  ` System knows The transaction is broken` 
- Duplicate transaction ( idempotency failure ) 
	- Scenario: Same request processed twice
	- Double entry detection Duplicate transaction IDs 
	- Ledger imbalance across expected flows
- Prevents Race Condition bugs
	- Only append entries 
	- No balance mutation

### Negative

- More complex to implement than single entry
	 - Every transfer requires minimum 2 journal entries

- More storage than single entry
	-   Every transaction creates multiple rows
  
- Harder to query current balance
  Must SUM journal entries vs reading one column
  SELECT SUM(amount) WHERE entry_type = 'CREDIT'
  minus SUM(amount) WHERE entry_type = 'DEBIT'
  
- Steeper learning curve for engineers unfamiliar with accounting concepts

### Neutral

- Industry standard for financial systems
  Stripe, Razorpay, every bank uses this
  
- Immutable audit trail by design
  Append-only means complete history always exists
  
- Regulatory compliance friendly
  Every rupee movement is traceable

## Alternatives Considered

| Option              | Reason Rejected                                                                                                            |
| ------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| Single Entry Ledger | Cannot detect imbalances + Cannot trace money flow + Cannot identify duplicate transaction + Causes  partial write problem |
