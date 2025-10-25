# TX-Fuzz

TX-Fuzz is a package containing helpful functions to create random transactions. 
It can be used to easily access fuzzed transactions from within other programs.

## Usage

```
cd cmd/livefuzzer
go build
```

Run an execution layer client such as [Geth][1] locally in a standalone bash window.
Tx-fuzz sends transactions to port `8545` by default.

```
geth --http --http.port 8545
```

Run livefuzzer.

```
./livefuzzer spam
```

Tx-fuzz allows for an optional seed parameter to get reproducible fuzz transactions

## Advanced usage
You can optionally specify a seed parameter or a secret key to use as a faucet

```
./livefuzzer spam --seed <seed> --sk <SK>
```

You can set the RPC to use with `--rpc <RPC>`.

## Fault Injection Tool

TX-Fuzz 现在包含一个故障注入工具，用于对以太坊客户端进行各种故障注入测试。

### 快速开始

```bash
# 编译故障注入工具
go build -o bin/injection ./cmd/injection/main.go

# RPC 故障注入：回退链头
./bin/injection --injection-mode=rpc --injection-target=setHead --injection-param=0x64

# OS 故障注入：重启容器
./bin/injection --injection-mode=os --injection-target=restart --container-id=geth_container
```

### 功能特性

- **RPC 故障注入**：setHead（回退链）、clearTxPool（清空交易池）、triggerFork（分叉选择）等
- **OS 故障注入**：restart、stop、pause、unpause、kill Docker 容器
- **完全集成**：复用项目现有的 CLI 框架、RPC 客户端、flag 系统
- **简单易用**：每个注入操作都是独立的原子函数

