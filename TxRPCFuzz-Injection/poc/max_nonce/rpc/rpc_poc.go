package main

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	// Configuration: rpc use geth or nethermind default port
	// pk use a random initial test key from kurtosis
	rpcURL := "http://127.0.0.1:8545"
	privateKeyHex := "0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"

	// Connect to node
	client, err := rpc.Dial(rpcURL)
	if err != nil {
		panic(fmt.Sprintf("Connection failed: %v", err))
	}
	defer client.Close()

	backend := ethclient.NewClient(client)

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(fmt.Sprintf("Private key parsing failed: %v", err))
	}

	sender := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("Sender address: %s\n", sender.Hex())

	// Get chain ID
	chainID, err := backend.ChainID(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Failed to get chain ID: %v", err))
	}

	// Get correct nonce (for comparison)
	correctNonce, err := backend.PendingNonceAt(context.Background(), sender)
	if err != nil {
		panic(fmt.Sprintf("Failed to get correct nonce: %v", err))
	}
	fmt.Printf("Correct nonce: %d\n", correctNonce)

	// Set nonce to maximum value
	maxNonce := uint64(math.MaxUint64)
	fmt.Printf("Using max nonce: %d\n", maxNonce)

	// Get gas price
	gasPrice, err := backend.SuggestGasPrice(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Failed to get gas price: %v", err))
	}

	// Construct transfer transaction
	to := common.HexToAddress("0xf93Ee4Cf8c6c40b329b0c0626F28333c132CF241") // random initial test address in kurtosis
	value := big.NewInt(1)
	gasLimit := uint64(21000)

	// Create transaction
	tx := types.NewTransaction(maxNonce, to, value, gasLimit, gasPrice, nil)

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		panic(fmt.Sprintf("Signing failed: %v", err))
	}

	// Serialize transaction
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf("Serialization failed: %v", err))
	}

	fmt.Printf("Raw transaction: %s\n", hexutil.Encode(rawTx))

	// Send transaction
	var txHash common.Hash
	err = client.CallContext(context.Background(), &txHash, "eth_sendRawTransaction", hexutil.Encode(rawTx))
	if err != nil {
		fmt.Printf("Send failed: %v\n", err)
	} else {
		fmt.Printf("Send successful! Transaction hash: %s\n", txHash.Hex())
	}
}
