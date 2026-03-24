# ADR-004: Why Chi Router
## Status

Accepted

## Date

2026-03-24 

## Context

Router is a component that decides  `Given an incoming request ( METHOD + URL ),which handler function should run` 
Go's default router is `net/http` 
So what's the problem with `net/http` ? 
	- Minimal by design
		- Cannot handle dynamic request
		- No route parameter
		- No middleware chaining  
		- No grouping / versioning 

We need a real world router which can handle dynamic request, middleware chaining and support route grouping.
We have a few option like Chi, Gin ,Fiber

## Decision

We decide to go with `Chi` router. 
- Why? 
	-  Minimal and design philosophy : Stay close to net/http, don't reinvent it
	- Idiomatic Go - refers to tht writing code in a way that alings with the Go programming's conventions, best practices, and community standards

## Consequences

### Positive

- Improved Maintainability 
- Better scalability of API Desgin

### Negative

- Slight learning curve for middleware patterns
- Less feature-rich than full frameworks

### Neutral

- Performance is comparable to other routers 

## Alternatives Considered

| Option     | Reason Rejected                                                                               |
| ---------- | --------------------------------------------------------------------------------------------- |
| `net/http` | Too Minimal + Lacks dynamic routing + No middleware chaining  +  No api grouping              |
| `Gin`      | Higher abstraction, uses custom context, less idiomatic, introduces tighter coupling          |
| `Fiber`    | Built on `fasthttp`, deviates from `net/http`, increases lock-in and reduces interoperability |
