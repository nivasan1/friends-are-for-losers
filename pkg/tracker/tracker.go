package tracker

import (
	"context"
	"fmt"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	fofl "github.com/nivasan1/friends-are-for-losers/pkg"

	"github.com/ethereum/go-ethereum"
	core "github.com/ethereum/go-ethereum/core/types"
	"github.com/nivasan1/friends-are-for-losers/pkg/types"
	"go.uber.org/zap"
)

const (
	FriendsTechAddress = "0xCF205808Ed36593aa40a44F10c7f7C2F67d4A4d4"
)

var Filter = ethereum.FilterQuery{
	Addresses: []common.Address{
		common.HexToAddress(FriendsTechAddress),
	},
	Topics: [][]common.Hash{{}},
}

var Codec abi.ABI

func init() {
	var err error
	Codec, err = abi.JSON(strings.NewReader(string(fofl.FriendsTechABI)))
	if err != nil {
		panic(err)
	}
}

//go:generate mockery --name EthClient --filename mock_ethclient.go
type EthClient interface {
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- core.Log) (ethereum.Subscription, error)
	Close()
}

// Tracker tracks the base-scan friends.tech contract for new registrations / accounts being registered.
//
//go:generate mockery --name Tracker --filename mock_tracker.go
type Tracker interface {
	Track(context.Context) <-chan *types.Registration
	Close() error
}

type trackerImpl struct {
	client        EthClient
	registrations chan *types.Registration
	logger        fofl.Logger
}

func DefaultTracker(ctx context.Context, dialAddr string, logger fofl.Logger) (Tracker, error) {
	client, err := ethclient.DialContext(ctx, dialAddr)
	if err != nil {
		return nil, err
	}

	return NewTracker(client, logger), nil
}

func NewTracker(client EthClient, logger fofl.Logger) Tracker {
	return &trackerImpl{
		client: client,
		logger: logger,
	}
}

func (t *trackerImpl) Track(ctx context.Context) <-chan *types.Registration {
	// create + return registration channel (non-blocking)
	registrations := make(chan *types.Registration)
	t.registrations = registrations
	// start a go-routine to listen to the events emitted from the contract
	go func() {
		// close the registrations channel when the context is cancelled
		defer close(registrations)

		logCh := make(chan core.Log)

		// create a new subscription to the contract
		sub, err := t.client.SubscribeFilterLogs(ctx, Filter, logCh)
		if err != nil {
			return
		}
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case log := <-logCh:
				t.logger.Info("received log", zap.String("address", log.Address.String()), zap.String("data", string(log.Data)), zap.String("topics", fmt.Sprintf("%v", log.Topics)))
				// decode the log data
				args, err := Codec.Unpack("Trade", log.Data)
				if err != nil {
					t.logger.Error("error decoding log data", zap.Error(err))
					continue
				}

				// check that the event is a registration event
				t.logger.Info("event args", zap.Any("args", args))

			case err := <-sub.Err():
				t.logger.Error("error subscribing to filter logs", zap.Error(err))
				return
			}
		}
	}()

	return registrations
}

func (t *trackerImpl) Close() error {
	if t.registrations == nil {
		return fmt.Errorf("tracker has not been initialized")
	}
	// check that the registration channel has been closed
	if fofl.CheckOpen(t.registrations) {
		close(t.registrations)
	}

	t.client.Close()
	return nil
}
