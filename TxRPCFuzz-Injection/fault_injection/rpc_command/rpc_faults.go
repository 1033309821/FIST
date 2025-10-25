// RPC Fault Injection Functions
// Reused: project's RPC client pattern from helper/helper.go and spammer packages
// Each function performs a specific fault injection via RPC
package rpc_command

import (
	"context"
	"fmt"
)

// SetHead calls debug_setHead to rewind the chain to a specific block
// Parameter: block number (hex string with 0x prefix or decimal)
func SetHead(ctx context.Context, endpoint, block string) error {
	fmt.Printf("[Inject] debug_setHead to block: %s\n", block)
	return callRPC(ctx, endpoint, "debug_setHead", nil, block)
}

// ClearTxPool attempts to clear the transaction pool (Geth-specific)
// May not be available in all client implementations
func ClearTxPool(ctx context.Context, endpoint string) error {
	fmt.Printf("[Inject] Clearing transaction pool\n")
	// Note: txpool_clear is not a standard method, trying txpool_content first to verify
	// Some clients use txpool_flush or don't expose this at all
	return callRPC(ctx, endpoint, "txpool_clear", nil)
}

// TriggerFork calls engine_forkchoiceUpdatedV2 to trigger a fork choice update
// Parameter: headBlockHash (hex string with 0x prefix)
func TriggerFork(ctx context.Context, endpoint, headHash string) error {
	fmt.Printf("[Inject] engine_forkchoiceUpdatedV2 with head: %s\n", headHash)
	// Minimal forkchoice state - can be extended as needed
	forkchoiceState := map[string]interface{}{
		"headBlockHash":      headHash,
		"safeBlockHash":      headHash,
		"finalizedBlockHash": headHash,
	}
	payloadAttributes := map[string]interface{}{}
	return callRPC(ctx, endpoint, "engine_forkchoiceUpdatedV2", nil, forkchoiceState, payloadAttributes)
}

// StopRPC calls admin_stopRPC to stop the RPC server (Geth-specific)
// WARNING: This will terminate the RPC endpoint, use with caution
func StopRPC(ctx context.Context, endpoint string) error {
	fmt.Printf("[Inject] admin_stopRPC - WARNING: This will stop the RPC server!\n")
	return callRPC(ctx, endpoint, "admin_stopRPC", nil)
}

// PauseRPC calls admin_sleep to pause RPC processing for a duration
// Parameter: duration in seconds as string
func PauseRPC(ctx context.Context, endpoint, duration string) error {
	fmt.Printf("[Inject] admin_sleep for %s seconds\n", duration)
	return callRPC(ctx, endpoint, "admin_sleep", nil, duration)
}

// DropPeers calls admin_removePeer to disconnect from peers
// Parameter: enode URL of the peer to remove
func DropPeers(ctx context.Context, endpoint, enodeURL string) error {
	fmt.Printf("[Inject] admin_removePeer: %s\n", enodeURL)
	return callRPC(ctx, endpoint, "admin_removePeer", nil, enodeURL)
}

func FreeMemory(ctx context.Context, endpoint string) error {
	fmt.Println("[Inject] debug_freeOSMemory - force GC and flush caches")
	return callRPC(ctx, endpoint, "debug_freeOSMemory", nil)
}
