package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	ENDPOINT_URL    string
	PRIVATE_KEY     string
	PUBLIC_ADDR     string
	CONTRACT_ADDRESS = common.HexToAddress("0xCF205808Ed36593aa40a44F10c7f7C2F67d4A4d4")
	hexcode         = "0x6945b123000000000000000000000000"
	u256            = "0000000000000000000000000000000000000000000000000000000000000001"
)

type Target struct {
	Address string `json:"address"`
	TwitterUsername string `json:"twitterUsername"`
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ENDPOINT_URL = os.Getenv("ENDPOINT_URL")
	PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
	PUBLIC_ADDR = os.Getenv("PUBLIC_KEY")
}

func userToAddr(user string) (string, error) {
	resp, err := http.Get("https://prod-api.kosetto.com/search/users?username=" + user)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var targets []Target
	err = json.Unmarshal(body, &targets)
	if err != nil {
		return "", err
	}

	for _, target := range targets {
		if target.TwitterUsername == user {
			return target.Address, nil
		}
	}

	return "", nil
}

func main() {
	loadEnv()
	client, err := ethclient.Dial(ENDPOINT_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// TODO: Continue with Ethereum logic using go-ethereum

	// For demonstration purposes, using the userToAddr function
	address, err := userToAddr("testUser")
	if err != nil {
		log.Fatalf("Failed to get address: %v", err)
	}
	fmt.Println(address)
}

// TODO: Add Ethereum functions, logic for buying shares, etc.
