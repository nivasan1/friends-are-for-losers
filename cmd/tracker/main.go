package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	fofl "github.com/nivasan1/friends-are-for-losers/pkg"
	"github.com/nivasan1/friends-are-for-losers/pkg/tracker"
	"go.uber.org/zap"
)

var addr = flag.String("addr", "ws://localhost:8546", "websocket address to dial for events")

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-sigCh
		cancel()
	}()

	tracker, err := tracker.DefaultTracker(ctx, *addr, zap.NewExample())
	if err != nil {
		panic(err)
	}
	defer tracker.Close()

	registrations := tracker.Track(ctx)

	select {
	case <-ctx.Done():
		return
	case <-registrations:
		if !fofl.CheckOpen(registrations) {
			return
		}
	}
}
