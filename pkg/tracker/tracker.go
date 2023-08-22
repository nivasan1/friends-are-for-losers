package tracker

import (
	"context"

	"github.com/nivasan1/friends-are-for-losers/pkg/types"
)

// Tracker tracks the base-scan friends.tech contract for new registrations / accounts being registered.
//
//go:generate mockery --name Tracker --filename mock_tracker.go
type Tracker interface {
	Track(context.Context) <-chan *types.Registration
	Close() error
}
