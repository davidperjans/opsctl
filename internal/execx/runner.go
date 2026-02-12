package execx

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
}

type OSRunner struct {
	Timeout time.Duration
}

func (r OSRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	// If the caller already set a deadline (from commands)
	// respect it and don't override

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		timeout := r.Timeout
		if timeout <= 5 {
			timeout = 5 * time.Second
		}

		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, name, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	if err != nil {
		// return stderr if stdout empty, both are useful for debugging
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = out
		}
		if msg == "" {
			return "", err
		}
		return msg, err
	}

	return out, nil
}
