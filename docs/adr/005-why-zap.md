# ADR-005: Why Zap

## Status

Accepted

## Date

2026-03-24
## Context

Log files are one of the most valuable assets that a developer has. Usually when something goes wrong in production the first thing which is checked are logs.
But logs need to be simple, filterable.

For log files to be machine-readable more advanced functionality,they need to be written in a structured format that can be easily parsed. This could be XML, JSON or formats.

We are discussing how to make logs more structured?
Zap is known for its high performance and structured logging capabilities, making it suitable for applications that require speed and efficiency. 
Logrus offers flexibility and ease of use with structured logging , but it may not match Zap's performance in high-load secnarios.

## Decision

We will use Zap logging which is used by Uber.

### Understanding Zero-Allocation Logging
- Zero-allocation logging is a crucial performance optimisation technique that minimises memory allocation during logging operations. This approach h significantly reduces garbage collection overhead and improves application performance, especially in high-throughput systems.
- The fundamental concept behind zero-allocation logging revolves around prevention heap allocation  during log operation s. Traditional logging often creates temporary strings, interfaces, and other objects that requires memory allocation. Zero-allocation logging eliminates these allocations through careful API design and memory reuse
- Memory allocation means the allocation of  memory required for a program to operate, and there are those for the stack area ( called static allocatino ) and those for the heap  ( called dynamic allocation )
	- When allocation memory on stack , the allocation size and the timing of allocation/release are statically determined when the program is written.
	- ON the other hand, the allocatio on heap area can be according to the situation when the program is executed without speicifying the maximum memmemory allocation size at the time of declaration, and the memory allocation size at the the time of declaration. 
## Consequences

### Positive

- Near Zero -Allocation logging 
	- Most go loggers uses `fmt.Sprintf` which creates heap allocations -> GC pressure  ( STW pauses) -> Latency spikes
-  Zap is consistently among the  fastest Go loggesrs. , Zap is 10x faster than Logrus in hotpaths
- Zap is built around structured llogs\
	- Logs become queryable data
- Strong Separation: Dev vs Production 
	- Zap gives you : 
		- zap.NewDevelopment() - > readalbe logs
		- zap.NewProduction()  - > JSON, optimized  
# ## Negative

-  Verbose, Less Ergonomic API 
	- Zap forces structured fields which may become slow 
- Learning Curve
- Overkill for small systems
### Neutral

- Strongly Opinionated Design
- Json first output
- Zap is not part of standard library.

## Alternatives Considered

| Option  | Reason Rejected                    |
| ------- | ---------------------------------- |
| Zerolog | Less flexible than Zap             |
| Logrus  | High allocations                   |
| log     | Very Basic + No structured Logging |
