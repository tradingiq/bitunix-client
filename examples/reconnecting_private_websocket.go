package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"go.uber.org/zap"
)

type ExampleBalanceSubscriber struct{}

func (s *ExampleBalanceSubscriber) SubscribeBalance(msg *model.BalanceChannelMessage) {
	log.Printf("Received Balance Update: %+v", msg.Data)
}

type ExamplePositionSubscriber struct{}

func (s *ExamplePositionSubscriber) SubscribePosition(msg *model.PositionChannelMessage) {
	log.Printf("Received Position Update: %+v", msg.Data)
}

type ExampleOrderSubscriber struct{}

func (s *ExampleOrderSubscriber) SubscribeOrder(msg *model.OrderChannelMessage) {
	log.Printf("Received Order Update: %+v", msg.Data)
}

type ExampleTpSlOrderSubscriber struct{}

func (s *ExampleTpSlOrderSubscriber) SubscribeTpSlOrder(msg *model.TpSlOrderChannelMessage) {
	log.Printf("Received TP/SL Order Update: %+v", msg.Data)
}

func main() {
	ctx := context.Background()

	// Get API credentials from environment
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		log.Fatal("API_KEY and SECRET_KEY environment variables must be set")
	}

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create base private websocket client
	baseClient, err := bitunix.NewPrivateWebsocket(ctx, apiKey, secretKey,
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

	// Subscribe to channels (subscriptions are automatically restored on reconnection)
	if err := reconnectingClient.SubscribeBalance(balanceSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to balance updates: %v", err)
	}

	if err := reconnectingClient.SubscribePositions(positionSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to position updates: %v", err)
	}

	if err := reconnectingClient.SubscribeOrders(orderSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to order updates: %v", err)
	}

	if err := reconnectingClient.SubscribeTpSlOrders(tpSlOrderSubscriber); err != nil {
		log.Fatalf("Failed to subscribe to TP/SL order updates: %v", err)
	}

	log.Println("Connected and subscribed to all channels. Starting stream...")

	// Start streaming (this will automatically reconnect on failures)
	go func() {
		if err := reconnectingClient.Stream(); err != nil {
			log.Printf("Stream ended with error: %v", err)
		}
	}()

	// Keep running for demo
	select {}
}
