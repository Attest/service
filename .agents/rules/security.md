# Security Rules

## General Security Principles

- This is a library package with no network exposure or secrets management
- Security concerns focus on safe API design and dependency management
- Follow secure Go coding practices

## Context Security

### Proper Context Cancellation
- Always respect context cancellation to prevent resource leaks
- Check `ctx.Done()` in long-running operations
- Never ignore context cancellation signals

```go
// Good: Respects cancellation
func process(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return nil // Clean shutdown
        default:
            // Do work
        }
    }
}
```

### Context Storage
- Never store contexts in structs (Go best practice)
- Pass contexts as function parameters
- Store cancel functions if needed, not contexts themselves

## Goroutine Safety

### Prevent Goroutine Leaks
- Every goroutine must have a clear termination path
- Use contexts for goroutine lifecycle management
- Always provide cleanup mechanisms

```go
// Good: Goroutine respects context
go func() {
    <-ctx.Done()
    _ = srv.Close() // Cleanup
}()
```

### Race Conditions
- Run tests with `-race` flag (enforced in Makefile)
- Protect shared state with proper synchronization
- Use channels for goroutine communication

## Signal Handling

### Safe Signal Processing
- Signal handlers must not panic
- Clean up resources before exit
- Use `signal.Stop()` to prevent leaks

```go
defer signal.Stop(c) // Always cleanup signal notifications
```

### Signal Channel Buffering
```go
c := make(chan os.Signal, len(sig)) // Buffered to prevent blocking
signal.Notify(c, sig...)
```

## Error Handling Security

### Don't Leak Sensitive Information
- Error messages should be descriptive but not expose internals
- In this library, errors are returned as-is (no sensitive data)
- Downstream users should wrap errors appropriately

### Fail Safely
- Default to stopping the service group if any component fails
- Don't continue processing if setup fails
- Ensure cleanup happens even on error paths

## Dependency Security

### Minimal Dependencies
- Only one external dependency: `golang.org/x/sync/errgroup`
- Review dependency updates carefully
- Keep `go.mod` and `go.sum` in sync

### Dependency Updates
```bash
# Check for updates
go list -u -m all

# Update dependencies
go get -u ./...
go mod tidy

# Verify go.sum
go mod verify
```

## HTTP Probe Security

### Probe Endpoints
- Liveness and readiness probes expose minimal information
- Only return HTTP status codes (200 or 503)
- No request body parsing or authentication needed
- Suitable for exposure to orchestration systems

### No Sensitive Data in Probes
```go
// Good: Only status code, no details
func (g Group) Liveness() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        if !g.started() {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
    }
}
```

## Testing Security

### Cleanup in Tests
- Always cleanup resources (contexts, goroutines, processes)
- Use `defer` for cleanup to ensure it happens
- Kill test processes properly to avoid zombies

```go
defer func() { _ = cmd.Process.Kill() }()
```

### Test Isolation
- Each test should be independent
- Don't rely on test execution order
- Clean up global state (signal handlers, etc.)

## Production Usage Security

### Library Consumer Responsibilities
When using this library, consumers should:
- Not expose probe endpoints publicly without authentication
- Handle secrets and configuration outside this library
- Implement proper logging without leaking sensitive data
- Use TLS for production HTTP servers

### Example Server Security
The `examples/server/main.go` is for demonstration only:
- No authentication or authorization
- No TLS
- Listens on all interfaces by default
- **Not production-ready without hardening**

## Resource Exhaustion Prevention

### Bounded Operations
- Use contexts with timeouts for bounded operations
- Don't create unbounded goroutines
- errgroup limits concurrent operations naturally

### Graceful Degradation
- Service group stops all components if one fails
- Prevents partial operation states
- Clean shutdown on errors

## Code Review Checklist

When reviewing changes:
- [ ] Are contexts propagated correctly?
- [ ] Are goroutines guaranteed to terminate?
- [ ] Is the race detector passing?
- [ ] Are resources cleaned up on all paths?
- [ ] Do error messages avoid leaking sensitive data?
- [ ] Are dependencies necessary and up-to-date?
