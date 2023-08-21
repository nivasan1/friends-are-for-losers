package main

import (
	"fmt"
	"log"
	"os"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
	"context"
)

const (
	ENDPOINT_URL  = "YOUR_ETHEREUM_NODE_ENDPOINT"  // Replace this with your Ethereum RPC URL
	PRIVATE_KEY   = "YOUR_PRIVATE_KEY"             // Replace this with your Ethereum account private key
	PUBLIC_ADDR   = "YOUR_PUBLIC_ADDRESS"          // Replace this with your Ethereum public address
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func userToAddr(userName string) (string, error) {
	// Here, you would have some logic to map usernames to Ethereum addresses.
	// Just for illustration purposes, we'll assume a simple case.
	switch userName {
	case "testUser":
		return "0xSomeEthereumAddressHere", nil
	default:
		return "", fmt.Errorf("User not found")
	}
}

func makeAuth(privKeyHex, accountAddress string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.From = common.HexToAddress(accountAddress)
	auth.Nonce = nil     // Set nil to use pending nonce
	auth.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), privateKey) // Note: replace 1337 with your chain ID
	}
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(20000000000) // in wei

	return auth, nil
}

func sendTransaction(client *ethclient.Client, toAddress common.Address, auth *bind.TransactOpts) error {
	rawTx := types.NewTransaction(auth.Nonce.Uint64(), toAddress, auth.Value, auth.GasLimit, auth.GasPrice, nil)

	signedTx, err := auth.Signer(auth.From, rawTx)
	if err != nil {
		return err
	}

	return client.SendTransaction(context.TODO(), signedTx)
}

func main() {
	loadEnv()

	client, err := ethclient.Dial(ENDPOINT_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address, err := userToAddr("testUser")
	if err != nil {
		log.Fatalf("Failed to get address: %v", err)
	}

	toAddress := common.HexToAddress(address)
	auth, err := makeAuth(PRIVATE_KEY, PUBLIC_ADDR)
	if err != nil {
		log.Fatalf("Failed to create auth: %v", err)
	}

	err = sendTransaction(client, toAddress, auth)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Println("Transaction sent to address:", address)
}

