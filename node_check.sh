#!/usr/bin/env bash
# node_check.sh
# 简单的以太坊节点RPC健康检测脚本（适用于 Geth/Erigon/Nethermind/其他）
#
# 用法:
#   ./node_check.sh <RPC_ENDPOINT> [BLOCK_CHECK_INTERVAL_SECONDS]
# 示例:
#   ./node_check.sh http://127.0.0.1:8545 12

set -euo pipefail

RPC_ENDPOINT="${1:-}"
BLOCK_CHECK_INTERVAL="${2:-5}"  # 秒，测量区块增长的间隔

if [ -z "$RPC_ENDPOINT" ]; then
  echo "Usage: $0 <RPC_ENDPOINT> [BLOCK_CHECK_INTERVAL_SECONDS]"
  exit 2
fi

CURL_OPTS=(-sS -X POST -H "Content-Type: application/json")

jq_or_exit() {
  if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq not installed. Install with 'apt-get install -y jq' or similar."
    exit 1
  fi
}

jq_or_exit

rpc_call() {
  local method="$1"
  shift
  local params_json
  if [ "$#" -gt 0 ]; then
    params_json="$*"
  else
    params_json="[]"
  fi
  # build payload
  local payload
  payload=$(jq -nc --arg m "$method" --argjson p "$params_json" '{"jsonrpc":"2.0","method":$m,"params":$p,"id":1}')
  curl "${CURL_OPTS[@]}" --data "$payload" "$RPC_ENDPOINT"
}

print_header() {
  echo "------------------------------------------------------------"
  echo "$1"
  echo "------------------------------------------------------------"
}

# 1) RPC reachable & client version
print_header "1) RPC reachable & client version"
resp=$(rpc_call "web3_clientVersion" '[]' 2>/dev/null) || {
  echo "[FAIL] Cannot reach RPC at $RPC_ENDPOINT"
  exit 3
}
client_version=$(echo "$resp" | jq -r '.result // "N/A"')
echo "RPC reachable. clientVersion: $client_version"

# 2) eth_syncing
print_header "2) eth_syncing"
resp=$(rpc_call "eth_syncing" '[]')
syncing=$(echo "$resp" | jq '.result')
if [ "$syncing" = "false" ]; then
  echo "Node reports: synced (eth_syncing = false)"
else
  echo "Node reports syncing: $(echo "$syncing" | jq -r tostring)"
fi

# 3) block number (and growth check)
print_header "3) block number and progress check"
blk_hex=$(rpc_call "eth_blockNumber" '[]' | jq -r '.result')
if [ -z "$blk_hex" ] || [ "$blk_hex" = "null" ]; then
  echo "[WARN] eth_blockNumber returned empty"
else
  blk_dec=$((blk_hex))
  echo "Current block: ${blk_dec} (hex:${blk_hex})"
  echo "Waiting ${BLOCK_CHECK_INTERVAL}s to check progress..."
  sleep "${BLOCK_CHECK_INTERVAL}"
  blk2_hex=$(rpc_call "eth_blockNumber" '[]' | jq -r '.result')
  blk2_dec=$((blk2_hex))
  if [ "$blk2_dec" -gt "$blk_dec" ]; then
    echo "Block progressed: ${blk_dec} -> ${blk2_dec} (+$((blk2_dec-blk_dec)))"
  elif [ "$blk2_dec" -eq "$blk_dec" ]; then
    echo "[WARN] Block did NOT progress in ${BLOCK_CHECK_INTERVAL}s (still ${blk2_dec})."
  else
    echo "[WARN] Block number decreased: ${blk_dec} -> ${blk2_dec} (possible reorg/rollback?)"
  fi
fi

# 4) check recent block txs and one tx receipt if exists
print_header "4) recent block transactions check"
if [ -n "${blk_hex}" ] && [ "${blk_hex}" != "null" ]; then
  block_json=$(rpc_call "eth_getBlockByNumber" "[\"${blk_hex}\", true]")
  tx_count=$(echo "$block_json" | jq '.result.transactions | length')
  echo "Transactions in latest block: ${tx_count}"
  if [ "$tx_count" -gt 0 ]; then
    # take first tx hash and fetch receipt
    first_tx_hash=$(echo "$block_json" | jq -r '.result.transactions[0].hash')
    echo "Sample tx hash: $first_tx_hash"
    receipt=$(rpc_call "eth_getTransactionReceipt" "[\"${first_tx_hash}\"]")
    echo "Transaction receipt:"
    echo "$receipt" | jq -C '.result'
  else
    echo "Latest block has no transactions."
  fi
else
  echo "[WARN] Skipping block tx check because eth_blockNumber failed."
fi

# 5) txpool status / content (Geth specific)
print_header "5) txpool status (if supported)"
txpool_status_resp=$(rpc_call "txpool_status" '[]' 2>/dev/null || true)
if [ -z "$txpool_status_resp" ] || [ "$txpool_status_resp" = "" ]; then
  echo "txpool_status not supported or not exposed by this node."
else
  echo "txpool_status:"
  echo "$txpool_status_resp" | jq -C '.result'
fi

# 6) peer count
print_header "6) peer count"
peer_count_hex=$(rpc_call "net_peerCount" '[]' | jq -r '.result')
if [ -n "$peer_count_hex" ] && [ "$peer_count_hex" != "null" ]; then
  peer_count=$((peer_count_hex))
  echo "Peer count: ${peer_count}"
else
  echo "net_peerCount unavailable."
fi

# 7) debug_getBadBlocks (if supported)
print_header "7) debug_getBadBlocks (if supported)"
badblocks_resp=$(rpc_call "debug_getBadBlocks" '[]' 2>/dev/null || true)
if [ -z "$badblocks_resp" ] || [ "$badblocks_resp" = "" ]; then
  echo "debug_getBadBlocks not supported or not exposed."
else
  echo "debug_getBadBlocks output:"
  echo "$badblocks_resp" | jq -C '.result'
fi

# 8) optional: engine/forkchoice check (non-invasive)
print_header "8) engine/forkchoiceUpdated sanity (read-only check)"
# try a read of engine API via engine_getPayloadV1? Many implementations require auth; we will just check if engine namespace exists by requesting an invalid payload and see error
engine_resp=$(rpc_call "engine_getPayloadV1" '[null]' 2>/dev/null || true)
if [ -z "$engine_resp" ] || [ "$engine_resp" = "" ]; then
  echo "engine namespace not available or restricted."
else
  echo "engine namespace responded (may be available):"
  echo "$engine_resp" | jq -C '.'
fi

# 9) final summary
print_header "SUMMARY"
echo "RPC endpoint: $RPC_ENDPOINT"
echo "Client: $client_version"
echo "Block: ${blk_hex} (now check finished)"
echo "Peer count: ${peer_count:-N/A}"
echo "txpool_status: $(echo "$txpool_status_resp" | jq -r '.result // "N/A"')"
echo "Notes:"
echo " - If txpool_status/debug_getBadBlocks/engine API returned 'not supported', it may be disabled or RPC namespace not exposed."
echo " - For storage-level checks (DB integrity, LevelDB/RocksDB errors), check node logs and consider using debug/block export commands."
echo ""
echo "Health checks completed."

exit 0