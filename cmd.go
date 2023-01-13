package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func runCmd(
	timeout time.Duration,
	workingDir,
	name string,
	args ...string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = filepath.Join(workingDir)

	output, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("err: %v, %s", err, output)
	}

	if cmd.ProcessState != nil &&
		cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("code: %v, %s, %s", cmd.ProcessState.ExitCode(), output, cmd.ProcessState.String())
	}
	return nil
}
