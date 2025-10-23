// Fault Injection CLI Entry Point
// Reused: github.com/urfave/cli/v2 (project's CLI framework from flags/flags.go)
// Reused: tx-fuzz/flags package for flag definitions
// Reused: tx-fuzz/fault_injection packages for injection logic
// Logging style: fmt.Printf (consistent with project's helper/helper.go and spammer/)
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MariusVanDerWijden/tx-fuzz/fault_injection/os_command"
	"github.com/MariusVanDerWijden/tx-fuzz/fault_injection/rpc_command"
	"github.com/MariusVanDerWijden/tx-fuzz/flags"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "injection",
		Usage: "Fault injection tool for Ethereum clients",
		Flags: flags.InjectionFlags,
		Action: func(c *cli.Context) error {
			return runInjection(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runInjection(c *cli.Context) error {
	ctx := context.Background()

	mode := c.String(flags.InjectionModeFlag.Name)
	target := c.String(flags.InjectionTargetFlag.Name)
	param := c.String(flags.InjectionParamFlag.Name)
	containerID := c.String(flags.ContainerIDFlag.Name)
	endpoint := c.String(flags.RpcFlag.Name)

	// Validate required flags
	if target == "" {
		return fmt.Errorf("--injection-target is required")
	}

	fmt.Printf("=== Fault Injection Tool ===\n")
	fmt.Printf("Mode: %s\n", mode)
	fmt.Printf("Target: %s\n", target)
	if param != "" {
		fmt.Printf("Parameter: %s\n", param)
	}
	fmt.Printf("============================\n\n")

	switch mode {
	case "os":
		return handleOSInjection(ctx, target, containerID, param)
	case "rpc":
		return handleRPCInjection(ctx, target, endpoint, param)
	default:
		return fmt.Errorf("invalid injection-mode: %s (must be 'os' or 'rpc')", mode)
	}
}

func handleOSInjection(ctx context.Context, target, containerID, param string) error {
	if containerID == "" {
		return fmt.Errorf("--container-id is required for os mode")
	}

	start := time.Now()
	var err error

	switch target {
	case "restart":
		err = os_command.RestartContainer(ctx, containerID)
	case "stop":
		err = os_command.StopContainer(ctx, containerID)
	case "pause":
		err = os_command.PauseContainer(ctx, containerID)
	case "unpause":
		err = os_command.UnpauseContainer(ctx, containerID)
	case "kill":
		err = os_command.KillContainer(ctx, containerID)
	default:
		return fmt.Errorf("unknown os target: %s (valid: restart, stop, pause, unpause, kill)", target)
	}

	if err != nil {
		return fmt.Errorf("os injection '%s' failed: %w", target, err)
	}

	fmt.Printf("\n✓ OS injection '%s' completed successfully in %v\n", target, time.Since(start))
	return nil
}

func handleRPCInjection(ctx context.Context, target, endpoint, param string) error {
	start := time.Now()
	var err error

	switch target {
	case "setHead":
		if param == "" {
			return fmt.Errorf("--injection-param (block number/hash) is required for setHead")
		}
		err = rpc_command.SetHead(ctx, endpoint, param)

	case "clearTxPool":
		err = rpc_command.ClearTxPool(ctx, endpoint)

	case "triggerFork":
		if param == "" {
			return fmt.Errorf("--injection-param (block hash) is required for triggerFork")
		}
		err = rpc_command.TriggerFork(ctx, endpoint, param)

	case "stopRPC":
		fmt.Println("WARNING: This will stop the RPC server!")
		err = rpc_command.StopRPC(ctx, endpoint)

	case "pauseRPC":
		if param == "" {
			return fmt.Errorf("--injection-param (duration in seconds) is required for pauseRPC")
		}
		err = rpc_command.PauseRPC(ctx, endpoint, param)

	case "dropPeers":
		if param == "" {
			return fmt.Errorf("--injection-param (enode URL) is required for dropPeers")
		}
		err = rpc_command.DropPeers(ctx, endpoint, param)

	case "freeMemory":
		err = rpc_command.FreeMemory(ctx, endpoint)

	default:
		return fmt.Errorf("unknown rpc target: %s (valid: setHead, clearTxPool, triggerFork, stopRPC, pauseRPC, dropPeers)", target)
	}

	if err != nil {
		return fmt.Errorf("rpc injection '%s' failed: %w", target, err)
	}

	fmt.Printf("\n✓ RPC injection '%s' completed successfully in %v\n", target, time.Since(start))
	return nil
}
