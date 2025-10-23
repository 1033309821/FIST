// Fault Injection OS Commands - Docker operations
// Fallback implementation: no existing docker utilities found in repository
// Uses standard library exec.Command to execute docker CLI commands
package os_command

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// RestartContainer restarts a docker container by ID
func RestartContainer(ctx context.Context, containerID string) error {
	return runDockerCmd(ctx, "restart", containerID)
}

// StopContainer stops a docker container by ID
func StopContainer(ctx context.Context, containerID string) error {
	return runDockerCmd(ctx, "stop", containerID)
}

// PauseContainer pauses a docker container by ID
func PauseContainer(ctx context.Context, containerID string) error {
	return runDockerCmd(ctx, "pause", containerID)
}

// UnpauseContainer unpauses a docker container by ID
func UnpauseContainer(ctx context.Context, containerID string) error {
	return runDockerCmd(ctx, "unpause", containerID)
}

// KillContainer kills a docker container by ID
func KillContainer(ctx context.Context, containerID string) error {
	return runDockerCmd(ctx, "kill", containerID)
}

// runDockerCmd executes a docker command and logs the output
// Uses project's logging style: fmt.Printf
func runDockerCmd(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "docker", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	fmt.Printf("[Docker] Executing: docker %v\n", args)
	if err := cmd.Run(); err != nil {
		fmt.Printf("[Docker] Error: %v; Output: %s\n", err, out.String())
		return err
	}
	fmt.Printf("[Docker] Success: %s\n", out.String())
	return nil
}
