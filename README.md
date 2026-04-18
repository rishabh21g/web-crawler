# Go Web Crawler

A simple concurrent web crawler in Go that:
- crawls **only one domain**,
- stops at a **max depth**,
- avoids **duplicate URLs**,
- fetches only **HTML (`text/html`)**,
- enforces **timeouts** and **response size limits**,
- uses a **worker pool** (goroutines) + **channel** + **WaitGroup**.

This repo is intentionally “learning-first”: the code is small, readable, and focused on core concurrency patterns.

## What this repo demonstrates

You’ll learn practical Go fundamentals by reading and extending this code:

- **Goroutines**: concurrent workers running in parallel
- **Channels**: passing `URLTask` work items between producers and consumers
- **WaitGroup**: tracking “in-flight” tasks until the crawl finishes
- **Mutex**: protecting shared state (`Seen` map) against races
- **context.Context**: request timeouts/cancellation for HTTP calls
- **net/http**: building requests, setting headers, handling status codes
- **HTML parsing**: extracting `<a href="...">` links with `x/net/html`
- **Defensive crawling**: validate content-type, limit body size, reject junk links

## Project layout

- Entry point: [cmd/crawler/main.go](cmd/crawler/main.go)
- Task model + config: [internal/models/types.go](internal/models/types.go)
- Scheduler (dedupe/domain/depth rules): [internal/scheduler/scheduler.go](internal/scheduler/scheduler.go)
- Fetch HTML with validation/limits: [internal/fetcher/fetcher.go](internal/fetcher/fetcher.go)
- Parse HTML and extract links: [internal/parser/parser.go](internal/parser/parser.go)
- Worker logic (fetch → parse → enqueue): [internal/worker/worker.go](internal/worker/worker.go)
- Storage placeholder: [internal/storage/storage.go](internal/storage/storage.go)

## How it works (end-to-end)

At a high level, the crawler is a loop:

1. Take a `URLTask` from the tasks channel
2. Validate & dedupe it via the Scheduler
3. Fetch HTML (`GET`, timeout, size limit, content-type check)
4. Parse HTML to discover links
5. Convert discovered links into new `URLTask`s (depth+1)
6. Send new tasks back into the channel (if allowed)
7. Mark the current task “done”

### Producer–consumer pattern (what to notice)

This repo is a classic **producer–consumer** design:

- **Work queue**: the `tasks chan models.URLTask`
- **Consumers**: worker goroutines (`worker.Worker`) range over the channel and consume tasks
- **Producers**: workers also *produce new tasks* (links discovered via parsing) and push them back into the same channel
- **Coordinator**: `main` seeds the first task and waits for completion using a `sync.WaitGroup`

Key point: the WaitGroup tracks “in-flight tasks”, not the number of worker goroutines.

### What stops the crawl?

In [cmd/crawler/main.go](cmd/crawler/main.go):

- `wg.Add(1)` is called when a task is accepted for processing (before enqueue)
- Each worker calls `wg.Done()` exactly once per task it processes (or drops)
- When the task graph is exhausted, `wg.Wait()` unblocks
- The channel is then closed (`close(tasks)`) so workers exit their `for task := range tasks` loops

## Running

From repo root:

```bash
go run ./cmd/crawler
```
