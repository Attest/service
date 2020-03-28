package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zknill/service"
)

func TestProbes(t *testing.T) {
	g := service.NewCtx(context.Background())

	releaseSetup := make(chan struct{})

	g.Setup(func(ctx context.Context) error {
		<-releaseSetup
		return nil
	})

	verifyStatusCode(t, "live", g.Liveness(), http.StatusServiceUnavailable)
	verifyStatusCode(t, "ready", g.Readiness(), http.StatusServiceUnavailable)

	go func() { _ = g.Start() }()

	verifyStatusCode(t, "live", g.Liveness(), http.StatusOK)
	verifyStatusCode(t, "ready", g.Readiness(), http.StatusServiceUnavailable)

	close(releaseSetup)

	verifyStatusCode(t, "live", g.Liveness(), http.StatusOK)
	verifyStatusCode(t, "ready", g.Readiness(), http.StatusOK)
}

func verifyStatusCode(t *testing.T, handlerName string, h http.Handler, status int) {
	t.Helper()

	var got *http.Response

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		r := &http.Request{}

		h.ServeHTTP(w, r)

		got = w.Result()

		if got.StatusCode == status {
			return
		}

		d := 200 * time.Millisecond
		<-time.After(d * time.Duration(i+1))
	}

	t.Errorf("%q expected %d, got: %d", handlerName, status, got.StatusCode)
}
