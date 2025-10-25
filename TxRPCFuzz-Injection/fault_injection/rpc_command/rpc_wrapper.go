// RPC Command Wrapper
// Reused: github.com/ethereum/go-ethereum/rpc (project standard RPC client)
// Similar to helper/helper.go GetRealBackend() pattern
package rpc_command

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

// dialRPC creates an RPC client connection
// Reuses project's standard: rpc.Dial() as seen in helper/helper.go and spammer/config.go
func dialRPC(endpoint string) (*rpc.Client, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC endpoint %s: %w", endpoint, err)
	}
	return client, nil
}

// callRPC makes a JSON-RPC call using the project's standard RPC client
// Similar to helper.go line 89: cl.CallContext(context.Background(), nil, "eth_sendRawTransaction", ...)
func callRPC(ctx context.Context, endpoint, method string, result interface{}, params ...interface{}) error {
	client, err := dialRPC(endpoint)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("[RPC] Calling %s with params: %v\n", method, params)
	if err := client.CallContext(ctx, result, method, params...); err != nil {
		fmt.Printf("[RPC] Error: %v\n", err)
		return err
	}
	fmt.Printf("[RPC] Success: %s completed\n", method)
	return nil
}
