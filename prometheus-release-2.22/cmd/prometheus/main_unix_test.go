// +build !windows

package main

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

// As soon as prometheus starts responding to http request it should be able to
// accept Interrupt signals for a graceful shutdown.
func TestStartupInterrupt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	prom := exec.Command(promPath, "-test.main", "--config.file="+promConfig, "--storage.tsdb.path="+promData)
	err := prom.Start()
	if err != nil {
		t.Errorf("execution error: %v", err)
		return
	}

	done := make(chan error)
	go func() {
		done <- prom.Wait()
	}()

	var startedOk bool
	var stoppedErr error

Loop:
	for x := 0; x < 10; x++ {
		// error=nil means prometheus has started so we can send the interrupt
		// signal and wait for the graceful shutdown.
		if _, err := http.Get("http://localhost:9090/graph"); err == nil {
			startedOk = true
			prom.Process.Signal(os.Interrupt)
			select {
			case stoppedErr = <-done:
				break Loop
			case <-time.After(10 * time.Second):
			}
			break Loop
		}
		time.Sleep(500 * time.Millisecond)
	}

	if !startedOk {
		t.Errorf("prometheus didn't start in the specified timeout")
		return
	}
	if err := prom.Process.Kill(); err == nil {
		t.Errorf("prometheus didn't shutdown gracefully after sending the Interrupt signal")
	} else if stoppedErr != nil && stoppedErr.Error() != "signal: interrupt" { // TODO - find a better way to detect when the process didn't exit as expected!
		t.Errorf("prometheus exited with an unexpected error:%v", stoppedErr)
	}
}
