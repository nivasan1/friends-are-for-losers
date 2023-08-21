package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	ENDPOINT_URL     = "YOUR_ETHEREUM_RPC_ENDPOINT"
	PRIVATE_KEY      = "YOUR_PRIVATE_KEY"
	PUBLIC_ADDR      = "YOUR_PUBLIC_ADDRESS"
	pollInterval     = 200 * time.Millisecond
	frenAPIEndpoint  = "https://prod-api.kosetto.com/search/users?username="
	GAS_LIMIT        = uint64(21000)
	GAS_PRICE        = uint64(50e9)
)

type Target struct {
	Amount int64   `json:"amount"`
	MaxBid float64 `json:"maxBid"`
}

type User struct {
	TwitterUsername string `json:"twitterUsername"`
	Address         string `json:"address"`
}

var targets = map[string]Target{
	"test1": {Amount: 3, MaxBid: 0.03},
	"test2": {Amount: 3, MaxBid: 0.07},
	"test3": {Amount: 1, MaxBid: 0.01},
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func makeAuth(privateKey, address string) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, err
	}
	auth := bind.NewKeyedTransactor(key)
	auth.From = common.HexToAddress(address)
	auth.GasLimit = GAS_LIMIT
	auth.GasPrice = big.NewInt(int64(GAS_PRICE))
	auth.Context = context.Background()
	return auth, nil
}

func sendTransaction(client *ethclient.Client, toAddress common.Address, auth *bind.TransactOpts) error {
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), auth.GasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, tx)
	if err != nil {
		return err
	}

	return client.SendTransaction(context.Background(), signedTx)
}

func userToAddr(userName string) (string, error) {
	resp, err := http.Get(frenAPIEndpoint + userName)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API response not OK: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", err
	}

	for _, u := range users {
		if u.TwitterUsername == userName {
			return u.Address, nil
		}
	}

	return "", fmt.Errorf("User not found")
}

func checkAndBuyShares() {
	client, err := ethclient.Dial(ENDPOINT_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	for target, details := range targets {
		address, err := userToAddr(target)
		if err != nil {
			log.Println("Error getting address for", target, ":", err)
			continue
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
		delete(targets, target)
	}
}

func main() {
	loadEnv()

	for {
		checkAndBuyShares()
		time.Sleep(pollInterval)
	}
}
