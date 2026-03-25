# Testing Rules

## Test Philosophy

- Tests should be fast, deterministic, and isolated
- Use table-driven tests for multiple similar cases
- Test both success and failure paths
- Integration tests verify real binary behavior

## Test Organization

### Unit Tests
- Located alongside source: `group_test.go`, `probes_test.go`
- Test package naming: `package service_test` (black-box testing)
- Each test file tests its corresponding source file

### System Tests
- Located in `system/` directory: `system/system_test.go`
- Test compiled binary behavior (integration tests)
- Verify signal handling and HTTP endpoints work end-to-end

## Running Tests

```bash
# All tests (lint + unit + system)
make test

# Unit tests only
make test-unit

# System tests only
make test-system
```

All tests run with:
- Race detector enabled (`-race`)
- 1-minute timeout (`-timeout 1m`)

## Test Naming

### Test Functions
- Start with `Test`: `TestGroup`, `TestProbes`, `TestBinary`
- Use subtests for variations: `t.Run("success", func(t *testing.T) {...})`
- Describe what is being tested, not implementation details

### Helper Functions
- Use descriptive names: `block`, `verifyStatusCode`, `mustGetStatus`
- Mark test helpers with `t.Helper()`
- Keep helpers in the same file or `_test.go` file

## Test Structure

### Subtests Pattern
```go
func TestGroup(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        // Test happy path
    })

    t.Run("setup before register", func(t *testing.T) {
        // Test specific behavior
    })
}
```

### Cleanup Pattern
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // Clean up resources

// Or use Close() method
defer g.Close()
```

## Assertions

- Use clear, descriptive error messages
- Include expected vs actual values: `"expected %d, got: %d"`
- Use `t.Helper()` in assertion functions
- Prefer `t.Error()` for non-fatal assertions
- Use `t.Fatal()` for setup failures that prevent test continuation

## Testing Concurrency

### Channels for Synchronization
```go
done := make(chan struct{})
go func() {
    defer close(done)
    // Async work
}()

select {
case <-done:
    // Success
case <-time.After(timeout):
    t.Fatal("timeout")
}
```

### Testing Blocking Behavior
```go
select {
case <-done:
    t.Error("expected not done")
default:
    // Still blocking, as expected
}
```

### Eventually Pattern
For asynchronous operations that need time to settle:
```go
func eventually(t *testing.T, attempts int, f func() bool) {
    t.Helper()
    for i := 0; i < attempts; i++ {
        if f() {
            return
        }
        <-time.After(500 * time.Millisecond * time.Duration(i+1))
    }
    t.Fatal("expected success")
}
```

## System Tests

### Binary Compilation
- System tests require building `examples/server` first
- Build target: `build/server`
- Command: `make build-test`

### Testing Binaries
```go
cmd := exec.Command("../build/server")
cmd.Env = []string{"EXAMPLE_HTTP_ADDR=:4444"}
if err := cmd.Start(); err != nil {
    t.Fatalf("starting command: %v", err)
}
defer cmd.Process.Kill()
```

### HTTP Endpoint Testing
Use retry logic for services that need startup time:
```go
func mustGetStatus(t *testing.T, url string, attempts int, status int) {
    t.Helper()
    var got int
    for i := 0; i < attempts; i++ {
        resp, err := http.Get(url)
        if err == nil && resp.StatusCode == status {
            return
        }
        if err == nil {
            got = resp.StatusCode
        }
        <-time.After(time.Duration(500*(i+1)) * time.Millisecond)
    }
    t.Errorf("url: %q, expected: %d, got: %d", url, status, got)
}
```

## Race Detector

- Always enabled via `-race` flag
- Catches data races in concurrent code
- May slow tests significantly (acceptable trade-off)
- Fix all race conditions immediately

## Test Coverage

- Aim for high coverage of public API
- Test error paths and edge cases
- Don't obsess over 100% coverage
- Focus on testing behavior, not implementation

## Common Patterns

### Testing Probes
```go
// Verify status code changes over time
verifyStatusCode(t, "live", g.Liveness(), http.StatusServiceUnavailable)
go func() { _ = g.Start() }()
verifyStatusCode(t, "live", g.Liveness(), http.StatusOK)
```

### Testing Context Cancellation
```go
ctx, cancel := context.WithCancel(context.Background())
g := service.NewCtx(ctx)
g.Register(func(ctx context.Context) error {
    <-ctx.Done()
    return nil
})

cancel() // Trigger shutdown
if err := g.Start(); err != nil {
    t.Error(err)
}
```

## Avoiding Flaky Tests

- Use proper synchronization (channels, contexts)
- Avoid hardcoded sleeps (use retry loops instead)
- Make timeouts generous enough for CI environments
- Test cleanup should always happen (use defer)
