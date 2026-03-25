# Code Style Rules

## General Principles

- Follow standard Go conventions and idioms
- Keep code simple, readable, and maintainable
- Prefer composition over inheritance
- Use short, descriptive variable names in small scopes

## Formatting

- Use `go fmt` for all code formatting (enforced via `make fmt`)
- Use tabs for indentation (Go standard)
- Keep line length reasonable (no hard limit, but use good judgment)

## Naming Conventions

### Variables
- Use short names in small scopes: `ctx`, `g`, `fn`, `err`
- Use descriptive names for package-level variables and types
- Boolean variables should read naturally: `started()`, `setup()`

### Functions
- Private helper functions: lowercase (e.g., `signalCtx`, `block`)
- Public API functions: PascalCase (e.g., `NewSignals`, `Liveness`)
- Function types: suffix with `Fn` (e.g., `runFn`)

### Types
- Exported types: PascalCase (e.g., `Group`)
- Private types: camelCase (rare in this codebase)

## Package Organization

- Keep related functionality in the same file
- File naming should reflect the main type or concept:
  - `group.go` - Group type and lifecycle
  - `probes.go` - Health check handlers
  - `signal.go` - Signal handling utilities
- Test files mirror source files: `group_test.go`, `probes_test.go`

## Error Handling

- Return errors, don't panic (unless truly exceptional)
- Don't ignore errors in production code
- In tests, use `t.Error()` or `t.Fatal()` appropriately
- Error returns should be the last return value

## Context Usage

- Always accept `context.Context` as the first parameter
- Use context for cancellation and deadline propagation
- Don't store contexts in structs (store cancel functions if needed)
- Check `ctx.Done()` in long-running operations

## Concurrency

- Use channels for coordination: `startedCh`, `setupCh`
- Use `sync.WaitGroup` or `errgroup.Group` for goroutine management
- Document goroutine lifecycle and ownership
- Prefer `errgroup.WithContext` for concurrent error handling

## Comments

- Document all exported functions, types, and methods
- Keep comments concise and up-to-date
- Use godoc conventions (package comment, function comments)
- Don't comment obvious code

## Testing Conventions

See `.agents/rules/testing.md` for detailed testing guidelines.

## Linting

All code must pass:
- `go vet ./...` - Standard Go static analysis
- `golangci-lint run` - Additional linters
- `go mod tidy` - Module dependencies are clean

Run `make lint` to check, `make fmt` to auto-fix issues.
