package main

import (
	"context"
	"fmt"

	"github.com/nivasan1/friends-are-for-losers/pkg/twitter" // Replace with the correct package path
)

func main() {
	bearerToken := "AAAAAAAAAAAAAAAAAAAAAMYdpgEAAAAA8slry2NwVWwh5hMFONM%2FK3Xbz6o%3DtvtpS38MSmRy9hk67L1vbza2PmlVk7EHldWxYTGT34EPOaGIam"
	tokenSecret := "YOUR_TOKEN_SECRET"
	accessToken := "YOUR_ACCESS_TOKEN"
	consumerSecret := "YOUR_CONSUMER_SECRET"
	consumerKey := "YOUR_CONSUMER_KEY"

	filter := twitter.NewTwitterFilter(bearerToken, tokenSecret, accessToken, consumerSecret, consumerKey)

	// Test username
	username := "yaper"

	isWorthy, err := filter.Filter(context.Background(), username) // or CheckWorthiness
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return
	}

	if isWorthy {
		fmt.Printf("The user %s is worth purchasing shares\n", username)
	} else {
		fmt.Printf("The user %s is not worth purchasing shares\n", username)
	}
}
