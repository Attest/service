package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Attest/service"
)

func TestGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		g := service.NewCtx(ctx)

		g.Register(block)

		errs := make(chan error)

		go func() {
			errs <- g.Start()
		}()

		cancel()
		if err := <-errs; err != nil {
			t.Errorf("expected no error: got %+v", err)
		}
	})

	t.Run("setup before register", func(t *testing.T) {
		g := service.NewCtx(context.Background())

		setupCalled := make(chan struct{})
		g.Setup(func(ctx context.Context) error {
			close(setupCalled)
			return nil
		})

		g.Register(func(ctx context.Context) error {
			select {
			case <-setupCalled:
			default:
				t.Error("expected setup called before register")
			}

			return nil
		})

		if err := g.Start(); err != nil {
			t.Error(err)
		}
	})

	t.Run("double call start", func(t *testing.T) {
		g := service.NewCtx(context.Background())
		g.Register(block)

		ch1 := make(chan struct{})
		ch2 := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer close(ch1)
			wg.Done()

			_ = g.Start()
		}()

		go func() {
			defer close(ch2)
			wg.Wait()

			_ = g.Start()
		}()

		select {
		case <-ch1:
			t.Error("expected 1 to be blocked")
		case <-ch2:
		}
	})

	t.Run("close", func(t *testing.T) {
		g := service.NewCtx(context.Background())
		g.Register(block)

		done := make(chan struct{})

		go func() {
			defer close(done)
			_ = g.Start()
		}()

		select {
		case <-done:
			t.Error("expected not done")
		default:
		}

		if err := g.Close(); err != nil {
			t.Error(err)
		}

		closed := func() bool {
			select {
			case <-done:
				return true
			default:
				return false
			}
		}

		eventually(t, 3, closed)
	})
}

func block(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func eventually(t *testing.T, attempts int, f func() bool) {
	t.Helper()

	for i := 0; i < 3; i++ {
		ok := f()
		if ok == true {
			return
		}

		d := 500 * time.Millisecond
		<-time.After(d * time.Duration(i+1))
	}

	t.Fatal("expected success")
}
