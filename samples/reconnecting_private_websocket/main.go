package main

import (
	"context"
	"log"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

type ExampleBalanceSubscriber struct{}

func (s *ExampleBalanceSubscriber) SubscribeBalance(msg *model.BalanceChannelMessage) {
	log.Printf("Received Balance Update: %+v", msg)
}

type ExamplePositionSubscriber struct{}

func (s *ExamplePositionSubscriber) SubscribePosition(msg *model.PositionChannelMessage) {
	log.Printf("Received Position Update: %+v", msg)
}

type ExampleOrderSubscriber struct{}

func (s *ExampleOrderSubscriber) SubscribeOrder(msg *model.OrderChannelMessage) {
	log.Printf("Received Order Update: %+v", msg)
}

type ExampleTpSlOrderSubscriber struct{}

func (s *ExampleTpSlOrderSubscriber) SubscribeTpSlOrder(msg *model.TpSlOrderChannelMessage) {
	log.Printf("Received TP/SL Order Update: %+v", msg)
}

func main() {
	ctx := context.Background()

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create base private websocket client
	baseClient, err := bitunix.NewPrivateWebsocket(ctx, samples.Config.ApiKey, samples.Config.SecretKey,
		bitunix.WithWebsocketLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create base private client: %v", err)
	}

	// Create reconnecting private websocket client
	reconnectingClient := bitunix.NewReconnectingPrivateWebsocket(ctx, baseClient,
		bitunix.WithPrivateMaxReconnectAttempts(0), // Infinite reconnect attempts
		bitunix.WithPrivateReconnectDelay(10*time.Second),
		bitunix.WithPrivateReconnectLogger(logger),
	)

	// Connect
	if err := reconnectingClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create subscribers
	balanceSubscriber := &ExampleBalanceSubscriber{}
	positionSubscriber := &ExamplePositionSubscriber{}
	orderSubscriber := &ExampleOrderSubscriber{}
	tpSlOrderSubscriber := &ExampleTpSlOrderSubscriber{}

	// Subscribe to all private channels
	if err := reconnectingClient.SubscribeBalance(balanceSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to balance: %v", err)
	}

	if err := reconnectingClient.SubscribePositions(positionSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to positions: %v", err)
	}

	if err := reconnectingClient.SubscribeOrders(orderSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to orders: %v", err)
	}

	if err := reconnectingClient.SubscribeTpSlOrders(tpSlOrderSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to TP/SL orders: %v", err)
	}

	log.Println("Successfully subscribed to all private channels")
	log.Println("The client will automatically reconnect and restore subscriptions on connection failures")

	// Start streaming (this will automatically reconnect on failures)
	go func() {
		if err := reconnectingClient.Stream(); err != nil {
			log.Printf("Stream ended with error: %v", err)
		}
	}()

	// Keep running for demo
	select {}
}
