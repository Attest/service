package system_test

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestBinary(t *testing.T) {
	startedCmd := func(env ...string) *exec.Cmd {
		cmd := exec.Command("../build/server")
		cmd.Env = env

		if err := cmd.Start(); err != nil {
			t.Fatalf("starting command: %v", err)
		}

		return cmd
	}

	t.Run("binary runs", func(t *testing.T) {
		cmd := startedCmd()

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			t.Fatalf("interrupting command: %v", err)
		}

		ps, err := cmd.Process.Wait()
		if err != nil {
			t.Fatalf("waiting for process %v", err)
		}

		if ps.ExitCode() != -1 {
			t.Errorf("expected exitcode -1, got: %d", ps.ExitCode())
		}
	})

	t.Run("liveness and readiness probes", func(t *testing.T) {
		cmd := startedCmd("EXAMPLE_HTTP_ADDR=:4444")
		cmd.Stdout = os.Stdout
		defer func() { _ = cmd.Process.Kill() }()

		mustGetStatus(t, "http://localhost:4444/live", 3, http.StatusOK)
		mustGetStatus(t, "http://localhost:4444/ready", 1, http.StatusOK)
	})
}

func mustGetStatus(t *testing.T, url string, attempts int, status int) {
	t.Helper()

	var got int

	for i := 0; i < attempts; i++ {
		resp, err := http.Get(url)
		if err == nil {
			if resp.StatusCode == status {
				return
			}

			got = resp.StatusCode
		}

		if err != nil {
			t.Log(err)
		}

		<-time.After(time.Duration(500*(i+1)) * time.Millisecond)
	}

	t.Errorf("url: %q, expected: %d, got: %d", url, status, got)
}
