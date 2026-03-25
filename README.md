# service

A lightweight Go library for managing service lifecycles with built-in signal handling, graceful shutdown, and health check probes.

## Features

- **Simple API**: Register setup and service functions with minimal boilerplate
- **Graceful Shutdown**: Automatic context cancellation on OS signals or errors
- **Health Probes**: Built-in HTTP liveness and readiness handlers for orchestration systems
- **Concurrent Execution**: Run multiple service components concurrently with `errgroup`
- **Two-Phase Initialization**: Setup phase runs before service components start

## Installation

```bash
go get github.com/Attest/service
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/Attest/service"
)

func main() {
    // Create a service group that handles SIGINT
    group := service.NewSignals(os.Interrupt)

    // Register setup functions (run sequentially before processes start)
    group.Setup(func(ctx context.Context) error {
        log.Println("Initializing database connection...")
        // Setup code here
        return nil
    })

    // Register service processes (run concurrently)
    group.Register(func(ctx context.Context) error {
        srv := &http.Server{Addr: ":8080"}

        // Graceful shutdown on context cancellation
        go func() {
            <-ctx.Done()
            srv.Close()
        }()

        return srv.ListenAndServe()
    })

    // Start the service group
    log.Fatal(group.Start())
}
```

## Health Check Probes

The service group exposes HTTP handlers for liveness and readiness probes:

```go
http.HandleFunc("/live", group.Liveness())
http.HandleFunc("/ready", group.Readiness())
```

### Probe Behavior

**Liveness Probe** (`/live`):
- Returns `503 Service Unavailable` until the group has started
- Returns `200 OK` once `Start()` is called
- Indicates the service process is alive

**Readiness Probe** (`/ready`):
- Returns `503 Service Unavailable` until setup is complete
- Returns `200 OK` once all setup functions have succeeded and processes are running
- Indicates the service is ready to accept traffic

## API Reference

### Creating a Group

```go
// Create group with OS signal handling
group := service.NewSignals(os.Interrupt, os.Kill)

// Create group with custom context
ctx := context.WithTimeout(context.Background(), 30*time.Second)
group := service.NewCtx(ctx)
```

### Registering Functions

```go
// Setup: runs sequentially before processes start
group.Setup(func(ctx context.Context) error {
    // Initialization logic
    return nil
})

// Register: runs concurrently after setup completes
group.Register(func(ctx context.Context) error {
    // Service logic
    <-ctx.Done() // Wait for cancellation
    return nil
})
```

### Starting and Stopping

```go
// Start the group (blocks until all processes complete)
if err := group.Start(); err != nil {
    log.Fatal(err)
}

// Manually trigger shutdown
if err := group.Close(); err != nil {
    log.Fatal(err)
}
```

## How It Works

1. **Setup Phase**: All registered setup functions run sequentially. If any fails, the group stops immediately.
2. **Process Phase**: All registered processes run concurrently using `errgroup`. If any process returns an error, the context is canceled and all processes are stopped.
3. **Shutdown**: When the context is canceled (via signal, error, or manual close), all processes receive the cancellation signal and should clean up gracefully.

## Examples

See the [examples/server](examples/server/main.go) directory for a complete working example of an HTTP server with health check probes.

## Development

### Requirements

- Go 1.13 or later
- `golangci-lint` for linting

### Running Tests

```bash
# Run all tests (lint + unit + system)
make test

# Run only unit tests
make test-unit

# Run only system tests
make test-system
```

### Code Formatting

```bash
# Check linting
make lint

# Auto-fix and format code
make fmt
```

## Architecture

For detailed architecture documentation including sequence diagrams and design decisions, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Contributing

This repository is maintained by [@Attest/audience-accounts](CODEOWNERS).

For AI agents and automated tools:
- See [AGENTS.md](AGENTS.md) for comprehensive repository documentation
- See [CLAUDE.md](CLAUDE.md) for Claude-specific instructions

## License

Copyright © Attest
