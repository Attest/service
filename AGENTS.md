# AGENTS.md

Single source of truth for AI agents working on this repository.

## Project Summary

`service` is a lightweight Go library that simplifies running and managing Go services. It provides a group-based lifecycle management system with built-in signal handling, graceful shutdown, and health check probes (liveness and readiness).

The library enables developers to:
- Register multiple service components that run concurrently
- Define setup functions that run before service components start
- Handle OS signals (SIGINT, SIGKILL) for graceful shutdown
- Expose HTTP liveness and readiness probes for orchestration systems
- Automatically stop all components if any one fails

## Tech Stack

| Category | Technology |
|----------|-----------|
| Language | Go 1.13+ |
| Concurrency | `golang.org/x/sync/errgroup` |
| Testing | Go standard testing, race detector |
| Build | Make, Go toolchain |
| Linting | `golangci-lint`, `go vet` |

## Repository Layout

```
service/
├── group.go           # Core Group type and lifecycle management
├── signal.go          # OS signal handling utilities
├── probes.go          # HTTP liveness and readiness handlers
├── group_test.go      # Unit tests for Group
├── probes_test.go     # Unit tests for probes
├── examples/
│   └── server/
│       └── main.go    # Example HTTP server implementation
├── system/
│   └── system_test.go # Integration tests against compiled binary
├── go.mod             # Go module definition
├── Makefile           # Build and test commands
├── README.md          # User-facing documentation
└── .gitignore         # Git ignore rules
```

## Key Entry Points

### Public API (`service` package)

**Core Types:**
- `Group` - Main service orchestration type
- `NewSignals(sig ...os.Signal) *Group` - Create group with signal handling
- `NewCtx(ctx context.Context) *Group` - Create group with custom context

**Lifecycle Methods:**
- `Setup(fn func(ctx context.Context) error)` - Register setup function (runs before processes)
- `Register(fn func(ctx context.Context) error)` - Register a service process
- `Start() error` - Start all setup functions, then all registered processes
- `Close() error` - Cancel context to trigger graceful shutdown

**Health Probes:**
- `Liveness() http.HandlerFunc` - Returns 503 until group started
- `Readiness() http.HandlerFunc` - Returns 503 until setup complete and group started

### Example Usage

See `examples/server/main.go` for a complete working example of an HTTP server with liveness/readiness probes.

## Development Workflow

### Setup
```bash
# Install dependencies
go mod download

# Install golangci-lint (required for linting)
# See: https://golangci-lint.run/usage/install/
```

### Testing
```bash
# Run all tests (lint + unit + system tests)
make test

# Run only unit tests
make test-unit

# Run only system tests (requires building example binary first)
make test-system

# Build example server binary
make build-test
```

### Linting and Formatting
```bash
# Run linters (go vet + golangci-lint + go mod tidy)
make lint

# Auto-fix linting issues and format code
make fmt
```

### Running Tests Manually
```bash
# Run tests with race detector
go test -race -timeout 1m ./...

# Run tests excluding system tests
go test -race -timeout 1m $(go list ./... | grep -v /system)

# Run only system tests
go test -race -timeout 1m ./system
```

## Agent Rules Reference

This repository follows specific coding conventions, testing practices, and security guidelines:

- **Code Style**: See `.agents/rules/code-style.md`
- **Testing**: See `.agents/rules/testing.md`
- **Security**: See `.agents/rules/security.md`

## Architecture

For detailed architecture documentation including sequence diagrams and component relationships, see `docs/ARCHITECTURE.md`.

## Updating This File

When making changes that affect repository structure, workflows, or core concepts:

1. Update the relevant section in this file
2. If adding new agent rules, create a new file in `.agents/rules/`
3. Update `CLAUDE.md` if new rule files are added
4. Keep this file as the single source of truth - avoid duplicating information in README.md

This file should be updated when:
- New major features are added
- Directory structure changes
- Development workflow changes
- Dependencies are added or significantly updated
- New conventions or rules are established
