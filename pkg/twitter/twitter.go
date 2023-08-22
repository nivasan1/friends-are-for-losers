package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type TwitterFilter struct {
	bearerToken    string
	tokenSecret    string
	accessToken    string
	consumerSecret string
	consumerKey    string
	client         *http.Client
}

func NewTwitterFilter(bearerToken, tokenSecret, accessToken, consumerSecret, consumerKey string) *TwitterFilter {
	// Define the OAuth2 configuration
	conf := &oauth2.Config{}

	// Create an OAuth2 token from the provided credentials
	token := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		RefreshToken: tokenSecret,
	}

	// Create an HTTP client that will authenticate using the OAuth2 token
	client := conf.Client(context.Background(), token)

	return &TwitterFilter{
		bearerToken:    bearerToken,
		tokenSecret:    tokenSecret,
		accessToken:    accessToken,
		consumerSecret: consumerSecret,
		consumerKey:    consumerKey,
		client:         client,
	}
}

// Filter is a wrapper around the twitter API that determines if a given
// twitter @ is worth purchasing shares
type Filter interface {
	Filter(ctx context.Context, address string) (bool, error)
}

func (f *TwitterFilter) Filter(ctx context.Context, username string) (bool, error) {
	url := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s?user.fields=public_metrics", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+f.bearerToken)

	resp, err := f.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			PublicMetrics struct {
				FollowersCount int `json:"followers_count"`
			} `json:"public_metrics"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Check if the followers count is greater than 5000
	if result.Data.PublicMetrics.FollowersCount > 5000 {
		return true, nil
	}

	return false, nil
}
