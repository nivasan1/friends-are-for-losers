package twitter

import (
	"context"
)

// Filter is a wrapper around the twitter API that determines if a given
// twitter @ is worth purchasing shares
//
//go:generate mockery --name Filter --filename mock_filter.go
type Filter interface {
	Filter(ctx context.Context, address string) (bool, error)
}
