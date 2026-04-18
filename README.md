# Go Web Crawler

A concurrent, multi-worker web crawler written in Go, designed to efficiently traverse websites by fetching and parsing HTML pages. It features configurable depth limits, domain restrictions, and a robust worker pool pattern to handle network operations concurrently.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Features](#features)
- [Current Limitations](#current-limitations)
- [Future Goals & Enhancements](#future-goals--enhancements)

## Overview

This project implements a web crawler that uses a pool of worker goroutines to process URLs concurrently. The main components include:
- A **Scheduler** to prevent duplicate work and enforce crawling rules.
- A **Fetcher** to handle HTTP requests with timeouts and content validation.
- A **Parser** to extract links from HTML documents.
- A **Worker Pool** to manage concurrent execution.

The crawler is designed to be simple, efficient, and easily extensible.

## Architecture

The crawler operates on a classic worker pool model. A central `main` function orchestrates the work by distributing `URLTask` items through a shared channel to a pool of worker goroutines.

```
┌──────────────────────────┐
│           Main           │
│  - Creates Scheduler     │
│  - Spawns Worker Pool    │
│  - Seeds Initial URL     │
└────────────┬─────────────┘
             │
     ┌───────▼─────────┐
     │   Tasks Channel │
     │ (models.URLTask)│
     └───────┬─────────┘
             │
  ┌──────────┼──────────┬──────────┐
  │          │          │          │
  ▼          ▼          ▼          ▼
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐
│Worker1│ │Worker2│ │...    │ │WorkerN│
└─┬─────┘ └─┬─────┘ └─┬─────┘ └─┬─────┘
  │         │         │         │
  └─────────┼─────────┼─────────┘
            │
    ┌───────▼────────────┐
    │  Processing Loop:  │
    │ 1. Check & Mark URL│
    │ 2. Fetch HTML      │
    │ 3. Parse for Links │
    │ 4. Add New Tasks   │
    └────────────────────┘
```

## Project Structure

```
.
├── cmd/server/main.go      # Main application entry point and worker logic
├── internal/
│   ├── crawler/
│   │   ├── crawler.go      # (Placeholder)
│   │   └── scheduler.go    # URL validation and duplicate prevention
│   ├── fetcher/
│   │   └── fetcher.go      # HTTP client for fetching HTML
│   ├── frontier/
│   │   └── queue.go        # A thread-safe queue for URLs
│   ├── models/
│   │   └── types.go        # Core data structures (URLTask, Config)
│   └── parser/
│       └── parser.go       # HTML parsing and link extraction
├── go.mod
└── README.md
```

## How It Works

1.  **Initialization**: The `main` function initializes the `scheduler`, the `tasks` channel, and a `sync.WaitGroup`. It then launches a fixed number of `Worker` goroutines.
2.  **Seeding**: An initial `URLTask` is sent to the `tasks` channel to begin the crawl.
3.  **Worker Execution**: Each worker pulls a `task` from the channel and performs the following steps:
    a. **Check**: It asks the `scheduler` if the URL should be processed. The scheduler checks if the URL has been seen before, is within the allowed domain, and has not exceeded the max depth.
    b. **Fetch**: If the check passes, the `fetcher` downloads the HTML content, respecting timeouts and validating the content type.
    c. **Parse**: The returned HTML is passed to the `parser`, which extracts all unique links.
    d. **Queue**: The newly found links are converted into `URLTask`s and sent back to the `tasks` channel for other workers to process.
4.  **Synchronization**: The `sync.WaitGroup` is used to track the number of active tasks. The `main` function waits for the `WaitGroup` counter to return to zero, indicating that all work is complete.
5.  **Termination**: Once the `WaitGroup` is unblocked, the `main` function closes the `tasks` channel, which signals the worker goroutines to terminate their loops and exit.

## Getting Started

### Prerequisites

- Go (version 1.18 or newer recommended)

### Running the Crawler

1.  Clone the repository.
2.  Navigate to the project directory.
3.  Run the main application from your terminal:
    ```bash
    go run cmd/server/main.go
    ```

## Configuration

Configuration is currently hardcoded in `cmd/server/main.go`.

- **`WORKERS`**: The number of concurrent worker goroutines (default: `10`).
- **`models.Config`**:
    - `Domain`: The only domain the crawler is allowed to visit.
    - `MaxDepth`: The maximum number of links to follow from the start page.
- **`startURL`**: The initial URL to begin crawling.

## Features

- **Concurrent Crawling**: Utilizes a worker pool to fetch multiple pages simultaneously.
- **Domain Scoping**: Restricts the crawler to a single domain to prevent it from wandering.
- **Depth Limiting**: Stops crawling after reaching a configured link depth.
- **Duplicate Prevention**: A `Seen` map in the scheduler ensures each URL is processed only once.
- **Content Validation**: The fetcher ensures that the downloaded content is `text/html` and respects a maximum file size.
- **Timeout Handling**: Network requests will time out gracefully to prevent workers from getting stuck.

## Current Limitations

1.  **No `robots.txt` Support**: The crawler does not respect the `robots.txt` file of websites, which is poor etiquette.
2.  **In-Memory State**: The list of "seen" URLs is stored in memory, which does not scale to very large crawls.
3.  **No Rate Limiting**: The crawler may make requests too quickly, which can overwhelm a server or get the crawler's IP address blocked.
4.  **No Persistence**: The results of the crawl (the links found) are printed to the console but not saved.

## Future Goals & Enhancements

- **`robots.txt` Compliance**: Implement a parser for `robots.txt` and respect its rules before fetching any URL.
- **Persistent Storage**:
    - Store the "seen" URLs in a database (like Redis or BadgerDB) to support larger crawls and resuming a stopped crawl.
    - Save the crawl results (page titles, links, content) to a database or file.
- **Rate Limiting**: Add a per-domain rate limiter (e.g., a token bucket) to be a more responsible web citizen.
- **Configuration from File**: Load crawler settings from a configuration file (e.g., YAML or JSON) instead of hardcoding them.
- **Metrics and Monitoring**: Add a simple endpoint or logger to report on the crawler's progress, such as pages per second, number of active workers, and queue size.
- **Advanced Link Parsing**: Improve the parser to handle links found in JavaScript files or sitemaps.
- **Distributed Crawling**: Refactor the architecture to allow multiple crawler instances to coordinate work, potentially using a distributed queue like RabbitMQ or Kafka.
