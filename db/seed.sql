-- ============================================================
-- Seed data for Level Up Backend
-- Run with: make db/seed
-- ============================================================

-- ── Modules ──────────────────────────────────────────────────
INSERT INTO modules (title, slug, description, order_index, estimated_hours) VALUES
(
    'Go Concurrency & Graceful Shutdown',
    'go-concurrency',
    'Master production-grade concurrency patterns: worker pools, context propagation, backpressure, goroutine leak prevention, and graceful shutdown. Build systems that handle load and fail cleanly.',
    1,
    8.0
),
(
    'Distributed Systems Fundamentals',
    'distributed-systems',
    'Idempotency, retries with exponential backoff, circuit breakers, rate limiting, message queues, and the difference between exactly-once vs at-least-once delivery.',
    2,
    10.0
),
(
    'Reliability & Observability',
    'reliability',
    'Structured logging, distributed tracing, SLOs, error budgets, deploy strategies (blue/green, canary), and building systems that are debuggable in production.',
    3,
    8.0
),
(
    'Architecture & Systems Thinking',
    'architecture',
    'Tradeoff simulations, CAP theorem in practice, data modeling under load, scaling decisions, and how senior engineers think about system design.',
    4,
    12.0
)
ON CONFLICT (slug) DO NOTHING;

-- ── Lessons: Module 1 — Go Concurrency ───────────────────────
INSERT INTO lessons (module_id, title, slug, content, order_index, estimated_minutes)
SELECT
    m.id,
    l.title,
    l.slug,
    l.content,
    l.order_index,
    l.estimated_minutes
FROM modules m
CROSS JOIN (VALUES
    (
        'Worker Pools',
        'worker-pools',
        '# Worker Pools

A worker pool bounds the number of goroutines doing concurrent work. Without it, an unbounded number of goroutines will exhaust memory under load.

## The Pattern

```go
func NewWorkerPool(ctx context.Context, workers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result, workers)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case results <- process(ctx, job):
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

## Key Observations
- Workers stop when `jobs` channel closes (range exits)
- Context cancellation exits workers early
- `wg.Wait()` in a goroutine closes results only after all workers finish
- Buffer size on results prevents workers blocking on slow consumers

## Senior Insight
The right number of workers depends on whether work is CPU-bound (use `runtime.GOMAXPROCS`) or IO-bound (higher concurrency, benchmark to find the ceiling).',
        1,
        45
    ),
    (
        'Context Propagation',
        'context-propagation',
        '# Context Propagation

Context is how you cancel work across goroutine boundaries. Threading it correctly is one of the most common things mid-level engineers get wrong.

## Rules

1. Always accept `ctx context.Context` as the **first parameter**
2. Never store context in a struct — pass it per-call
3. Always check `ctx.Err()` or `<-ctx.Done()` in long-running loops
4. Wrap with `context.WithTimeout` at the boundary of external calls

## The Leak Pattern (Wrong)

```go
// BAD: ignores context, goroutine leaks if caller cancels
go func() {
    result, err := db.QueryRow("SELECT ...")  // blocks forever
    ch <- result
}()
```

## Correct

```go
// GOOD: context threaded through
go func() {
    result, err := db.QueryRowContext(ctx, "SELECT ...")
    select {
    case ch <- result:
    case <-ctx.Done():
    }
}()
```

## Timeouts at System Boundaries

```go
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
resp, err := httpClient.Do(req.WithContext(ctx))
```

`defer cancel()` is non-negotiable — it prevents context leak even on the happy path.',
        2,
        40
    ),
    (
        'Backpressure & Rate Limiting',
        'backpressure',
        '# Backpressure

Backpressure is the signal from a downstream system that it is overwhelmed. Without it, fast producers destroy slow consumers.

## Channel-Based Backpressure

```go
// Bounded channel = natural backpressure
jobs := make(chan Job, 100)

// Producer: non-blocking send with backpressure signal
select {
case jobs <- job:
    // accepted
default:
    // queue full — return 429 to caller
    return ErrBackpressure
}
```

## Token Bucket Rate Limiter

```go
limiter := rate.NewLimiter(rate.Every(time.Second/100), 10) // 100 rps, burst 10

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    // handle request
}
```

## Senior Insight
Backpressure is a feature, not a failure. A system that returns 429 gracefully is more reliable than one that accepts everything and crashes. Design your ingestion layer to signal overload explicitly.',
        3,
        35
    ),
    (
        'Goroutine Leak Detection',
        'goroutine-leaks',
        '# Goroutine Leak Detection

A goroutine leak is a goroutine that is no longer doing useful work but is never garbage collected. In production, leaks accumulate until the process runs out of memory.

## Common Causes

1. **Blocking channel send with no receiver**
```go
// BAD: if nobody reads results, this goroutine leaks
go func() {
    result := compute()
    results <- result  // blocks forever if caller returned early
}()
```

2. **Waiting on a channel that is never closed**
```go
done := make(chan struct{})
go func() {
    <-done  // leaks if done is never closed
}()
```

3. **Ticker without stop**
```go
// BAD
go func() {
    t := time.NewTicker(time.Second)
    for range t.C { work() }  // leaks: t.Stop() never called
}()
```

## Detection with goleak

```go
func TestWorkerPool(t *testing.T) {
    defer goleak.VerifyNone(t)  // fails test if goroutines leak
    // ... test code
}
```

## The Fix: Always Pair Goroutines with Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // guarantees goroutines exit

go func() {
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-jobs:
            process(job)
        }
    }
}()
```',
        4,
        40
    ),
    (
        'Graceful Shutdown',
        'graceful-shutdown',
        '# Graceful Shutdown

Graceful shutdown means: stop accepting new work, finish in-flight work, clean up resources — in that order.

## The Pattern

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: router}

    // Channel to receive OS signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // Start server in background
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Block until signal
    <-quit
    log.Println("shutting down...")

    // Give in-flight requests 15 seconds to complete
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // 1. Stop accepting new requests + drain in-flight
    srv.Shutdown(ctx)

    // 2. Drain background workers (close job channel, wait for wg)
    workerPool.Shutdown()

    // 3. Close DB pool (returns connections, flushes prepared statements)
    dbPool.Close()

    log.Println("stopped")
}
```

## The Order Matters
1. **HTTP server shutdown first** — stops new requests from creating new work
2. **Workers second** — drains queued jobs
3. **DB last** — workers may need DB during drain

## Senior Insight
Your readiness probe (`/healthz`) should return 503 as soon as shutdown begins. This tells the load balancer to stop routing traffic before `Shutdown()` is called, giving you a clean drain.',
        5,
        50
    ),
    (
        'DB Connection Pool Tuning',
        'db-pool-tuning',
        '# DB Connection Pool Tuning

The connection pool is the bridge between your app and Postgres. Misconfiguration here causes cascading failures under load.

## Key Parameters (pgxpool)

```go
cfg.MaxConns = 25              // max open connections
cfg.MinConns = 5               // keep warm connections alive
cfg.MaxConnLifetime = 1*time.Hour    // recycle connections (avoids stale)
cfg.MaxConnIdleTime = 30*time.Minute // return idle connections
cfg.HealthCheckPeriod = 1*time.Minute
```

## How to Size MaxConns

**Formula:** `MaxConns = (num_cores * 2) + effective_spindle_count`

For a 4-core app server hitting Postgres: `MaxConns ≈ 10–25`

**Do not set it to 100+.** Postgres handles connections as OS processes. At 200 connections Postgres spends more time on connection overhead than actual queries.

## The Thundering Herd Problem

If all `MaxConns` connections are saturated, new requests queue waiting for a free connection. This queue can grow unboundedly.

Fix: set an acquisition timeout:

```go
cfg.ConnConfig.ConnectTimeout = 3 * time.Second
```

Return 503 to the caller instead of waiting indefinitely. A fast failure is better than a slow queue.

## Senior Insight
Watch `pg_stat_activity` in production. If you see hundreds of `idle in transaction` connections, you have a missing `defer tx.Rollback()` somewhere — a transaction leak.',
        6,
        45
    )
) AS l(title, slug, content, order_index, estimated_minutes)
WHERE m.slug = 'go-concurrency'
ON CONFLICT (module_id, slug) DO NOTHING;

-- ── Skills: Module 1 ─────────────────────────────────────────
INSERT INTO skills (module_id, skill_name, order_index)
SELECT m.id, s.skill_name, s.order_index
FROM modules m
CROSS JOIN (VALUES
    ('Implement a bounded worker pool with context cancellation', 1),
    ('Thread context correctly through all function calls', 2),
    ('Implement channel-based backpressure with 429 response', 3),
    ('Identify and fix goroutine leaks in a code review', 4),
    ('Implement graceful shutdown with correct drain order', 5),
    ('Tune DB connection pool for a given workload', 6),
    ('Write a leak-free ticker-based background worker', 7)
) AS s(skill_name, order_index)
WHERE m.slug = 'go-concurrency'
ON CONFLICT DO NOTHING;

-- ── Assignment: Module 1 ──────────────────────────────────────
INSERT INTO assignments (module_id, title, description, rubric, estimated_hours)
SELECT
    m.id,
    'Build a Production-Grade Job Runner',
    E'Build a job runner in Go that processes tasks from a queue with the following constraints:\n\n**Requirements:**\n- Bounded worker pool (configurable size, default 5)\n- Each job has a timeout (context with deadline)\n- Backpressure: return error when queue is full (do not block)\n- Graceful shutdown: drain in-flight jobs on SIGTERM, max 30s\n- No goroutine leaks under any condition\n- Structured JSON logging for job start, complete, fail, timeout\n- DB pool for job persistence (mark jobs as processing/done/failed)\n\n**Scenarios your implementation must handle:**\n1. 1000 jobs submitted concurrently with 5 workers\n2. SIGTERM mid-processing — in-flight jobs complete, queued jobs are re-queued\n3. A job that hangs — context deadline fires, job marked as failed, worker continues\n4. DB unavailable — jobs fail with error, system continues running\n\n**Deliverable:** GitHub repo with README explaining your design decisions.',
    E'## Evaluation Rubric\n\n### Worker Pool (25 points)\n- [ ] Pool size is configurable and respected under load\n- [ ] Workers exit cleanly when context is cancelled\n- [ ] No goroutine leak verified with goleak in tests\n\n### Context & Timeouts (20 points)\n- [ ] Per-job timeout implemented with context.WithDeadline\n- [ ] Context threaded to all DB calls and external calls\n- [ ] Timeout fires correctly on a hanging job\n\n### Backpressure (15 points)\n- [ ] Non-blocking submit with clear error when queue full\n- [ ] Caller receives actionable error (not a hang)\n\n### Graceful Shutdown (25 points)\n- [ ] SIGTERM triggers shutdown\n- [ ] In-flight jobs complete (up to 30s)\n- [ ] Shutdown sequence: HTTP → workers → DB\n- [ ] Process exits cleanly (exit code 0)\n\n### Code Quality (15 points)\n- [ ] Structured logging on all job lifecycle events\n- [ ] Tests cover the goroutine leak, timeout, and shutdown scenarios\n- [ ] README explains design decisions and tradeoffs',
    8.0
FROM modules m
WHERE m.slug = 'go-concurrency'
ON CONFLICT (module_id) DO NOTHING;
