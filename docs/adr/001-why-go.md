# ADR-001: Use Go for Core Ledger Services

## Status
Accepted

## Date
2026-03-24

## Context
I am building ledger-core as a solo project targeting
fintech backend roles at companies like Razorpay and
Zerodha. The ledger needs to handle concurrent transfers
safely without race conditions.

The system will process financial transactions where correctness and consistency are critical, while also needing to handle high concurrency.

## Decision
Use Go as the primary language for implementing core backend and ledger services.

## **Consequences**

### Positive
-  In-built garbage collector improves developers’ efficiency by automatically freeing up no longer needed memory allocations and keeping those still in use.
- Golang advantages is goroutines - threads that allow multiple process to run concurrently. Compare to multithreading used in other technologies,goroutines are extremely lightweight ( only a few kBs), meaning an app could use hundreds or thousands of threads at the same time without overloading the hardware.
- Goroutines can effectively communicate between  with each other with the help of go Channels,which results in low-latency communication. 
- Go has a strong support for concurrency and non-blocking I/O ( goroutines do not wait for I/O operation to finish ), Make resource-sharing efficient and is perfect for building large,distributed system. With this in mind, it is no surprise that over 75% of Cloud Native Computing Foundation projects [are written in Go](https://thenewstack.io/go-language-fuels-cloud-native-development/).
- [Why fintech companies use GoLang](https://surf.dev/go-for-fintech-projects-use-cases-by-global-companies/)
- Golang is a compiled language, offering performance close to C/C++.
- Golang's standard library includes packages for encryption, hashing, and secure communication.
- Strong typing that reduces runtime errors
- Go services are lightweight and efficient, they require fewer compute resources to handle the same workload compared to heavier runtimes. [Future Proofing](https://metadesignsolutions.com/future-proofing-your-fintech-app-why-go-is-the-choice-for-2026s-secure-high-speed-transactions/)
### Negative
-  Garbage collection may introduce latency in edge cases
	- Go uses a concurrent GC, i.e. It runs alongside the application ( not full stop-the-world like old JVM )
	- But it still consumes CPU cycles 
	- When allocation rates is high GC runs more frequently, GC competes with your application for CPU
- Everything that escaptes to the heap is GC-managed.
	- If our system creates a lot of objects ( per request, per tick, per message)
	- Heap grows rapidly 
	- GC triggers more often
- Stop-the-World still exist
	- Go's GC is mostly concurrent, but it still has Stop-the-world pauses for scanning which takes upto 2ms which is a lot if we want to achieve micro-second latency 
- Unlike C++/Rust:
	- We can't control allocation granularity
	- We can't place objects in specific memory regions
- C/C++ has no GC -> so no STW pauses
	- This results in manual memory management which is cumbersome
	- This make harder to scale
- Rust
	- No GC → no STW
	- Memory managed at compile time
	- Tradeoff: complexity
- Every other language JAVA( JVM )  and Python has one or the other form of GC 

### Neutral
- Opinionated Language Design , no inheritance , no  generics, no operator overloading. 
-  Go compiles into a single binary -> No runtime dependency.
- Go encourages using stdlib over third-party libraries -> resulting in stability  but feels very basic compare to Python
## Alternatives Considered

| Option               | Reason Rejected                                                                                                                                                                                                                                                                                                                                                                                                                             |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Node.js (TypeScript) | Simpler development but not ideal for CPU-bound and high-concurrency workloads  since, One main thread handles all execution, If one task is CPU-heavy -> everything else waits. Node shines at I/O, since it registers a callback and moves on immediately. **JS execution is tied to one thread**                                                                                                                                         |
| Python               | Uses GIL ( Global Interpreter Lock ) which ensures only one thread executes python bytecode at a time per process.                                                                 - **Concurrency** = multiple tasks _in progress_ (interleaved)<br>- **Parallelism** = multiple tasks _executing at the same time on multiple cores_<br>Python has **concurrency**  <br>Python (CPython) lacks **true parallelism for CPU-bound threads** |
| C++                  | High Complexity + Manual Memory Management + Memory Leaks. Does not aligns well with the ne need for rapid iteration and maintainability                                                                                                                                                                                                                                                                                                    |
