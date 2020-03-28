## service

_make it easy to run go services_

### How it works

The service groups allows registration of the different parts of a service.
These components are run in a group, if any one them fails, the group is stopped.

```go
group := service.NewSignals(os.Interrupt)

group.Setup(func(ctx context.Context) error) {
	// register a setup function.
	// all setup functions are run before registered functions.
	// if it returns an error, the group is stopped.
})

group.Register(func(ctx context.Context) error {
	// register a function that accepts a context
	// if it returns an error, the group is stopped
})

// the group exposes two http probes
// liveness reports service unavailable until all the setup functions have completed
// readiness reports service unavailable until all the registered functions have started
group.Liveness()
group.Readiness()

group.Register(httpServer())

// run the service
log.Fatal(group.Start())
```
