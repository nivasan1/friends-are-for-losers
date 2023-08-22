package driver

import (
	"context"
	"sync"

	"github.com/nivasan1/friends-are-for-losers/pkg/tracker"
	"github.com/nivasan1/friends-are-for-losers/pkg/twitter"
	"github.com/nivasan1/friends-are-for-losers/pkg/types"
	"go.uber.org/zap"
)

type Driver struct {
	Tracker tracker.Tracker
	Twitter twitter.Filter
	logger  Logger
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

func NewDriver(tracker tracker.Tracker, twitter twitter.Filter, logger Logger) *Driver {
	return &Driver{
		Tracker: tracker,
		Twitter: twitter,
		logger:  logger,
	}
}

func (d *Driver) Run(ctx context.Context) error {
	registrations := d.Tracker.Track(ctx)
	done := make(chan struct{})

	// wait for all threads to finish after closing the tracker
	defer func() {
		// close the tracker when finished
		d.Tracker.Close()

		// wait for all registration threads to finish
		<-done
	}()

	errCh := make(chan error)

	go func() {
		// start wait-group
		wg := sync.WaitGroup{}

		// index through all registrations (fan-out)
		for registration := range registrations {
			wg.Add(1)

			// start a go-routine for each registration
			go func(registration *types.Registration) {
				defer wg.Done()

				// check if this twitter acct. is worth purchasing shares
				ok, err := d.Twitter.Filter(ctx, registration.Address)
				if err != nil {
					d.logger.Error("error filtering twitter", zap.Error(err))

					// send the error back to the main thread if there is an error here
					nonBlockingSend(errCh, err)
					return
				}

				if !ok {
					d.logger.Info("twitter filter returned false", zap.String("address", registration.Address))
					return
				}

				d.logger.Info("twitter filter returned true", zap.String("address", registration.Address))
				// TODO: purchase shares
			}(registration)
		}

		// wait for all threads to finish (fan-in)
		wg.Wait()
		close(done)
	}()

	// wait for all threads to finish
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func nonBlockingSend[A any](ch chan A, a A) {
	// first check that the channel isn't closed
	if !checkOpen(ch) {
		return
	}

	select {
	case ch <- a:
	default:
	}
}

func checkOpen[A any](ch chan A) bool {
	select {
	case _, ok := <-ch:
		return ok
	default: // this channel is open
		return true
	}
}
