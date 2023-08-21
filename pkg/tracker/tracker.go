package tracker

import (
	"context"

	"github.com/nivasan1/friends-are-for-losers/pkg/types"
)

// Tracker tracks the base-scan friends.tech contract for new registrations / accounts being registered
type Tracker interface {
	Track(context.Context) <-chan *types.Registration
}
