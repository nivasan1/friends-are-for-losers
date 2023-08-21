package driver

import (
	"context"

	"github.com/nivasan1/friends-are-for-losers/pkg/tracker"
	"github.com/nivasan1/friends-are-for-losers/pkg/twitter"
	"go.uber.org/zap"
)

type Driver struct {
	Tracker tracker.Tracker
	Twitter twitter.Filter
	logger  zap.Logger
}

func (d *Driver) Run(ctx context.Context) error {
	registrations := d.Tracker.Track(ctx)
	errCh := make(chan error)

	for registration := range registrations {
		go func() {
			ok, err := d.Twitter.Filter(ctx, registration.Address)
			if err != nil {
				d.logger.Error("error filtering twitter", zap.Error(err))
				errCh <- err
				return
			}

			if !ok {
				d.logger.Info("twitter filter returned false", zap.String("address", registration.Address))
				return
			}

			d.logger.Info("twitter filter returned true", zap.String("address", registration.Address))
			// TODO: purchase shares
		}()
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			close(registrations)
			return err
		}
	}
	return nil
}
