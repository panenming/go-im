package exec

import (
	"os/exec"
	"time"
)

func ExecTimeout(d time.Duration, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	if d <= 0 {
		return cmd.Wait()
	}

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(d):
		cmd.Process.Kill()
		return <-done
	case err := <-done:
		return err
	}
}
