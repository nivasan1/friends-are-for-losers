package driver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nivasan1/friends-are-for-losers/pkg/driver"
	tracker_mocks "github.com/nivasan1/friends-are-for-losers/pkg/tracker/mocks"
	twitter_mocks "github.com/nivasan1/friends-are-for-losers/pkg/twitter/mocks"
	"github.com/nivasan1/friends-are-for-losers/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type DriverTestSuite struct {
	suite.Suite

	driver      *driver.Driver
	mockTracker *tracker_mocks.Tracker
	mockTwitter *twitter_mocks.Filter
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}

func (s *DriverTestSuite) SetupTest() {
	s.mockTracker = tracker_mocks.NewTracker(s.T())
	s.mockTwitter = twitter_mocks.NewFilter(s.T())

	s.driver = driver.NewDriver(s.mockTracker, s.mockTwitter, zap.NewNop())
}

// test that for each account that the tracker returns, the driver calls the twitter filter,
// and when the context for Run is cancelled the driver returns the context error.
func (s *DriverTestSuite) TestTrackerRegistrations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)

	registrations := make(chan *types.Registration)
	s.Run("when a tracker returns a registration, the driver calls the twitter filter", func() {
		s.mockTracker.On("Track", ctx).Return(receiver(registrations))
		s.mockTwitter.On("Filter", mock.Anything, "0x123").Return(true, nil)
		s.mockTwitter.On("Filter", mock.Anything, "0x456").Return(true, nil)
		s.mockTracker.On("Close").Return(nil).Run(func(args mock.Arguments) {
			close(registrations)
		}).Once()

		// start the driver in gr
		go func() {
			errCh <- s.driver.Run(ctx)
		}()

		registrations <- &types.Registration{
			Address: "0x123",
		}
	})

	s.Run("do another registration for shits and giggles", func() {
		registrations <- &types.Registration{
			Address: "0x456",
		}
	})

	s.Run("when the context is cancelled, the driver returns the context error", func() {
		cancel()
		s.EqualError(<-errCh, context.Canceled.Error())
		s.mockTracker.AssertExpectations(s.T())
	})
}

// test that if the twitter filter errors out, the driver returns the error and closes the tracker
func (s *DriverTestSuite) TestDriverFailures() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	registrations := make(chan *types.Registration)

	s.Run("when the twitter filter errors out, the driver returns the error and closes the tracker", func() {
		s.mockTracker.On("Track", ctx).Return(receiver(registrations))
		s.mockTwitter.On("Filter", mock.Anything, "0x123").Return(false, fmt.Errorf("error"))
		s.mockTracker.On("Close").Return(nil).Run(func(args mock.Arguments) {
			close(registrations)
		}).Once()

		// start the driver in gr
		go func() {
			errCh <- s.driver.Run(ctx)
		}()

		registrations <- &types.Registration{
			Address: "0x123",
		}

		s.EqualError(<-errCh, "error")
		s.mockTracker.AssertExpectations(s.T())
	})
}

func (s *DriverTestSuite) TestDriverQuitsWhenTrackerQuitsUnexpectedly() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registrations := make(chan *types.Registration)
	s.Run("when the Tracker quits unexpectedly, the driver quits", func() {
		s.mockTracker.On("Track", ctx).Return(receiver(registrations)).Run(func(args mock.Arguments) {
			// force tracker to quit immediately
			close(registrations)
		})
		s.mockTracker.On("Close").Return(fmt.Errorf("error")).Once()

		// start the driver in gr
		s.EqualError(s.driver.Run(ctx), "unexpected quit")
	})
}

func receiver[A any](ch chan A) <-chan A {
	return ch
}
