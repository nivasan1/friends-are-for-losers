package twitter

import (
	"context"
)

// Filter is a wrapper around the twitter API that determines if a given
// twitter @ is worth purchasing shares
type Filter interface {
	Filter(ctx context.Context, address string) (bool, error)
}
