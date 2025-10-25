#!/bin/bash

# 故障注入测试脚本示例
# 此脚本展示如何使用故障注入工具进行各种测试

set -e

INJECTION_BIN="${INJECTION_BIN:-./bin/injection}"
RPC_ENDPOINT="${RPC_ENDPOINT:-http://127.0.0.1:8545}"

echo "========================================"
echo "TX-Fuzz 故障注入测试脚本"
echo "========================================"
echo "使用二进制: $INJECTION_BIN"
echo "RPC 端点: $RPC_ENDPOINT"
echo "========================================"
echo ""

# 检查二进制是否存在
if [ ! -f "$INJECTION_BIN" ]; then
    echo "错误: 找不到 $INJECTION_BIN"
    echo "请先编译: go build -o bin/injection ./cmd/injection/main.go"
    exit 1
fi

# 函数: RPC 健康检查
check_rpc() {
    echo ">>> 检查 RPC 连接..."
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$RPC_ENDPOINT" > /dev/null 2>&1; then
        echo "✓ RPC 连接正常"
        return 0
    else
        echo "✗ RPC 连接失败，请确保节点正在运行"
        return 1
    fi
}

# 测试 1: 获取当前区块号（验证 RPC 可用）
test_current_block() {
    echo ""
    echo "=== 测试 1: 获取当前区块号 ==="
    BLOCK=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$RPC_ENDPOINT" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$BLOCK" ]; then
        BLOCK_DEC=$((BLOCK))
        echo "当前区块: $BLOCK (十进制: $BLOCK_DEC)"
    else
        echo "无法获取区块号"
        return 1
    fi
}

# 测试 2: RPC 注入 - 回退链头（仅打印，不实际执行）
test_sethead() {
    echo ""
    echo "=== 测试 2: RPC 注入 - SetHead (演示模式) ==="
    echo "命令: $INJECTION_BIN --injection-mode=rpc --injection-target=setHead --injection-param=0x64 --rpc=$RPC_ENDPOINT"
    echo ""
    echo "⚠️  注意: SetHead 会回退区块链，影响节点状态"
    echo "如需实际执行，请取消下面的注释："
    echo "# $INJECTION_BIN --injection-mode=rpc --injection-target=setHead --injection-param=0x64 --rpc=$RPC_ENDPOINT"
}

# 测试 3: RPC 注入 - 清空交易池（仅打印，不实际执行）
test_clear_txpool() {
    echo ""
    echo "=== 测试 3: RPC 注入 - ClearTxPool (演示模式) ==="
    echo "命令: $INJECTION_BIN --injection-mode=rpc --injection-target=clearTxPool --rpc=$RPC_ENDPOINT"
    echo ""
    echo "⚠️  注意: 此操作会清空待处理交易"
    echo "如需实际执行，请取消下面的注释："
    echo "# $INJECTION_BIN --injection-mode=rpc --injection-target=clearTxPool --rpc=$RPC_ENDPOINT"
}

# 测试 4: OS 注入 - Docker 容器操作（需要容器 ID）
test_docker_ops() {
    echo ""
    echo "=== 测试 4: OS 注入 - Docker 操作 (演示模式) ==="
    echo ""
    echo "可用的 Docker 容器:"
    if command -v docker &> /dev/null; then
        docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}" 2>/dev/null || echo "无法列出容器"
    else
        echo "Docker 未安装或不可用"
    fi
    echo ""
    echo "示例命令:"
    echo "  重启容器: $INJECTION_BIN --injection-mode=os --injection-target=restart --container-id=<CONTAINER_ID>"
    echo "  暂停容器: $INJECTION_BIN --injection-mode=os --injection-target=pause --container-id=<CONTAINER_ID>"
    echo "  恢复容器: $INJECTION_BIN --injection-mode=os --injection-target=unpause --container-id=<CONTAINER_ID>"
    echo ""
    echo "⚠️  将 <CONTAINER_ID> 替换为实际的容器 ID"
}

# 测试 5: 显示帮助信息
test_help() {
    echo ""
    echo "=== 测试 5: 显示帮助信息 ==="
    $INJECTION_BIN --help || true
}

# 主流程
main() {
    # 只做非破坏性检查
    if check_rpc; then
        test_current_block
    fi
    
    test_sethead
    test_clear_txpool
    test_docker_ops
    test_help
    
    echo ""
    echo "========================================"
    echo "测试完成！"
    echo "========================================"
    echo ""
    echo "提示："
    echo "1. 上面的大部分测试是演示模式，未实际执行故障注入"
    echo "2. 如需实际测试，请编辑此脚本并取消相应命令的注释"
    echo "3. 建议在测试网络或独立环境中运行破坏性操作"
    echo "4. 详细文档: cmd/injection/README.md"
}

main "$@"

